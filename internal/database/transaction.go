package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (p *Postgres) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBeginTx, err)
	}

	defer func() { _ = tx.Rollback() }()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %w", ErrCommitTx, err)
	}

	return nil
}
