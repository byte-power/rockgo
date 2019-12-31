package rock

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const apiPrefix = "api."

type panicHandlerProvider interface {
	PanicHandler() PanicHandler
}

func newRockMiddleware(provider panicHandlerProvider) context.Handler {
	return func(ctx iris.Context) {
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
				var fn PanicHandler
				if provider != nil {
					fn = provider.PanicHandler()
				}
				if fn != nil {
					fn(ctx, err)
				} else {
					ctx.StatusCode(http.StatusInternalServerError)
					ctx.JSON(util.AnyMap{"error": err.Error()})
				}
			}
			recordMetric(ctx, startHandleTime)
		}()
		ctx.Next()
	}
}

func recordMetric(ctx iris.Context, startHandleTime time.Time) {
	route := ctx.GetCurrentRoute()
	if route == nil || Metric() == nil {
		return
	}
	name := route.MainHandlerName()
	code := ctx.GetStatusCode()
	dur := startHandleTime.Sub(time.Now())
	method := strings.ToLower(ctx.Method())
	var codeExpl string
	if code > 100 && code < 400 {
		codeExpl = ".ok"
	} else if code >= 400 && code < 500 {
		codeExpl = ".4xx"
	} else if code >= 500 {
		codeExpl = ".5xx"
	}
	prefixName := apiPrefix + name
	prefixAll := apiPrefix + "all"
	// record: [{appName}.api.{path}.{method} | {appName}.api.{path}.all | {appName}.api.all] * [status 100~399 | 4xx | 5xx]
	MetricIncrease(prefixName + "." + method)
	MetricIncrease(prefixName + ".all")
	MetricIncrease(prefixAll)
	if codeExpl != "" {
		MetricIncrease(prefixName + "." + method + codeExpl)
		MetricIncrease(prefixName + ".all" + codeExpl)
		MetricIncrease(prefixAll + codeExpl)
	}
	// record: time cost - [{appName}.api.all | {appName}.api.{path}.all] * [status 100~399 | all]
	// 以便后期统计 min, mean, max, all, 90%
	MetricTiming(prefixAll, dur)
	MetricTiming(prefixName+".all", dur)
	if codeExpl != "" {
		MetricTiming(prefixAll+codeExpl, dur)
		MetricTiming(prefixName+".all"+codeExpl, dur)
	}
}
