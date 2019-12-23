package rock

import (
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
)

type Application interface {
	// 载入config_rock，以初始化各个内部模块，并自动添加用于异常恢复的recover和sentry中间件、记录基础metric信息的中间件
	// InitWithConfig would load config from [path] and then initialize each internal components.
	InitWithConfig(path string) error

	// 取得iris实例，以注册控制器、中间件等
	Iris() *iris.Application

	// 生成服务对象，用以注册各种请求的handler
	NewService(name, path string) *Service
	NewServiceGroup(name, path string) *ServiceGroup

	// 启动服务，host可包含端口号，如省略域名或ip将与0.0.0.0等效
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
		return err
	}
	// create loggers
	if cfgIt := util.AnyToAnyMap(cfg["log"]); cfgIt != nil {
		for name, v := range cfgIt {
			vs := util.AnyToAnyMap(v)
			if vs == nil {
				// TODO: warn parse failed
				continue
			}
			logger, err := parseLogger(name, vs)
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
		initMetric(cfgIt)
	}
	if cfgIt := util.AnyToAnyMap(cfg["sentry"]); cfgIt != nil {
		sentryMW, err := initSentryMiddleware(parseSentryOption(cfgIt))
		if err != nil {
			return err
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
