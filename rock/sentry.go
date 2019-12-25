package rock

import (
	"time"

	"github.com/byte-power/rockgo/util"
	"github.com/getsentry/sentry-go"
	sentryiris "github.com/getsentry/sentry-go/iris"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type PanicHandler func(ctx iris.Context, err error)

var panicHandler PanicHandler

// 设置panic信息获取器，仅当sentry.repanic=true时生效
func SetPanicHandler(fn PanicHandler) {
	panicHandler = fn
}

func parseSentryOption(m util.AnyMap) (sentry.ClientOptions, sentryiris.Options) {
	opt := sentry.ClientOptions{}
	mwOPT := sentryiris.Options{
		Repanic: util.AnyToBool(m["repanic"]),
	}
	if it := util.AnyToString(m["dsn"]); it != "" {
		opt.Dsn = it
	}
	if it, ok := m["debug"]; ok {
		opt.Debug = util.AnyToBool(it)
	}
	if it, ok := m["sample_rate"]; ok {
		opt.SampleRate = util.AnyToFloat64(it)
	}
	if it, ok := m["ignore_errors"].([]interface{}); ok {
		opt.IgnoreErrors = make([]string, len(it))
		for i, name := range it {
			opt.IgnoreErrors[i] = util.AnyToString(name)
		}
	}
	if it, ok := m["server_name"]; ok {
		opt.ServerName = util.AnyToString(it)
	}
	if it, ok := m["release"]; ok {
		opt.Release = util.AnyToString(it)
	}
	if it, ok := m["dist"]; ok {
		opt.Dist = util.AnyToString(it)
	}
	if it, ok := m["environment"]; ok {
		opt.Environment = util.AnyToString(it)
	}
	if it, ok := m["max_breadcrumbs"]; ok {
		opt.MaxBreadcrumbs = int(util.AnyToInt64(it))
	}
	if it, ok := m["timeout_seconds"]; ok {
		mwOPT.Timeout = time.Second * time.Duration(util.AnyToInt64(it))
	}
	return opt, mwOPT
}

func newSentryMiddleware(opt sentry.ClientOptions, mwOPT sentryiris.Options) (context.Handler, error) {
	err := sentry.Init(opt)
	if err != nil {
		return nil, err
	}
	return sentryiris.New(mwOPT), nil
}
