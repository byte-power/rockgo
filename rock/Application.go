package rock

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	// 启动服务，host可包含端口号，如省略域名或ip将与0.0.0.0等效
	Run(host string, conf ...host.Configurator)
}

func NewApplication() Application {
	a := &application{
		iris:          iris.New(),
		services:      map[string]*Service{},
		handlerStatus: map[string]bool{},
	}
	a.iris.Use(func(ctx iris.Context) {
		startHandleTime := time.Now()
		defer func() {
			recovered := recover()
			if recovered != nil {
				var err error
				switch v := recovered.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = errors.New(util.AnyToString(v))
				}
				if fn := panicHandler; fn != nil {
					fn(ctx, err)
				} else {
					ctx.StatusCode(http.StatusInternalServerError)
					ctx.JSON(util.AnyMap{"error": err.Error()})
				}
			}
			route := ctx.GetCurrentRoute()
			if route != nil {
				name := route.MainHandlerName()
				appPrefix := ""
				code := ctx.GetStatusCode()
				var codeExpl string
				if code > 100 && code < 400 {
					codeExpl = ".ok"
				} else if code >= 400 && code < 500 {
					codeExpl = ".4xx"
				} else if code >= 500 {
					codeExpl = ".5xx"
				}
				// record: [{appName}.api.{path}.{method} | {appName}.api.{path}.all | {appName}.api.all] * [status 100~399 | 4xx | 5xx]
				statsPrefix := appPrefix + ".api."
				MetricIncrease(statsPrefix + name + "." + strings.ToLower(ctx.Method()))
				MetricIncrease(statsPrefix + name + ".all")
				MetricIncrease(statsPrefix + "all")
				if codeExpl != "" {
					MetricIncrease(statsPrefix + name + "." + strings.ToLower(ctx.Method()) + codeExpl)
					MetricIncrease(statsPrefix + name + ".all" + codeExpl)
					MetricIncrease(statsPrefix + "all" + codeExpl)
				}
				// record: time cost - [{appName}.api.all | {appName}.api.{path}.all] * [status 100~399 | all]
				// 以便后期统计 min, mean, max, all, 90%
				dur := startHandleTime.Sub(time.Now())
				MetricTiming(statsPrefix+"all", dur)
				MetricTiming(statsPrefix+name+".all", dur)
				if codeExpl != "" {
					MetricTiming(statsPrefix+"all"+codeExpl, dur)
					MetricTiming(statsPrefix+name+".all"+codeExpl, dur)
				}
			}
		}()
		ctx.Next()
	})
	return a
}

var _ Application = (*application)(nil)

type application struct {
	iris *iris.Application

	services      map[string]*Service
	handlerStatus map[string]bool
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
	return nil
}

func (a *application) NewService(name, path string) *Service {
	if name == "" {
		defaultLogger.Warn("Service should named with non-zero length")
	}
	if path == "" || path[0] != '/' {
		panic(fmt.Errorf("Service(%s) path(%s) should start with '/'", name, path))
	}
	if _, ok := a.services[name]; ok {
		// TODO: warn name exist
	}
	s := &Service{app: a, path: path, name: name}
	a.services[name] = s
	return s
}

func (a *application) Iris() *iris.Application {
	return a.iris
}

func (a *application) Run(host string, conf ...host.Configurator) {
	a.iris.Run(iris.Addr(host, conf...))
}
