package rock

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byte-power/rockgo/log"
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const (
	timeFormat = "02/Jan/2006:15:04:05 -0700"
	apiPrefix  = "api."
)

// NewAccessLogMiddleware make iris middleware to log access log with Info.
func NewAccessLogMiddleware(logger *log.Logger) context.Handler {
	return func(ctx iris.Context) {
		ctx.Next()
		if logger != nil {
			l := MakeAccessLog(ctx, time.Now())
			logger.Info(l)
		}
	}
}

// ref: https://en.wikipedia.org/wiki/Common_Log_Format
func MakeAccessLog(ctx iris.Context, t time.Time) string {
	req := ctx.Request()
	return fmt.Sprintf("%s - - [%s] \"%s %s %s\" %v %v",
		ctx.RemoteAddr(),
		t.Format(timeFormat),
		req.Method,
		req.RequestURI,
		req.Proto,
		ctx.GetStatusCode(),
		ctx.ResponseWriter().Header().Get("Content-Length"))
}

func newRockMiddleware(app *application) context.Handler {
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
				if fn := app.panicHandler; fn != nil {
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
