package rock

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
)

type Service struct {
	app   *application
	group *ServiceGroup

	name string
	path string
}

func (s *Service) Get(fn ...iris.Handler) *Service {
	return s.handle("GET", fn)
}
func (s *Service) Post(fn ...iris.Handler) *Service {
	return s.handle("POST", fn)
}
func (s *Service) Put(fn ...iris.Handler) *Service {
	return s.handle("PUT", fn)
}
func (s *Service) Option(fn ...iris.Handler) *Service {
	return s.handle("OPTION", fn)
}
func (s *Service) Delete(fn ...iris.Handler) *Service {
	return s.handle("DELETE", fn)
}

func (s *Service) handle(method string, fn []iris.Handler) *Service {
	path := s.path
	if s.group.party == nil {
		s.group.registerHandlerStatus(method, path)
		route := s.app.iris.Handle(method, path, fn...)
		route.MainHandlerName = s.name
	} else {
		path = s.group.path+path
		s.group.registerHandlerStatus(method, path)
		s.group.party.Handle(method, s.path, fn...)
	}
	defaultLogger.Infof("Service.handle %s %s %s", s.name, method, path)
	return s
}

type ServiceGroup struct {
	app   *application
	party iris.Party

	name string
	path string

	services      map[string]*Service
	handlerStatus map[string]bool
}

func (g *ServiceGroup) Use(mw ...iris.Handler) *ServiceGroup {
	g.party.Use(mw...)
	return g
}

func (g *ServiceGroup) NewService(name, path string) *Service {
	if name == "" {
		defaultLogger.Warn("Service should named with non-zero length")
	}
	if path == "" || path[0] != '/' {
		panic(fmt.Errorf("Service(%s) path(%s) should start with '/'", name, path))
	}
	if g.name != "" {
		name = g.name + "." + name
	}
	if _, ok := g.services[name]; ok {
		// TODO: warn name exist
	}
	s := &Service{app: g.app, group: g, path: path, name: name}
	if g.services == nil {
		g.services = map[string]*Service{}
	}
	g.services[name] = s
	return s
}

func (g *ServiceGroup) NewServiceGroup(name, path string) *ServiceGroup {
	if name == "" {
		defaultLogger.Warn("ServiceGroup should named with non-zero length")
	}
	if path == "" || path[0] != '/' {
		panic(fmt.Errorf("ServiceGroup(%s) path(%s) should start with '/'", name, path))
	}
	if g.name != "" {
		name = g.name + "." + name
	}
	newOne := &ServiceGroup{app: g.app, name: name, path: path}
	if g.party == nil {
		newOne.party = g.app.Iris().Party(path)
	} else {
		newOne.party = g.party.Party(path)
	}
	return newOne
}

func (g *ServiceGroup) registerHandlerStatus(method, path string) {
	key := method + path
	if g.handlerStatus[key] {
		panic(errors.New("duplicate handle " + method + " for " + path))
	}
	if g.handlerStatus == nil {
		g.handlerStatus = map[string]bool{}
	}
	g.handlerStatus[key] = true
}
