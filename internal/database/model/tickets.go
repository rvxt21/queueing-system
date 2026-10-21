package model

import (
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

// TicketToDomain converts one DB row into a domain ticket.
func TicketToDomain(t model.Tickets) domain.Ticket {
	return domain.Ticket{
		ID:          t.ID,
		Number:      t.Number,
		ServiceCode: t.ServiceCode,
		Status:      domain.TicketStatus(t.Status),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
