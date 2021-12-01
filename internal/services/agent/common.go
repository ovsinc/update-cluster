package agent

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

type IpubSubHandler interface {
	Handle(msg *nats.Msg)
}

// map:: channel : handler
type PubSubHandlers map[string]IpubSubHandler

func sendReply(enc common.Encoder, msg *nats.Msg, resp *common.HelloResponse) {
	reply, err := enc.Encode(msg.Reply, resp)
	if err != nil {
		log.Printf("Err while encode reply msg: %v", err)
	}

	if err := msg.Respond([]byte(reply)); err != nil {
		log.Printf("Err while send reply: %v", err)
	}
}
