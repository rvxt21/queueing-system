package outbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type Repository interface {
	GetOutboxes(ctx context.Context, limit int64) ([]domain.Outbox, error)
	MarkOutboxComplete(ctx context.Context, id uuid.UUID) error
	MarkOutboxFailed(ctx context.Context, id uuid.UUID, publishErr error) error
}

type Broker interface {
	Publish(ctx context.Context, msg *domain.Outbox) error
}
