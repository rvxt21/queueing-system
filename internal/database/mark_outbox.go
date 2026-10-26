package database

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
)

func (p *Postgres) MarkOutboxComplete(ctx context.Context, id uuid.UUID) error {
	stmt := table.Outbox.UPDATE().
		SET(table.Outbox.PublishedAt.SET(postgres.NOW())).
		WHERE(table.Outbox.ID.EQ(postgres.UUID(id)))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("%w: %w", ErrMarkOutbox, err)
	}

	return nil
}

func (p *Postgres) MarkOutboxFailed(ctx context.Context, id uuid.UUID, publishErr error) error {
	stmt := table.Outbox.UPDATE().
		SET(
			table.Outbox.Attempts.SET(table.Outbox.Attempts.ADD(postgres.Int(1))),
			table.Outbox.LastError.SET(postgres.String(publishErr.Error())),
		).
		WHERE(table.Outbox.ID.EQ(postgres.UUID(id)))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("%w: %w", ErrMarkOutbox, err)
	}

	return nil
}
