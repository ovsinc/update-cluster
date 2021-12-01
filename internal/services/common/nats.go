package common

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

func ConnectNats(lc fx.Lifecycle) (*nats.Conn, error) {
	nc, err := nats.Connect(
		Config.NatsURL,

		nats.Name("Agent Responder"),

		// nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(3),
		nats.ReconnectWait(1*time.Second),
		nats.FlusherTimeout(1*time.Second),

		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Got disconnected! Reason: %q\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("Connection closed. Reason: %q\n", nc.LastError())
		}),
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				if nc != nil {
					nc.Close()
				}
				return nil
			},
		},
	)

	return nc, err
}
