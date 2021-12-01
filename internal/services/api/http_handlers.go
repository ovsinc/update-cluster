package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/ovsinc/update-cluster/internal/services/common"
)

const timeout = 2 * time.Second

type Handler interface {
	Handle(c *fiber.Ctx) error
}

type StopHandler interface {
	Stop()
}

//

type healthHandler struct {
	nc     *nats.Conn
	subj   string
	enc    common.Encoder
	status string
	code   int
	mu     *sync.Mutex
	done   chan struct{}
}

func NewHealthHandler(
	nc *nats.Conn,
	subj string,
	enc common.Encoder,
) Handler {
	h := &healthHandler{
		nc:     nc,
		subj:   subj,
		enc:    enc,
		code:   StatusUNKNOWN,
		status: "unknown",
		mu:     new(sync.Mutex),
		done:   make(chan struct{}),
	}
	h.checker()
	return h
}

const (
	StatusOK               = 200
	StatusNotResponse      = 600
	StatusBadRequest       = 700
	StatusUNKNOWN          = 0
	StatusSideCarUnhealthy = 900
	StatusExit             = 800
)

func (h *healthHandler) checker() {
	req := common.HelloRequest{
		Msg: "hello",
	}
	payload, _ := h.enc.Encode(h.subj, &req)
	timeout := time.Second

	go func() {
		setStatus := func(code int, status string) {
			h.mu.Lock()
			h.code = code
			h.status = status
			h.mu.Unlock()

			select {
			case <-h.done:
			default:
				time.Sleep(2 * time.Second)
			}
		}

		for {
			msg, err := h.nc.Request(h.subj, payload, timeout)

			select {
			case <-h.done:

			default:
				switch {
				case err != nil:
					setStatus(StatusNotResponse, err.Error())

				case h.nc.LastError() != nil:
					setStatus(StatusNotResponse, h.nc.LastError().Error())

				default:
					var resp common.HelloResponse
					err := h.enc.Decode(h.subj, msg.Data, &resp)

					switch {
					case err != nil:
						setStatus(StatusBadRequest, err.Error())

					case !resp.OK:
						errs := make([]string, 0, len(resp.Errors))
						for k, e := range resp.Errors {
							errs = append(errs, k+": "+e)
						}
						setStatus(StatusSideCarUnhealthy, strings.Join(errs, "; "))

					default:
						setStatus(StatusOK, "OK")
					}
				}
			}

		}
	}()
}

func (h *healthHandler) Handle(c *fiber.Ctx) error {
	return c.Status(h.code).SendString(h.status)
}

func (h *healthHandler) Stop() {
	close(h.done)
}

//

type helloHandler struct {
	nc   *nats.Conn
	subj string
	enc  common.Encoder
}

func NewHelloHandler(
	nc *nats.Conn,
	subj string,
	enc common.Encoder,
) Handler {
	return &helloHandler{
		nc:   nc,
		subj: subj,
		enc:  enc,
	}
}

func (h *helloHandler) Handle(c *fiber.Ctx) error {
	req := common.HelloRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(
				common.HelloResponse{
					CommonResp: common.CommonResp{
						OK:     false,
						Errors: map[string]string{"common": err.Error()},
					},
				},
			)
	}

	req = req.Sanitize()

	payload, err := h.enc.Encode(h.subj, &req)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(
				common.HelloResponse{
					CommonResp: common.CommonResp{
						OK:     false,
						Errors: map[string]string{"common": err.Error()},
					},
				},
			)
	}

	msg, err := h.nc.Request(h.subj, payload, timeout)
	if err != nil {
		if h.nc.LastError() != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(
					common.HelloResponse{
						CommonResp: common.CommonResp{
							OK:     false,
							Errors: map[string]string{"common": h.nc.LastError().Error()},
						},
					},
				)
		}

		return c.Status(http.StatusInternalServerError).
			JSON(
				common.HelloResponse{
					CommonResp: common.CommonResp{
						OK:     false,
						Errors: map[string]string{"common": err.Error()},
					},
				},
			)
	}

	var resp common.HelloResponse

	if err := h.enc.Decode(h.subj, msg.Data, &resp); err != nil {
		log.Printf("Err. Response decode error: '%v'", err)
		return c.Status(http.StatusInternalServerError).
			JSON(
				common.HelloResponse{
					CommonResp: common.CommonResp{
						OK:     false,
						Errors: map[string]string{"common": err.Error()},
					},
				},
			)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

//

type whoamiHandler struct{}

func NewWhoamiHandler() Handler {
	return &whoamiHandler{}
}

func (h *whoamiHandler) Handle(c *fiber.Ctx) error {
	hostname, _ := os.Hostname()

	headers := make(map[string]string)

	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	data := struct {
		Hostname  string   `json:"hostname,omitempty"`
		ClientIPs []string `json:"client_ips,omitempty"`
		HostIPs   []string `json:"host_ips,omitempty"`
		Headers   string   `json:"headers,omitempty"`
		URL       string   `json:"url,omitempty"`
		Host      string   `json:"host,omitempty"`
		Method    string   `json:"method,omitempty"`
		Envs      []string `json:"envs,omitempty"`
	}{
		Hostname:  hostname,
		ClientIPs: c.IPs(),
		Headers:   c.Request().Header.String(),
		URL:       c.OriginalURL(),
		Host:      c.Hostname(),
		Method:    string(c.Request().Header.Method()),
		Envs:      os.Environ(),
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				data.HostIPs = append(
					data.HostIPs,
					fmt.Sprintf("%s:%s", i.Name, ip.String()),
				)
			}
		}
	}

	return c.
		Status(http.StatusOK).
		JSON(data)
}
