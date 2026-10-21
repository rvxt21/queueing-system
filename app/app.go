package app

import (
	"github.com/rvxt21/queueing-system/config"
	"github.com/rvxt21/queueing-system/internal/broker"
	"github.com/rvxt21/queueing-system/internal/database"
	"github.com/rvxt21/queueing-system/internal/service/services"
	"github.com/rvxt21/queueing-system/pkg/logger"
	"github.com/rvxt21/queueing-system/server"
	"github.com/rvxt21/queueing-system/server/handlers"
	"go.uber.org/fx"
)

func App(cfg *config.Config) *fx.App {
	return fx.New(
		fx.Supply(
			cfg,
		),

		fx.Provide(
			fx.Annotate(
				database.NewPostgres,
				fx.As(fx.Self()),
				fx.As(new(services.Repository)),
			),
			broker.NewRabbitMQ,
			fx.Annotate(
				services.NewService,
				fx.As(new(handlers.ServicesService)),
			),
			handlers.NewHandler,
			server.NewServer,
			logger.NewLogger,
		),

		fx.Invoke(invokeHooks),
	)
}

type hooks struct {
	fx.In

	Postgres *database.Postgres
	RabbitMQ *broker.RabbitMQ
	Server   *server.Server
}

func invokeHooks(lc fx.Lifecycle, h hooks) {
	lc.Append(fx.Hook{OnStart: h.Postgres.Start, OnStop: h.Postgres.Stop})
	lc.Append(fx.Hook{OnStart: h.RabbitMQ.Start, OnStop: h.RabbitMQ.Stop})
	lc.Append(fx.Hook{OnStart: h.Server.Start, OnStop: h.Server.Stop})
}
