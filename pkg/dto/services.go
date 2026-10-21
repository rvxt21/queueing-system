package dto

import (
	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type ServicesResponse struct {
	ID          uuid.UUID `example:"01a15d53-a46f-762b-bebd-4232ae144e8a" format:"uuid"       json:"id"`
	ServiceCode string    `example:"passport_new"                         json:"service_code"`
	Name        string    `example:"Apply for a new passport"             json:"name"`
	IsActive    bool      `example:"true"                                 json:"is_active"`
}

func NewServicesResponse(services []domain.Service) []ServicesResponse {
	res := make([]ServicesResponse, len(services))

	for i := range services {
		res[i] = ServicesResponse{
			ID:          services[i].ID,
			ServiceCode: services[i].ServiceCode,
			Name:        services[i].Name,
			IsActive:    services[i].IsActive,
		}
	}

	return res
}
