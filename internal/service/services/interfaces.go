package services

import (
	"context"

	"github.com/rvxt21/queueing-system/internal/domain"
)

type Repository interface {
	GetServices(ctx context.Context) ([]domain.Service, error)
	CreateTicket(
		ctx context.Context,
		serviceCode string,
		createFn func(service *domain.Service, counter *domain.TicketCounter) (*domain.Ticket, error),
	) (*domain.Ticket, error)
}
