package app

import (
	"github.com/rvxt21/queueing-system/config"
	"github.com/rvxt21/queueing-system/internal/broker"
	"github.com/rvxt21/queueing-system/internal/database"
	"go.uber.org/fx"
)

func App(cfg *config.Config) *fx.App {
	return fx.New(
		fx.Supply(
			cfg,
		),

		fx.Provide(
			database.NewPostgres,
			broker.NewRabbitMQ,
		),

		fx.Invoke(invokeHooks),
	)
}

type hooks struct {
	fx.In

	Postgres *database.Postgres
	RabbitMQ *broker.RabbitMQ
}

func invokeHooks(lc fx.Lifecycle, h hooks) {
	lc.Append(fx.Hook{OnStart: h.Postgres.Start, OnStop: h.Postgres.Stop})
	lc.Append(fx.Hook{OnStart: h.RabbitMQ.Start, OnStop: h.RabbitMQ.Stop})
}
