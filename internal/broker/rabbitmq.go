package broker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rvxt21/queueing-system/config"
)

type RabbitMQ struct {
	config config.RabbitMQ
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func (r *RabbitMQ) Start(ctx context.Context) error {
	dsn := (&url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(r.config.User, r.config.Password),
		Host:   net.JoinHostPort(r.config.Host, r.config.Port),
	}).String()

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", errors.Join(err, conn.Close()))
	}

	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("enable publisher confirms: %w", errors.Join(err, ch.Close(), conn.Close()))
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
		return fmt.Errorf("error exchange declare: %w", errors.Join(err, ch.Close(), conn.Close()))
	}

	if err := declareQueues(ch); err != nil {
		return errors.Join(err, ch.Close(), conn.Close())
	}

	if err := bindQueues(ch, r.config.Exchange); err != nil {
		return errors.Join(err, ch.Close(), conn.Close())
	}

	r.conn, r.ch = conn, ch

	return nil
}

func (r *RabbitMQ) Stop(ctx context.Context) error {
	var errs []error
	if r.ch != nil {
		errs = append(errs, r.ch.Close())
	}

	if r.conn != nil {
		errs = append(errs, r.conn.Close())
	}

	return errors.Join(errs...)
}

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	return &RabbitMQ{
		config: cfg.RabbitMQ,
	}
}
