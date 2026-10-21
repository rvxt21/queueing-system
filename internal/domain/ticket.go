package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusNew       TicketStatus = "new"
	TicketStatusSent      TicketStatus = "sent"
	TicketStatusDelivered TicketStatus = "delivered"
	TicketStatusResolved  TicketStatus = "resolved"
)

type Ticket struct {
	ID          uuid.UUID
	Number      string
	ServiceCode string
	Status      TicketStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func FormatTicketNumber(prefix string, n int32) string {
	return fmt.Sprintf("%s-%03d", prefix, n)
}
