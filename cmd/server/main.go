package main

import (
	"context"

	"github.com/m2tx/gofxtest/domain/event"
	"github.com/m2tx/gofxtest/internal/env"
	"github.com/m2tx/gofxtest/internal/http"
	"github.com/m2tx/gofxtest/internal/queue"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type AppConfig struct {
	ServiceName string `env:"NAME" default:"gofxtest"`
}

func main() {
	fx.New(
		fx.Provide(context.Background),
		fx.Provide(env.New[AppConfig]),
		fx.Provide(env.New[http.HttpConfig]),
		fx.Provide(env.New[mongo.MongoConfig]),
		fx.Provide(env.New[queue.RabbitMQConfig]),
		fx.Provide(AsQueue(queue.NewRabbitMQ)),
		fx.Provide(mongo.NewClient),
		fx.Provide(mongo.NewEventRepository),
		fx.Provide(event.NewEventService),
		fx.Provide(http.NewServer, fx.Annotate(
			http.NewHandler,
			fx.ParamTags(`group:"routeHandlers"`),
		)),
		fx.Provide(AsRouteHandler(http.NewHealthcheckRoute)),
		fx.Provide(AsRouteHandler(http.NewSwaggerRoute)),
		fx.Provide(AsRouteHandler(http.NewEventRoute)),
		fx.Provide(func(config AppConfig) (*zap.Logger, error) {
			logger, err := zap.NewProduction()
			if err != nil {
				return nil, err
			}

			return logger.Named(config.ServiceName), nil
		}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: log,
			}
		}),
		fx.Invoke(func(lc fx.Lifecycle, client mongo.MongoClient) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return client.Disconnect(ctx)
				},
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, srv http.HttpServer) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return srv.Start()
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, srv queue.Subscriber, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := srv.Subscribe(ctx, event.EventCreatedTopic, func(msg queue.Message) {
						log.Info("consume event-created", zap.String("data", string(msg.Data)))
					})
					if err != nil {
						return err
					}
					err = srv.Subscribe(ctx, event.EventUpdatedTopic, func(msg queue.Message) {
						log.Info("consume event-updated", zap.String("data", string(msg.Data)))
					})
					if err != nil {
						return err
					}
					err = srv.Subscribe(ctx, event.EventDeletedTopic, func(msg queue.Message) {
						log.Info("consume event-deleted", zap.String("data", string(msg.Data)))
					})
					if err != nil {
						return err
					}
					return nil
				},
			})
		}),
	).Run()
}

func AsRouteHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(http.RouteHandler)),
		fx.ResultTags(`group:"routeHandlers"`),
	)
}

func AsQueue(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(queue.Publisher)),
		fx.As(new(queue.Subscriber)),
	)
}
