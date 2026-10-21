package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/model"
	"github.com/rvxt21/queueing-system/internal/database/gen/queueing/ticketing/table"
	dbmodel "github.com/rvxt21/queueing-system/internal/database/model"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (p *Postgres) CreateTicket(
	ctx context.Context,
	serviceCode string,
	createFn func(service *domain.Service, counter *domain.TicketCounter) (*domain.Ticket, error),
) (*domain.Ticket, error) {
	var created *domain.Ticket

	err := p.withTx(ctx, func(tx *sql.Tx) error {
		service, err := getServiceByCode(ctx, tx, serviceCode)
		if err != nil {
			return err
		}

		counter, err := lockTicketCounter(ctx, tx, serviceCode)
		if err != nil {
			return err
		}

		ticket, err := createFn(service, counter)
		if err != nil {
			return err
		}

		if ticket == nil {
			return ErrNoTicketToCreate
		}

		if err := saveTicketCounter(ctx, tx, counter); err != nil {
			return err
		}

		created, err = insertTicket(ctx, tx, ticket)

		return err
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func insertTicket(ctx context.Context, db qrm.Queryable, ticket *domain.Ticket) (*domain.Ticket, error) {
	stmt := table.Tickets.
		INSERT(table.Tickets.Number, table.Tickets.ServiceCode).
		VALUES(ticket.Number, ticket.ServiceCode).
		RETURNING(table.Tickets.AllColumns)

	var row model.Tickets
	if err := stmt.QueryContext(ctx, db, &row); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateTicket, err)
	}

	res := dbmodel.TicketToDomain(row)

	return &res, nil
}
