package services

import (
	"context"
	"fmt"

	"github.com/rvxt21/queueing-system/internal/domain"
)

func (s *Service) RequestService(ctx context.Context, serviceCode string) (*domain.Ticket, error) {
	ticket, err := s.repo.CreateTicket(ctx, serviceCode,
		func(service *domain.Service, counter *domain.TicketCounter) (*domain.Ticket, error) {
			if !service.IsActive {
				return nil, domain.ErrServiceInactive
			}

			return &domain.Ticket{
				Number:      domain.FormatTicketNumber(service.Prefix, counter.Next()),
				ServiceCode: service.ServiceCode,
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequestService, err)
	}

	return ticket, nil
}
