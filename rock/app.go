package rock

import (
	"fmt"

	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
)

const ConfigFilename = "rockgo"

type PanicHandler func(ctx iris.Context, err error)

type Application interface {
	// Name returns app_name in rockgo config.
	Name() string

	// Iris returns iris application.
	Iris() *iris.Application

	// Set panic handler, only work on sentry.repanic is true.
	SetPanicHandler(fn PanicHandler)

	// NewService make Service to register handler.
	//
	// - Parameters:
	//   - name: for statsd
	//   - path: path from root, e.g. "foo" for "/foo"
	NewService(name, path string) *Service

	// NewServiceGroup make ServiceGroup to handle multiple Services.
	//
	// - Parameters:
	//   - name: for statsd
	//   - path: directory name in path from root
	NewServiceGroup(name, path string) *ServiceGroup

	// Run server with <host> and multiple [conf].
	//
	// - Parameters:
	//   - host: [server name] with port. e.g. "mydomain.com:80" or ":8080" (equal "0.0.0.0:8080")
	Run(host string, conf ...host.Configurator)
}

// NewApplication would load config files from <configDir>, and then make Application with rockgo.yaml (or json) in <configDir>.
func NewApplication(configDir string) (Application, error) {
	err := ImportConfigFilesFromDirectory(configDir)
	if err != nil {
		return nil, util.NewError(ErrNameApplicationInitFailure, err)
	}
	cfg := util.AnyToAnyMap(sharedConfig[ConfigFilename])
	if cfg == nil {
		return nil, fmt.Errorf("%s %s/%s.yaml (or json) not exists", ErrNameApplicationInitFailure, configDir, ConfigFilename)
	}
	a := &application{iris: iris.New()}
	a.rootGroup = ServiceGroup{app: a}
	err = a.init(cfg)
	if err != nil {
		return nil, util.NewError(ErrNameApplicationInitFailure, err)
	}
	a.iris.Use(newRockMiddleware(a))
	return a, nil
}

type application struct {
	name      string
	iris      *iris.Application
	rootGroup ServiceGroup

	panicHandler PanicHandler
}

func (a *application) Name() string {
	return a.name
}

func (a *application) init(cfg util.AnyMap) error {
	appName := util.AnyToString(cfg["app_name"])
	// create loggers
	if cfgIt := util.AnyToAnyMap(cfg["log"]); cfgIt != nil {
		for name, v := range cfgIt {
			vs := util.AnyToAnyMap(v)
			if vs == nil {
				continue
			}
			logger, err := parseLogger(appName, name, vs)
			if err != nil {
				return err
			}
			loggers[name] = logger
		}
	}
	if it := loggers["default"]; it != nil {
		defaultLogger = *it
	}
	if cfgIt := util.AnyToAnyMap(cfg["metric"]); cfgIt != nil {
		initMetric(appName, cfgIt)
	}
	if cfgIt := util.AnyToAnyMap(cfg["sentry"]); cfgIt != nil {
		sentryMW, err := newSentryMiddleware(parseSentryOption(cfgIt))
		if err != nil {
			return err
		}
		a.iris.Use(sentryMW)
	}
	if appName == "" {
		defaultLogger.Warnf("'app_name' was not defined or it is empty in config (statsd and fluent depended on it).")
	}
	a.name = appName
	return nil
}

func (a *application) NewService(name, path string) *Service {
	return a.rootGroup.NewService(name, path)
}

func (a *application) NewServiceGroup(name, path string) *ServiceGroup {
	return a.rootGroup.NewServiceGroup(name, path)
}

func (a *application) Iris() *iris.Application {
	return a.iris
}

func (a *application) SetPanicHandler(fn PanicHandler) {
	a.panicHandler = fn
}

func (a *application) Run(host string, conf ...host.Configurator) {
	a.iris.Run(iris.Addr(host, conf...))
}
