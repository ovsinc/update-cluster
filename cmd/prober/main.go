package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

const (
	name    = "prober"
	timeout = time.Second
)

func say(msg string) {
	log.Printf("[%s] %s", strings.ToUpper(name), msg)
}

func die(msg string, code int) {
	say(msg)
	os.Exit(code)
}

func main() {
	subj := common.Config.HelloSubject
	opts := common.ConfigNats()
	enc := common.GetEncoder()
	req := common.HelloRequest{Msg: "hello"}
	payload, _ := enc.Encode(subj, &req)
	var resp common.HelloResponse

	nc, err := nats.Connect(common.Config.NatsURL, opts...)
	if err != nil {
		die(fmt.Sprintf("connect error: %s", err.Error()), 1)
	}

	msg, err := nc.Request(subj, payload, timeout)
	if err != nil {
		die(fmt.Sprintf("send request error: %s", err.Error()), 1)
	}

	if err := enc.Decode(subj, msg.Data, &resp); err != nil {
		die(fmt.Sprintf("response decode error: '%v'", err), 1)
	}

	if !resp.OK {
		str := make([]string, 0, len(resp.Errors))
		for k, v := range resp.Errors {
			str = append(str, fmt.Sprintf("%s:%s", k, v))
		}
		die(fmt.Sprintf("error in response: %s", strings.Join(str, "; ")), 1)
	}

	die("ok", 0)
}
