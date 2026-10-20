package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"

	_ "github.com/lib/pq" // postgres driver
	"github.com/rvxt21/queueing-system/config"
)

type Postgres struct {
	config config.Postgres
	db     *sql.DB
}

func (p *Postgres) Start(ctx context.Context) error {
	dsn := (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.config.User, p.config.Password),
		Host:     net.JoinHostPort(p.config.Host, p.config.Port),
		Path:     p.config.Database,
		RawQuery: url.Values{"sslmode": {p.config.SSLMode}}.Encode(),
	}).String()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOpenPostgresConnection, err)
	}

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrPingingDatabase, errors.Join(err, db.Close()))
	}

	p.db = db

	return nil
}

func (p *Postgres) Stop(ctx context.Context) error {
	if p.db != nil {
		return p.db.Close()
	}

	return nil
}

func NewPostgres(cfg *config.Config) *Postgres {
	return &Postgres{
		config: cfg.Postgres,
	}
}
