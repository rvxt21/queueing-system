package database

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	dbmodel "github.com/rvxt21/queueing-system/internal/database/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func lockTicketCounter(ctx context.Context, db qrm.Queryable, serviceCode string) (*domain.TicketCounter, error) {
	counters := table.TicketCounters

	stmt := counters.
		INSERT(counters.ServiceCode, counters.Day, counters.LastNumber).
		VALUES(serviceCode, postgres.CURRENT_DATE(), 0).
		ON_CONFLICT(counters.ServiceCode, counters.Day).
		DO_UPDATE(postgres.SET(counters.LastNumber.SET(counters.LastNumber))).
		RETURNING(counters.AllColumns)

	var row model.TicketCounters
	if err := stmt.QueryContext(ctx, db, &row); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLockTicketCounter, err)
	}

	res := dbmodel.TicketCounterToDomain(row)

	return &res, nil
}

func saveTicketCounter(ctx context.Context, db qrm.Executable, counter *domain.TicketCounter) error {
	counters := table.TicketCounters

	stmt := counters.
		UPDATE(counters.LastNumber).
		SET(counter.LastNumber).
		WHERE(
			counters.ServiceCode.EQ(postgres.String(counter.ServiceCode)).
				AND(counters.Day.EQ(postgres.DateT(counter.Day))),
		)

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveTicketCounter, err)
	}

	return nil
}
