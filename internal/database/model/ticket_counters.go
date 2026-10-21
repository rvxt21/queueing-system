package model

import (
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func TicketCounterToDomain(c model.TicketCounters) domain.TicketCounter {
	return domain.TicketCounter{
		ServiceCode: c.ServiceCode,
		Day:         c.Day,
		LastNumber:  c.LastNumber,
	}
}
