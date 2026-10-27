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
)

type Worker struct {
	br     Broker
	repo   Repository
	log    *slog.Logger
	queues []string
	count  int

	cancel context.CancelFunc
	wg     sync.WaitGroup
	done   chan struct{}
}

func NewWorker(br Broker, repo Repository, log *slog.Logger, cfg *config.Config) *Worker {
	return &Worker{
		br:     br,
		log:    log,
		repo:   repo,
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

							w.log.Warn("consumer lost, resubscribing", slog.String("tag", tag))

							backoff := reconnectMinBackoff

							for attempt := 1; ; attempt++ {
								deliveries, err = w.br.Consume(ctx, queue, tag)
								if err == nil {
									w.log.Info("consumer started",
										slog.String("queue", queue),
										slog.String("tag", tag),
										slog.Int("attempts", attempt))

									break
								}

								w.log.Warn("resubscribe failed",
									slog.String("tag", tag),
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

							continue
						}

						var event domain.TicketCreatedEvent
						if err := json.Unmarshal(msg.Body, &event); err != nil {
							w.log.Error("decode ticket event", slog.String("message_id", msg.MessageID), slog.Any("err", err))
							_ = msg.Nack(false)

							continue
						}

						started := time.Now()

						w.log.Info("ticket processing started",
							slog.String("ticket_id", event.ID.String()),
							slog.String("ticket_number", event.Number),
							slog.String("tag", tag))

						select {
						case <-time.After(simulatedWork):
						case <-ctx.Done():
							_ = msg.Nack(true) // shutting down, put it back in the queue
							return
						}

						err := w.repo.UpdateTicketStatus(ctx, event.ID, domain.TicketStatusResolved)
						if err != nil {
							w.log.Error("update ticket status",
								slog.String("ticket_id", event.ID.String()),
								slog.String("ticket_number", event.Number),
								slog.Any("err", err))
							_ = msg.Nack(!errors.Is(err, domain.ErrTicketNotFound))

							continue
						}

						if err := msg.Ack(); err != nil {
							w.log.Warn("ack failed, ticket will be redelivered",
								slog.String("ticket_id", event.ID.String()),
								slog.String("ticket_number", event.Number),
								slog.Any("err", err))

							continue
						}

						w.log.Info("ticket resolved",
							slog.String("ticket_id", event.ID.String()),
							slog.String("ticket_number", event.Number),
							slog.Duration("took", time.Since(started).Round(time.Millisecond)))
					}
				}
			})
		}
	}
}
