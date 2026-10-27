package database

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (p *Postgres) UpdateTicketStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error {
	stmt := table.Tickets.UPDATE().
		SET(
			table.Tickets.Status.SET(postgres.String(string(status))),
			table.Tickets.UpdatedAt.SET(postgres.NOW()),
		).
		WHERE(table.Tickets.ID.EQ(postgres.UUID(id)))

	res, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateTicketStatus, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateTicketStatus, err)
	}

	if rows == 0 {
		return domain.ErrTicketNotFound
	}

	return nil
}
