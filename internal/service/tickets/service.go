package tickets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rvxt21/queueing-system/config"
	"github.com/rvxt21/queueing-system/internal/broker"
	"github.com/rvxt21/queueing-system/internal/domain"
)

const (
	simulatedWork       = 15 * time.Second // simulatedWork stands in for real ticket processing.
	reconnectMinBackoff = time.Second
	reconnectMaxBackoff = 30 * time.Second
	reconnectBackoffMul = 2
	revertTimeout       = 5 * time.Second
)

type Worker struct {
	br     Broker
	repo   Repository
	notify Notifier
	log    *slog.Logger
	queues []string
	count  int

	cancel context.CancelFunc
	wg     sync.WaitGroup
	done   chan struct{}
}

func NewWorker(br Broker, repo Repository, notify Notifier, log *slog.Logger, cfg *config.Config) *Worker {
	return &Worker{
		br:     br,
		log:    log,
		repo:   repo,
		notify: notify,
		count:  cfg.Count,
		queues: broker.Queues(),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	w.cancel = cancel
	w.done = make(chan struct{})

	go w.work(ctx)

	return nil
}

func (w *Worker) Stop(ctx context.Context) error {
	w.cancel()

	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Worker) work(ctx context.Context) {
	defer close(w.done)
	defer w.wg.Wait()

	for i := range w.count {
		for _, queue := range w.queues {
			w.wg.Go(func() {
				tag := fmt.Sprintf("%s-worker-%d", queue, i) // e.g. "tickets.tax-worker-0"

				deliveries, err := w.br.Consume(ctx, queue, tag)
				if err != nil {
					w.log.Error("start consumer", slog.String("tag", tag), slog.Any("err", fmt.Errorf("%w: %w", ErrConsuming, err)))

					return
				}

				w.log.Info("consumer started", slog.String("queue", queue), slog.String("tag", tag))

				for {
					select {
					case <-ctx.Done():
						return
					case msg, ok := <-deliveries:
						if !ok {
							if ctx.Err() != nil {
								return
							}

							deliveries, ok = w.resubscribe(ctx, queue, tag)
							if !ok {
								return
							}

							continue
						}

						w.processMessage(ctx, msg, tag)
					}
				}
			})
		}
	}
}

// resubscribe calls Consume again with backoff until it succeeds.
// It returns false if ctx is cancelled first.
func (w *Worker) resubscribe(ctx context.Context, queue, tag string) (<-chan domain.Delivery, bool) {
	w.log.Warn("consumer lost, resubscribing", slog.String("tag", tag))

	backoff := reconnectMinBackoff

	for attempt := 1; ; attempt++ {
		deliveries, err := w.br.Consume(ctx, queue, tag)
		if err == nil {
			w.log.Info("consumer started",
				slog.String("queue", queue),
				slog.String("tag", tag),
				slog.Int("attempts", attempt))

			return deliveries, true
		}

		w.log.Warn("resubscribe failed",
			slog.String("tag", tag),
			slog.Int("attempt", attempt),
			slog.Duration("retry_in", backoff),
			slog.Any("err", err))

		select {
		case <-ctx.Done():
			return nil, false
		case <-time.After(backoff):
		}

		backoff = min(backoff*reconnectBackoffMul, reconnectMaxBackoff)
	}
}

// processMessage handles one ticket and acks or nacks it.
func (w *Worker) processMessage(ctx context.Context, msg domain.Delivery, tag string) {
	var event domain.TicketCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		w.log.Error("decode ticket event", slog.String("message_id", msg.MessageID), slog.Any("err", err))
		_ = msg.Nack(false)

		return
	}

	started := time.Now()

	err := w.repo.UpdateTicketStatus(ctx, event.ID, domain.TicketStatusDelivered)
	if err != nil {
		w.log.Error("mark ticket delivered",
			slog.String("ticket_id", event.ID.String()),
			slog.String("ticket_number", event.Number),
			slog.Any("err", err))
		_ = msg.Nack(!errors.Is(err, domain.ErrTicketNotFound))

		return
	}

	w.publishStatus(event, domain.TicketStatusDelivered, tag)

	w.log.Info("ticket processing started",
		slog.String("ticket_id", event.ID.String()),
		slog.String("ticket_number", event.Number),
		slog.String("tag", tag))

	select {
	case <-time.After(simulatedWork):
	case <-ctx.Done():
		w.revertToSent(ctx, event)

		_ = msg.Nack(true) // shutting down, put it back in the queue

		return
	}

	err = w.repo.UpdateTicketStatus(ctx, event.ID, domain.TicketStatusResolved)
	if err != nil {
		w.log.Error("mark ticket resolved",
			slog.String("ticket_id", event.ID.String()),
			slog.String("ticket_number", event.Number),
			slog.Any("err", err))

		notFound := errors.Is(err, domain.ErrTicketNotFound)
		if !notFound {
			w.revertToSent(ctx, event)
		}

		_ = msg.Nack(!notFound)

		return
	}

	w.publishStatus(event, domain.TicketStatusResolved, tag)

	if err := msg.Ack(); err != nil {
		w.log.Warn("ack failed, ticket will be redelivered",
			slog.String("ticket_id", event.ID.String()),
			slog.String("ticket_number", event.Number),
			slog.Any("err", err))

		return
	}

	w.log.Info("ticket resolved",
		slog.String("ticket_id", event.ID.String()),
		slog.String("ticket_number", event.Number),
		slog.Duration("took", time.Since(started).Round(time.Millisecond)))
}

func (w *Worker) revertToSent(ctx context.Context, event domain.TicketCreatedEvent) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), revertTimeout)
	defer cancel()

	err := w.repo.UpdateTicketStatus(ctx, event.ID, domain.TicketStatusSent)
	if err != nil {
		w.log.Warn("revert ticket to sent",
			slog.String("ticket_id", event.ID.String()),
			slog.String("ticket_number", event.Number),
			slog.Any("err", err))

		return
	}

	w.publishStatus(event, domain.TicketStatusSent, "")
}

func (w *Worker) publishStatus(event domain.TicketCreatedEvent, status domain.TicketStatus, worker string) {
	w.notify.Publish(domain.TicketStatusChanged{
		ID:          event.ID,
		Number:      event.Number,
		ServiceCode: event.ServiceCode,
		Status:      status,
		Worker:      worker,
	})
}
