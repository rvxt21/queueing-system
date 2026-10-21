package broker

import (
	"time"

	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type TicketMessage struct {
	ID          uuid.UUID `json:"id"`
	Number      string    `json:"number"`
	ServiceCode string    `json:"service_code"`
	CreatedAt   time.Time `json:"created_at"`
}

func newTicketMessage(t *domain.Ticket) TicketMessage {
	return TicketMessage{
		ID:          t.ID,
		Number:      t.Number,
		ServiceCode: t.ServiceCode,
		CreatedAt:   t.CreatedAt,
	}
}
