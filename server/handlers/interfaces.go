package handlers

import (
	"context"

	"github.com/rvxt21/queueing-system/internal/domain"
	"github.com/rvxt21/queueing-system/internal/service/hub"
)

type ServicesService interface {
	ListServices(ctx context.Context) ([]domain.Service, error)
	RequestService(ctx context.Context, serviceCode string) (*domain.Ticket, error)
}

type Hub interface {
	Subscribe() *hub.Subscriber
	Unsubscribe(sub *hub.Subscriber)
}
