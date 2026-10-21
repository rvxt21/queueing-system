package model

import (
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type Services []model.Services

func (s Services) ToDomain() []domain.Service {
	res := make([]domain.Service, len(s))

	for i := range s {
		res[i] = domain.Service{
			ID:          s[i].ID,
			ServiceCode: s[i].ServiceCode,
			Name:        s[i].Name,
			IsActive:    s[i].IsActive,
			CreatedAt:   s[i].CreatedAt,
			UpdatedAt:   s[i].UpdatedAt,
		}
	}

	return res
}
