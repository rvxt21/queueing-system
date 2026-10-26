package database

import (
	"context"
	"fmt"

	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	"github.com/rvxt21/queueing-system/internal/database/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (p *Postgres) GetOutboxes(ctx context.Context, limit int64) ([]domain.Outbox, error) {
	stmt := table.Outbox.
		SELECT(
			table.Outbox.AllColumns,
		).
		FROM(table.Outbox).
		WHERE(table.Outbox.PublishedAt.IS_NULL()).
		ORDER_BY(table.Outbox.CreatedAt.ASC()).
		LIMIT(limit)

	var outboxes model.Outboxes

	if err := stmt.QueryContext(ctx, p.db, &outboxes); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetOutboxes, err)
	}

	return outboxes.ToDomain(), nil
}
