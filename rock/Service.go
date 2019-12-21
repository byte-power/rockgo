package rock

import (
	"errors"

	"github.com/kataras/iris/v12"
)

type Service struct {
	app  *application
	name string
	path string
}

func (s *Service) Get(fn func(iris.Context)) *Service {
	return s.handle("GET", fn)
}
func (s *Service) Post(fn func(iris.Context)) *Service {
	return s.handle("POST", fn)
}
func (s *Service) Put(fn func(iris.Context)) *Service {
	return s.handle("PUT", fn)
}
func (s *Service) Option(fn func(iris.Context)) *Service {
	return s.handle("OPTION", fn)
}
func (s *Service) Delete(fn func(iris.Context)) *Service {
	return s.handle("DELETE", fn)
}

func (s *Service) handle(method string, fn func(iris.Context)) *Service {
	key := method + s.path
	if s.app.handlerStatus[key] {
		panic(errors.New("duplicate handle " + method + " for " + s.path))
	}
	s.app.handlerStatus[key] = true
	defaultLogger.Infof("Service.handle %s %s %s", s.name, method, s.path)
	route := s.app.iris.Handle(method, s.path, fn)
	route.MainHandlerName = s.name
	return s
}
