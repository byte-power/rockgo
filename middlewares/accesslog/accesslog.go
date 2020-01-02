package accesslog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/byte-power/rockgo/log"
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const timeFormat = "02/Jan/2006:15:04:05 -0700"

// New make iris middleware to log access log with <logger> by log.LevelInfo.
func New(logger *log.Logger) context.Handler {
	return func(ctx iris.Context) {
		ctx.Next()
		if logger != nil {
			length := ctx.ResponseWriter().Header().Get("Content-Length")
			l := Sprint(ctx.Request(), ctx.RemoteAddr(), time.Now(), ctx.GetStatusCode(), util.AnyToInt64(length))
			logger.Info(l)
		}
	}
}

// Sprint with iris.Context and time.
// ref: https://en.wikipedia.org/wiki/Common_Log_Format
// - Parameters:
//   - req: HTTP reqeust
//   - remoteAddr: Requester address (or client address), e.g. iris.Context#RemoteAddr()
//   - t: Time on responsing, e.g. time.Now()
//   - statusCode: Response status code
//   - contentLength: Response content length, e.g. convert ResponseWriter#.Header().Get("Content-Length") to int64
func Sprint(req *http.Request, remoteAddr string, t time.Time, statusCode int, contentLength int64) string {
	return fmt.Sprintf("%s - - [%s] \"%s %s %s\" %v %v",
		remoteAddr,
		t.Format(timeFormat),
		req.Method,
		req.RequestURI,
		req.Proto,
		statusCode,
		contentLength)
}
