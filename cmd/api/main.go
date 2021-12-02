package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"

	"github.com/ovsinc/update-cluster/internal/services/api"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

const (
	_heath_path = "/health"
	_api_path   = "/hello"
	_whoami     = "/"
)

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Test API 1.0",
		AppName:       "Test API 1.0",
		//
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})
	app.Use(logger.New())
	return app
}

func newEncoder() common.Encoder {
	return common.GetEncoder()
}

//

func run(lifecycle fx.Lifecycle, app *fiber.App, nc *nats.Conn, enc common.Encoder) error {
	svc := api.NewService(app)

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				svc.Registry(api.Routers{
					_heath_path: api.Route{
						Method:  http.MethodGet,
						Handler: api.NewHealthHandler(nc, common.Config.HelloSubject, enc),
					},
					_api_path: api.Route{
						Method:  http.MethodPost,
						Handler: api.NewHelloHandler(nc, common.Config.HelloSubject, enc),
					},
					_whoami: api.Route{
						Method:  http.MethodGet,
						Handler: api.NewWhoamiHandler(),
					},
				})

				go func() {
					_ = app.Listen(fmt.Sprintf(":%d", common.Config.Port))
				}()

				log.Println("[API] server started")

				return nil
			},

			OnStop: func(context.Context) error {
				log.Println("[API] stops server...")

				svc.Ungregistry()
				err := app.Shutdown()

				time.Sleep(common.Config.APIShutdownDur)

				log.Println("[API] server stopped")
				return err
			},
		},
	)
	return nil
}

func main() {
	appCtx := fx.New(
		// options
		fx.StartTimeout(common.Config.GracefulStartTimeout),
		fx.StopTimeout(common.Config.GracefulStopTimeout),

		fx.Provide(
			newApp,
			newEncoder,
			common.ConfigNats,
			common.ConnectNats,
		),

		fx.Invoke(
			run,
		),
	)

	appCtx.Run()

	if err := appCtx.Err(); err != nil {
		log.Println(err)
	}
}
