package rock

import (
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
	// TODO: ...
	return opt, mwOPT
}

func initSentryMiddleware(opt sentry.ClientOptions, mwOPT sentryiris.Options) (context.Handler, error) {
	err := sentry.Init(opt)
	if err != nil {
		return nil, err
	}
	return sentryiris.New(mwOPT), nil
}
