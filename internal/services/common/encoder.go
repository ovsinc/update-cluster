package common

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/builtin"
)

type Encoder nats.Encoder

func GetEncoder() Encoder {
	return &builtin.JsonEncoder{}
}
