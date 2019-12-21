package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/byte-power/rockgo/rock"
	"github.com/kataras/iris/v12"
)

func main() {
	app := rock.NewApplication()
	if err := app.InitWithConfig("config.yaml"); err != nil {
		panic(err)
	}

	logger := rock.Logger("Main")
	logger.Debug("should not display")

	app.Iris().OnErrorCode(http.StatusNotFound, func(ctx iris.Context) {
		println("application.on404", ctx.Request().Method, ctx.Request().URL.String())
	})

	app.Iris().Use(rock.NewAccessLogMiddleware(rock.Logger("Access")))

	rock.SetPanicHandler(func(ctx iris.Context, err error) {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Text(string(debug.Stack()))
		rock.Logger("Main").Error(err.Error())
	})

	app.NewService("", "/").
		Get(func(ctx iris.Context) {
			ctx.StatusCode(http.StatusOK)
		})

	app.NewService("user", "/user/{id:int min(1)}").
		Get(func(ctx iris.Context) {
			id, err := ctx.Params().GetInt("id")
			ctx.ResponseWriter().Write([]byte(fmt.Sprintf("get user:%v(%v)", id, err)))
		}).
		Post(func(ctx iris.Context) {
			id, _ := ctx.Params().GetInt("id")
			ctx.ResponseWriter().Write([]byte(fmt.Sprintf("post user:%v", id)))
		})
	handleWorkspaces(app)

	app.NewService("fatal", "/fatal").
		Get(func(ctx iris.Context) {
			panic("PanicErrInfo")
		})

	logger.Info("Server running...")
	app.Run(":8080")
}
