package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ovsinc/update-cluster/internal/services/agent"
	"github.com/ovsinc/update-cluster/internal/services/common"
	"go.uber.org/fx"
)

func newEncoder() common.Encoder {
	return common.GetEncoder()
}

func newHandlers(enc common.Encoder) agent.PubSubHandlers {
	return agent.PubSubHandlers{
		common.Config.HelloSubject: agent.NewPubSubHello(enc),
	}
}

//

func registryService(
	nc *nats.Conn, enc common.Encoder, handlers agent.PubSubHandlers,
) (*agent.PubSub, error) {
	svc := agent.NewService(nc, common.Config.QueueGroup)
	return svc, svc.Subscribes(handlers)
}

func run(lifecycle fx.Lifecycle, svc *agent.PubSub) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				return nil
			},

			OnStop: func(c context.Context) error {
				log.Println("[BACKEND] stops server...")

				_ = svc.Unsubscribes()

				time.Sleep(common.Config.BackendShutdownDur)

				log.Println("[BACKEND] server stopped")
				return nil
			},
		},
	)
	return nil
}

//

func die() { //nolint:deadcode,unused
	log.Println("[BACKEND] Simulates fail")
	os.Exit(255)
}

func main() {
	// die()

	appCtx := fx.New(
		// options
		fx.StartTimeout(common.Config.GracefulStartTimeout),
		fx.StopTimeout(common.Config.GracefulStopTimeout),

		fx.Provide(
			common.ConfigNats,
			common.ConnectNats,
			newEncoder,
			newHandlers,
			registryService,
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
