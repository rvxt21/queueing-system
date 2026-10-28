package tickets

import (
	"context"

	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type Broker interface {
	Consume(ctx context.Context, queue, consumerTag string) (<-chan domain.Delivery, error)
}

// Notifier tells live subscribers (the dashboard) about ticket status changes.
type Notifier interface {
	Publish(msg domain.TicketStatusChanged)
}

type Repository interface {
	UpdateTicketStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error
}
