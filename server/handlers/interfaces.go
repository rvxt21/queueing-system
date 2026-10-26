package handlers

import (
	"context"

	"github.com/rvxt21/queueing-system/internal/domain"
)

type ServicesService interface {
	ListServices(ctx context.Context) ([]domain.Service, error)
	RequestService(ctx context.Context, serviceCode string) (*domain.Ticket, error)
}
