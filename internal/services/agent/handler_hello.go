package agent

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

type pubSubHello struct {
	enc common.Encoder
}

func NewPubSubHello(enc common.Encoder) IpubSubHandler {
	return &pubSubHello{
		enc: enc,
	}
}

func (ps *pubSubHello) Handle(msg *nats.Msg) {
	var message common.HelloRequest

	err := ps.enc.Decode(msg.Subject, msg.Data, &message)
	if err != nil {
		log.Printf("Err while decode response msg: %v", err)
		sendReply(ps.enc, msg, &common.HelloResponse{
			CommonResp: common.CommonResp{
				OK:     false,
				Errors: map[string]string{"common": err.Error()},
			},
		})
		return
	}

	log.Printf(
		"Got Request '%v' from '%s' channel",
		message,
		msg.Subject,
	)

	resp := common.HelloResponse{
		CommonResp: common.CommonResp{
			OK: true,
		},
		Msg: fmt.Sprintf("Answer to '%s' channel with message '%s'", msg.Subject, message.Msg),
	}

	log.Printf(
		"Send Reply '%v' to '%s' channel",
		resp, msg.Subject,
	)

	sendReply(ps.enc, msg, &resp)
}
