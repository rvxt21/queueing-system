package services

import (
	"context"

	"github.com/rvxt21/queueing-system/internal/domain"
)

// Repository is the storage the services service depends on.
type Repository interface {
	GetServices(ctx context.Context) ([]domain.Service, error)
}
