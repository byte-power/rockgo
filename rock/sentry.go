package rock

import (
	"time"

	"github.com/byte-power/rockgo/util"
	"github.com/getsentry/sentry-go"
	sentryiris "github.com/getsentry/sentry-go/iris"
	"github.com/kataras/iris/v12/context"
)

func parseSentryOption(cfg util.AnyMap) (sentry.ClientOptions, sentryiris.Options) {
	opt := sentry.ClientOptions{}
	mwOPT := sentryiris.Options{
		Repanic: util.AnyToBool(cfg["repanic"]),
	}
	if it := util.AnyToString(cfg["dsn"]); it != "" {
		opt.Dsn = it
	}
	if it, ok := cfg["debug"]; ok {
		opt.Debug = util.AnyToBool(it)
	}
	if it, ok := cfg["sample_rate"]; ok {
		opt.SampleRate = util.AnyToFloat64(it)
	}
	if it, ok := cfg["ignore_errors"].([]interface{}); ok {
		opt.IgnoreErrors = make([]string, len(it))
		for i, name := range it {
			opt.IgnoreErrors[i] = util.AnyToString(name)
		}
	}
	if it, ok := cfg["server_name"]; ok {
		opt.ServerName = util.AnyToString(it)
	}
	if it, ok := cfg["release"]; ok {
		opt.Release = util.AnyToString(it)
	}
	if it, ok := cfg["dist"]; ok {
		opt.Dist = util.AnyToString(it)
	}
	if it, ok := cfg["environment"]; ok {
		opt.Environment = util.AnyToString(it)
	}
	if it, ok := cfg["max_breadcrumbs"]; ok {
		opt.MaxBreadcrumbs = int(util.AnyToInt64(it))
	}
	if it, ok := cfg["timeout_seconds"]; ok {
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
