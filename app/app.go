package app

import (
	"github.com/rvxt21/queueing-system/config"
	"github.com/rvxt21/queueing-system/internal/broker"
	"github.com/rvxt21/queueing-system/internal/database"
	"github.com/rvxt21/queueing-system/internal/service/hub"
	"github.com/rvxt21/queueing-system/internal/service/outbox"
	"github.com/rvxt21/queueing-system/internal/service/services"
	"github.com/rvxt21/queueing-system/internal/service/tickets"
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
				fx.As(new(outbox.Repository)),
				fx.As(new(tickets.Repository)),
			),
			fx.Annotate(
				broker.NewRabbitMQ,
				fx.As(fx.Self()),
				fx.As(new(outbox.Broker)),
				fx.As(new(tickets.Broker)),
			),
			fx.Annotate(
				services.NewService,
				fx.As(new(handlers.ServicesService)),
			),
			fx.Annotate(
				hub.NewHub,
				fx.As(fx.Self()),
				fx.As(new(tickets.Notifier)),
			),
			handlers.NewHandler,
			server.NewServer,
			logger.NewLogger,
			outbox.NewWorker,
			tickets.NewWorker,
		),

		fx.Invoke(invokeHooks),
	)
}

type hooks struct {
	fx.In

	Postgres *database.Postgres
	RabbitMQ *broker.RabbitMQ
	Server   *server.Server
	Outbox   *outbox.Worker
	Tickets  *tickets.Worker
}

func invokeHooks(lc fx.Lifecycle, h hooks) {
	lc.Append(fx.Hook{OnStart: h.Postgres.Start, OnStop: h.Postgres.Stop})
	lc.Append(fx.Hook{OnStart: h.RabbitMQ.Start, OnStop: h.RabbitMQ.Stop})
	lc.Append(fx.Hook{OnStart: h.Server.Start, OnStop: h.Server.Stop})
	lc.Append(fx.Hook{OnStart: h.Outbox.Start, OnStop: h.Outbox.Stop})
	lc.Append(fx.Hook{OnStart: h.Tickets.Start, OnStop: h.Tickets.Stop})
}
