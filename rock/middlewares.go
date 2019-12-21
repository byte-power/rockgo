package rock

import (
	"fmt"
	"time"

	"github.com/byte-power/rockgo/log"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const timeFormat = "02/Jan/2006:15:04:05 -0700"

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
		ctx.RemoteAddr(), t.Format(timeFormat),
		req.Method, req.URL.Path, req.Proto,
		ctx.GetStatusCode(), ctx.ResponseWriter().Header().Get("Content-Length"))
}
