package outbox

import (
	"context"
	"log/slog"
	"time"

	"github.com/rvxt21/queueing-system/config"
)

type Worker struct {
	br   Broker
	repo Repository
	log  *slog.Logger

	interval  time.Duration
	batchSize int64
	cancel    context.CancelFunc
	done      chan struct{}
}

func NewWorker(br Broker, repo Repository, log *slog.Logger, cfg *config.Config) *Worker {
	return &Worker{
		br:        br,
		repo:      repo,
		log:       log,
		interval:  cfg.PollInterval,
		batchSize: cfg.BatchSize,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	w.cancel = cancel
	w.done = make(chan struct{})

	go w.run(ctx)

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

func (w *Worker) run(ctx context.Context) {
	defer close(w.done)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.publishMessages(ctx)
		}
	}
}
