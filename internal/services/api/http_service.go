package api

import (
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	Method  string
	Handler Handler
}

type Routers map[string]Route

//

//

type Service struct {
	handlers Routers
	app      *fiber.App
}

func NewService(app *fiber.App) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Registry(handlers Routers) {
	s.handlers = handlers
	for path, handhandler := range handlers {
		s.app.Add(handhandler.Method, path, handhandler.Handler.Handle)
	}
}

func (s *Service) Ungregistry() {
	for _, handhandler := range s.handlers {
		if stop, ok := handhandler.Handler.(StopHandler); ok {
			stop.Stop()
		}
	}
}
