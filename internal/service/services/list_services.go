package services

import (
	"context"
	"fmt"

	"github.com/rvxt21/queueing-system/internal/domain"
)

func (s *Service) ListServices(ctx context.Context) ([]domain.Service, error) {
	services, err := s.repo.GetServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListServices, err)
	}

	return services, nil
}
