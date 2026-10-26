package domain

import (
	"time"

	"github.com/google/uuid"
)

const EventTicketCreated = "ticket.created"

type TicketCreatedEvent struct {
	ID          uuid.UUID `json:"id"`
	Number      string    `json:"number"`
	ServiceCode string    `json:"service_code"`
	CreatedAt   time.Time `json:"created_at"`
}
