package domain

import (
	"time"

	"github.com/google/uuid"
)

type Outbox struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	RoutingKey  string
	Payload     string
	CreatedAt   time.Time
	PublishedAt *time.Time
	Attempts    int
	LastError   *string
}
