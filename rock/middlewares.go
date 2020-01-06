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
	dur := time.Now().Sub(startHandleTime)
	method := strings.ToLower(ctx.Method())
	var codeExpl []byte
	if code > 100 && code < 400 {
		codeExpl = []byte(".ok")
	} else if code >= 400 && code < 500 {
		codeExpl = []byte(".4xx")
	} else if code >= 500 {
		codeExpl = []byte(".5xx")
	}
	// record for analytic: min, mean, max, all, 90%
	// count: [api.{path}.{method} | api.{path}.all | api.all] * (status 100~399 | 4xx | 5xx)
	// time cost: [api.all | api.{path}.all] * [(status 100~399) | all]

	var buf strings.Builder
	var key string
	for _, suffix := range []string{method, "all"} {
		buf.Reset()
		buf.WriteString(apiPrefix)
		buf.WriteString(name)
		buf.WriteByte('.')
		buf.WriteString(suffix)
		MetricIncrease(buf.String())
		if codeExpl != nil {
			buf.Write(codeExpl)
			MetricIncrease(buf.String())
		}
	}
	buf.Reset()
	buf.WriteString(apiPrefix)
	buf.WriteString("all")
	key = buf.String()
	MetricIncrease(key)
	MetricTiming(key, dur)
	if codeExpl != nil {
		buf.Write(codeExpl)
		key = buf.String()
		MetricTiming(key, dur)
		MetricIncrease(key)
	}
	buf.Reset()
	buf.WriteString(apiPrefix)
	buf.WriteString(name)
	buf.WriteString(".all")
	MetricTiming(buf.String(), dur)
	if codeExpl != nil {
		buf.Write(codeExpl)
		MetricTiming(buf.String(), dur)
	}
}
