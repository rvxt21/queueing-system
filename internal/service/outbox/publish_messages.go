package outbox

import (
	"context"
	"log/slog"
)

func (w *Worker) publishMessages(ctx context.Context) {
	messages, err := w.repo.GetOutboxes(ctx, w.batchSize)
	if err != nil {
		w.log.Error("fetch outbox events", slog.Any("err", err))

		return
	}

	for _, msg := range messages {
		if err := w.br.Publish(ctx, &msg); err != nil {
			if ctx.Err() != nil {
				return
			}

			w.log.Error("publish outbox event", slog.Any("id", msg.ID), slog.Any("err", err))

			if err := w.repo.MarkOutboxFailed(ctx, msg.ID, err); err != nil {
				w.log.Error("mark outbox failed", slog.Any("id", msg.ID), slog.Any("err", err))
			}

			return
		}

		if err := w.repo.MarkOutboxComplete(ctx, msg.ID); err != nil {
			w.log.Error("mark outbox complete", slog.Any("id", msg.ID), slog.Any("err", err))

			return
		}
	}
}
