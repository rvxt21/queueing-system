package model

import (
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

type Outboxes []model.Outbox

func (o Outboxes) ToDomain() []domain.Outbox {
	res := make([]domain.Outbox, len(o))
	for i := range o {
		res[i] = OutboxToDomain(o[i])
	}

	return res
}

func OutboxToDomain(s model.Outbox) domain.Outbox {
	return domain.Outbox{
		ID:          s.ID,
		AggregateID: s.AggregateID,
		RoutingKey:  s.RoutingKey,
		EventType:   s.EventType,
		Payload:     s.Payload,
		CreatedAt:   s.CreatedAt,
		PublishedAt: s.PublishedAt,
		Attempts:    int(s.Attempts),
		LastError:   s.LastError,
	}
}
