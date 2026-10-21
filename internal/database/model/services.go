package model

import (
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type Services []model.Services

func (s Services) ToDomain() []domain.Service {
	res := make([]domain.Service, len(s))

	for i := range s {
		res[i] = ServiceToDomain(s[i])
	}

	return res
}

func ServiceToDomain(s model.Services) domain.Service {
	return domain.Service{
		ID:          s.ID,
		ServiceCode: s.ServiceCode,
		Name:        s.Name,
		Prefix:      s.Prefix,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
