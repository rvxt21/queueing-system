package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func insertTicketCreatedEvent(ctx context.Context, db qrm.Executable, ticket *domain.Ticket) error {
	payload, err := json.Marshal(domain.TicketCreatedEvent{
		ID:          ticket.ID,
		Number:      ticket.Number,
		ServiceCode: ticket.ServiceCode,
		CreatedAt:   ticket.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("%w: marshal payload: %w", ErrCreateOutbox, err)
	}

	stmt := table.Outbox.INSERT(
		table.Outbox.AggregateID,
		table.Outbox.EventType,
		table.Outbox.RoutingKey,
		table.Outbox.Payload,
	).VALUES(
		ticket.ID,
		domain.EventTicketCreated,
		ticket.ServiceCode,
		string(payload),
	)

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return fmt.Errorf("%w: %w", ErrCreateOutbox, err)
	}

	return nil
}
