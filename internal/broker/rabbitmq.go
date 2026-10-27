package broker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rvxt21/queueing-system/config"
)

const (
	reconnectMinBackoff = time.Second
	reconnectMaxBackoff = 30 * time.Second
	reconnectBackoffMul = 2
)

type RabbitMQ struct {
	config config.RabbitMQ
	log    *slog.Logger

	mu   sync.RWMutex // guards conn and ch, replaced on reconnect
	conn *amqp.Connection
	ch   *amqp.Channel

	cancel context.CancelFunc
	done   chan struct{}
}

func NewRabbitMQ(cfg *config.Config, log *slog.Logger) *RabbitMQ {
	return &RabbitMQ{
		config: cfg.RabbitMQ,
		log:    log,
	}
}

func (r *RabbitMQ) Start(ctx context.Context) error {
	conn, ch, err := r.connect()
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.conn, r.ch = conn, ch
	r.mu.Unlock()

	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	r.cancel = cancel
	r.done = make(chan struct{})

	go r.watch(ctx)

	r.log.Info("rabbitmq connected",
		slog.String("host", r.config.Host),
		slog.String("exchange", r.config.Exchange))

	return nil
}

func (r *RabbitMQ) Stop(ctx context.Context) error {
	// Stop watch first, so it doesn't treat our own close as a lost connection.
	r.cancel()

	select {
	case <-r.done:
	case <-ctx.Done():
		return ctx.Err()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var errs []error
	if r.ch != nil {
		errs = append(errs, r.ch.Close())
	}

	if r.conn != nil {
		errs = append(errs, r.conn.Close())
	}

	err := errors.Join(errs...)
	if err == nil {
		r.log.Info("rabbitmq connection closed")
	}

	return err
}

// connect dials RabbitMQ and declares the exchange, queues and bindings.
func (r *RabbitMQ) connect() (*amqp.Connection, *amqp.Channel, error) {
	dsn := (&url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(r.config.User, r.config.Password),
		Host:   net.JoinHostPort(r.config.Host, r.config.Port),
	}).String()

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("open channel: %w", errors.Join(err, conn.Close()))
	}

	if err := ch.Confirm(false); err != nil {
		return nil, nil, fmt.Errorf("enable publisher confirms: %w", errors.Join(err, ch.Close(), conn.Close()))
	}

	err = ch.ExchangeDeclare(
		r.config.Exchange,   // name
		amqp.ExchangeDirect, // type
		true,                // durability
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error exchange declare: %w", errors.Join(err, ch.Close(), conn.Close()))
	}

	if err := declareQueues(ch); err != nil {
		return nil, nil, errors.Join(err, ch.Close(), conn.Close())
	}

	if err := bindQueues(ch, r.config.Exchange); err != nil {
		return nil, nil, errors.Join(err, ch.Close(), conn.Close())
	}

	return conn, ch, nil
}

// watch reconnects whenever the connection drops, until ctx is cancelled.
func (r *RabbitMQ) watch(ctx context.Context) {
	defer close(r.done)

	for {
		closed := r.currentConn().NotifyClose(make(chan *amqp.Error, 1))

		select {
		case <-ctx.Done():
			return
		case reason := <-closed:
			logConnectionLost(r.log, reason)
		}

		lostAt := time.Now()
		backoff := reconnectMinBackoff

		for attempt := 1; ; attempt++ {
			conn, ch, err := r.connect()
			if err == nil {
				r.mu.Lock()
				r.conn, r.ch = conn, ch
				r.mu.Unlock()
				r.log.Info("rabbitmq reconnected",
					slog.Int("attempts", attempt),
					slog.Duration("downtime", time.Since(lostAt).Round(time.Millisecond)))

				break
			}

			r.log.Warn("rabbitmq reconnect failed",
				slog.Int("attempt", attempt),
				slog.Duration("retry_in", backoff),
				slog.Any("err", err))

			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}

			backoff = min(backoff*reconnectBackoffMul, reconnectMaxBackoff)
		}
	}
}

// logConnectionLost logs why the connection closed. reason is nil when it
// closed without an error from the server.
func logConnectionLost(log *slog.Logger, reason *amqp.Error) {
	if reason == nil {
		log.Warn("rabbitmq connection lost")

		return
	}

	log.Warn("rabbitmq connection lost",
		slog.Int("code", reason.Code),
		slog.String("reason", reason.Reason),
		slog.Bool("server_initiated", reason.Server))
}

func (r *RabbitMQ) currentConn() *amqp.Connection {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.conn
}

func (r *RabbitMQ) currentCh() *amqp.Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ch
}
