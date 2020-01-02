package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/byte-power/rockgo/rock"
	"github.com/byte-power/rockgo/middlewares/accesslog"
	"github.com/kataras/iris/v12"
)

func main() {
	// make new application with config directory contains rockgo.yaml
	app, err := rock.NewApplication("settings")
	if err != nil {
		panic(err)
	}

	// get the "Main" logger was defined in rockgo.yaml, use it output something
	logger := rock.Logger("Main")
	logger.Debug("should not display")
	// you can get config by Config(), or ConfigIn(keyPath)
	logger.Infof("loaded config: %v", rock.Config())
	logger.Info("loaded config:", "a", rock.ConfigIn("jd.a"), "z", rock.ConfigIn("yd.z"))

	// register 404 handler
	app.Iris().OnErrorCode(http.StatusNotFound, func(ctx iris.Context) {
		println("application.on404", ctx.Request().Method, ctx.Request().URL.String())
	})

	// register middleware - use logger named "Access" to record access log
	app.Iris().Use(accesslog.New(rock.Logger("Access")))

	// use panic handler control behavior on panic
	app.SetPanicHandler(func(ctx iris.Context, err error) {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Text(string(debug.Stack()))
		rock.Logger("Main").Error(err.Error())
	})

	// handle get /
	app.Serve("root", "/").
		Get(func(ctx iris.Context) {
			ctx.StatusCode(http.StatusOK)
			ctx.Text("Hello World")
		})

	// handle get and post /user/{id}
	app.Serve("user", "/user/{id:int min(1)}").
		Get(func(ctx iris.Context) {
			id, err := ctx.Params().GetInt("id")
			ctx.ResponseWriter().Write([]byte(fmt.Sprintf("get user:%v(%v)", id, err)))
		}).
		Post(func(ctx iris.Context) {
			id, _ := ctx.Params().GetInt("id")
			ctx.ResponseWriter().Write([]byte(fmt.Sprintf("post user:%v", id)))
		})
	// handle some API defined in other file (article.go)
	handleArticles(app)

	// make a group for directory based middleware
	g1 := app.ServeGroup("g1", "/g")
	g1.Use(func(ctx iris.Context) {
		println("g1 only middleware")
		ctx.Next()
	})
	g1.Serve("a", "/a").Get(func(ctx iris.Context) {
		ctx.Text("get /g/a")
	})
	// group in group
	g1.ServeGroup("2", "/2").Serve("3", "/3").Get(func(ctx iris.Context) {
		ctx.Text("get /g/2/3")
	})

	// request /fatal would make panic, you can test panic handler with this
	app.Serve("fatal", "/fatal").
		Get(func(ctx iris.Context) {
			panic("PanicErrInfo")
		})

	// Handle directory access
	app.Serve("settings", "/conf").HandleDir("settings")

	logger.Infof("Server %v running...", app.Name())
	// make application run up, you can change host or port, or append some iris configuration
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
