package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"

	"github.com/ovsinc/update-cluster/internal/services/api"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

const (
	_portEnv    = "API_PORT"
	_port       = 8000
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
	port, err := strconv.Atoi(os.Getenv(_portEnv))
	if err != nil {
		log.Printf("port error: %v", err)
		log.Printf("use default port: %d", _port)
		port = _port
	}
	svc := api.NewService(app)

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Printf("HTTP server listen :%d", port)

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
					_ = app.Listen(fmt.Sprintf(":%d", port))
				}()
				return nil
			},
			OnStop: func(context.Context) error {
				log.Printf("Stop server on :%d", port)
				svc.Ungregistry()
				time.Sleep(common.Config.APIShutdownDur)
				return app.Shutdown()
			},
		},
	)
	return nil
}

var port int

func main() {
	appCtx := fx.New(
		// options
		fx.StartTimeout(common.Config.GracefulStartTimeout),
		fx.StopTimeout(common.Config.GracefulStopTimeout),

		fx.Provide(
			newApp,
			newEncoder,
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
