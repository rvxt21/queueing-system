package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type RequestServiceRequest struct {
	ServiceCode string `binding:"required" example:"passport_new" json:"service_code"`
}

type TicketResponse struct {
	ID          uuid.UUID `json:"id"`
	Number      string    `example:"PN-042"       json:"number"`
	ServiceCode string    `example:"passport_new" json:"service_code"`
	Status      string    `example:"new"          json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTicketResponse(t domain.Ticket) TicketResponse {
	return TicketResponse{
		ID:          t.ID,
		Number:      t.Number,
		ServiceCode: t.ServiceCode,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
	}
}
