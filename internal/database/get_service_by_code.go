package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	dbmodel "github.com/rvxt21/queueing-system/internal/database/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (p *Postgres) GetServiceByCode(ctx context.Context, serviceCode string) (*domain.Service, error) {
	return getServiceByCode(ctx, p.db, serviceCode)
}

func getServiceByCode(ctx context.Context, db qrm.Queryable, serviceCode string) (*domain.Service, error) {
	stmt := table.Services.
		SELECT(table.Services.AllColumns).
		WHERE(table.Services.ServiceCode.EQ(postgres.String(serviceCode)))

	var service model.Services
	if err := stmt.QueryContext(ctx, db, &service); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, domain.ErrServiceNotFound
		}

		return nil, fmt.Errorf("%w: %w", ErrGetServiceByCode, err)
	}

	res := dbmodel.ServiceToDomain(service)

	return &res, nil
}
