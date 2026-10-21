package database

import (
	"context"
	"fmt"

	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	"github.com/rvxt21/queueing-system/internal/database/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (p *Postgres) GetServices(ctx context.Context) ([]domain.Service, error) {
	stmt := table.Services.
		SELECT(table.Services.AllColumns).
		FROM(table.Services).
		ORDER_BY(table.Services.Name.ASC())

	var services model.Services
	if err := stmt.QueryContext(ctx, p.db, &services); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetServices, err)
	}

	return services.ToDomain(), nil
}
