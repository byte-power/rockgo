package rock

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const apiPrefix = "api."
const timecostSuffix = ".timecost"
const statusCodeSuffixOK = ".ok"

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
	name := route.Path()
	code := ctx.GetStatusCode()
	dur := time.Now().Sub(startHandleTime)
	method := strings.ToLower(ctx.Method())
	var codeExpl string
	if code > 100 && code < 400 {
		codeExpl = statusCodeSuffixOK
	} else if code >= 400 && code < 500 {
		codeExpl = ".4xx"
	} else if code >= 500 {
		codeExpl = fmt.Sprintf(".%v", code)
	}

	var buf strings.Builder
	buf.Reset()
	// count: api.all
	MetricIncrease(writeStrings(&buf, apiPrefix, "all"))
	// timecost: api.all.timecost
	buf.WriteString(timecostSuffix)
	MetricTimeDuration(buf.String(), dur)

	buf.Reset()
	for _, suffix := range []string{method, "all"} {
		buf.Reset()
		// count: api.{path}.[{method}|all]
		MetricIncrease(writeStrings(&buf, apiPrefix, name, ".", suffix))
		if codeExpl != "" {
			// count: api.{path}.[{method}|all].[ok|4xx|5xx]
			buf.WriteString(codeExpl)
			MetricIncrease(buf.String())
		}
	}
	// timecost: api.{path}.all.timecost
	MetricTimeDuration(writeStrings(&buf, apiPrefix, name, ".all", timecostSuffix), dur)
	if codeExpl != "" {
		buf.Reset()
		// count: api.all.[ok|4xx|5xx]
		MetricIncrease(writeStrings(&buf, apiPrefix, "all", codeExpl))
		// count: api.all.[ok|4xx|5xx].timecost
		buf.WriteString(timecostSuffix)
		MetricTimeDuration(buf.String(), dur)
		if codeExpl == statusCodeSuffixOK {
			buf.Reset()
			// timecost: api.{path}.all.ok.timecost
			MetricTimeDuration(writeStrings(&buf, apiPrefix, name, ".all", codeExpl, timecostSuffix), dur)
		}
	}
}

func writeStrings(buf *strings.Builder, strs ...string) string {
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf.String()
}
