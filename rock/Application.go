package rock

import (
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
)

type Application interface {
	// 载入config_rock，以初始化各个内部模块，并自动添加用于异常恢复的recover和sentry中间件、记录基础metric信息的中间件
	//
	// InitWithConfig would load config from <path>, and then init each internal components.
	//
	// config file at <path> should have extension 'json' for JSON format, or 'yaml' / 'yml' for YAML format.
	InitWithConfig(path string) error

	// Get iris application.
	Iris() *iris.Application

	// NewService make Service to register handler.
	//
	// Parameters:
	//   - name: for statsd
	//   - path: path from root, e.g. "foo" for "/foo"
	NewService(name, path string) *Service

	// NewServiceGroup make ServiceGroup to handle multiple Services.
	//
	// Parameters:
	//   - name: for statsd
	//   - path: directory name in path from root
	NewServiceGroup(name, path string) *ServiceGroup

	// Run server with <host> and multiple [conf].
	//
	// Parameters:
	//   - host: [server name] with port. e.g. "mydomain.com:80" or ":8080" (equal "0.0.0.0:8080")
	Run(host string, conf ...host.Configurator)
}

func NewApplication() Application {
	a := &application{iris: iris.New()}
	a.rootGroup = ServiceGroup{app: a}
	a.iris.Use(rockMiddleware)
	return a
}

var _ Application = (*application)(nil)

type application struct {
	iris      *iris.Application
	rootGroup ServiceGroup
}

func (a *application) InitWithConfig(path string) error {
	// read config
	var cfg util.AnyMap
	err := LoadConfigFromFile(path, &cfg)
	if err != nil {
		return util.NewError(ErrNameApplicationInitFailure, err)
	}
	// create loggers
	if cfgIt := util.AnyToAnyMap(cfg["log"]); cfgIt != nil {
		for name, v := range cfgIt {
			vs := util.AnyToAnyMap(v)
			if vs == nil {
				continue
			}
			logger, err := parseLogger(name, vs)
			if err != nil {
				return util.NewError(ErrNameApplicationInitFailure, err)
			}
			loggers[name] = logger
		}
	}
	if it := loggers["default"]; it != nil {
		defaultLogger = *it
	}
	if cfgIt := util.AnyToAnyMap(cfg["metric"]); cfgIt != nil {
		initMetric(cfgIt)
	}
	if cfgIt := util.AnyToAnyMap(cfg["sentry"]); cfgIt != nil {
		sentryMW, err := initSentryMiddleware(parseSentryOption(cfgIt))
		if err != nil {
			return util.NewError(ErrNameApplicationInitFailure, err)
		}
		a.iris.Use(sentryMW)
	}
	a.rootGroup.name = metricPrefix
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

func (a *application) Run(host string, conf ...host.Configurator) {
	a.iris.Run(iris.Addr(host, conf...))
}
