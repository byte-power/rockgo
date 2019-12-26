package rock

import (
	"errors"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

type Service struct {
	app   *application
	group *ServiceGroup

	name string
	path string
}

func (s *Service) Get(fn ...iris.Handler) *Service {
	return s.handle(http.MethodGet, fn)
}
func (s *Service) Post(fn ...iris.Handler) *Service {
	return s.handle(http.MethodPost, fn)
}
func (s *Service) Put(fn ...iris.Handler) *Service {
	return s.handle(http.MethodPut, fn)
}
func (s *Service) Connect(fn ...iris.Handler) *Service {
	return s.handle(http.MethodConnect, fn)
}
func (s *Service) Head(fn ...iris.Handler) *Service {
	return s.handle(http.MethodHead, fn)
}
func (s *Service) Option(fn ...iris.Handler) *Service {
	return s.handle(http.MethodOptions, fn)
}
func (s *Service) Patch(fn ...iris.Handler) *Service {
	return s.handle(http.MethodPatch, fn)
}
func (s *Service) Trace(fn ...iris.Handler) *Service {
	return s.handle(http.MethodTrace, fn)
}
func (s *Service) Delete(fn ...iris.Handler) *Service {
	return s.handle(http.MethodDelete, fn)
}

func (s *Service) handle(method string, fn []iris.Handler) *Service {
	if s.group.party == nil {
		s.group.registerHandlerStatus(method, s.path)
		route := s.app.iris.Handle(method, s.path, fn...)
		route.MainHandlerName = s.name
	} else {
		s.group.registerHandlerStatus(method, s.group.path+s.path)
		s.group.party.Handle(method, s.path, fn...)
	}
	return s
}

// HandleDir would serve static files.
// - Parameters:
//   - dir: Directory path
func (s *Service) HandleDir(dir string, opts ...router.DirOptions) {
	if s.group.party == nil {
		s.group.registerHandlerStatus("DIR", s.path)
		s.app.iris.HandleDir(s.path, dir, opts...)
	} else {
		s.group.registerHandlerStatus("DIR", s.group.path+s.path)
		s.group.party.HandleDir(s.path, dir, opts...)
	}
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
		defaultLogger.Warn("Service should named with non-empty string")
	}
	if g.name != "" {
		name = g.name + "." + name
	}
	if _, ok := g.services[name]; ok {
		defaultLogger.Warn("Service name duplicated", "name", name)
	}
	if path == "/" || path[0] != '/' {
		path = "/" + path
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
		defaultLogger.Warn("ServiceGroup should named with non-empty string")
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
