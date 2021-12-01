package agent

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/multierr"
)

type PubSub struct {
	nc    *nats.Conn
	queue string
	subs  []*nats.Subscription
}

func NewService(nc *nats.Conn, queue string) *PubSub {
	return &PubSub{
		nc:    nc,
		queue: queue,
	}
}

func (ps *PubSub) Subscribes(handlers PubSubHandlers) error {
	ps.subs = make([]*nats.Subscription, 0, len(handlers))

	for subj, handle := range handlers {
		s, err := ps.nc.QueueSubscribe(subj, ps.queue, handle.Handle)
		if err != nil {
			return err
		}
		ps.subs = append(ps.subs, s)
	}

	return nil
}

func (ps *PubSub) Unsubscribes() error {
	var err error
	for _, sub := range ps.subs {
		err = multierr.Append(err, sub.Unsubscribe())
	}
	return err
}
