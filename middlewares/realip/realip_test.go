package realip

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func newApp() *iris.Application {
	app := iris.New()

	clientApi := app.Party("/")
	clientApi.Use(NewForClientApi(1))
	clientApi.Get("/", set_ip_handler)

	serverApi := app.Party("/serverapi")
	serverApi.Use(NewForServerApi())
	serverApi.Get("/", set_ip_handler)

	return app
}

func set_ip_handler(ctx iris.Context) {
	realIp := ""
	if value := ctx.Values().Get(RealIpKey); value != nil {
		realIp = value.(string)
	}
	ctx.Writef(realIp)
	ctx.Next()
}

func TestClientApi(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/").Expect().Body().Equal("")
	e.GET("/").WithHeader("X-Forwarded-For", "1.1.1.1").Expect().Body().Equal("1.1.1.1")
	e.GET("/").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2").Expect().Body().Equal("2.2.2.2")
	e.GET("/").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2,3.3.3.3").Expect().Body().Equal("3.3.3.3")
	e.GET("/").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2,3.3.3.3").WithHeader("User-Agent", "Amazon CloudFront").Expect().Body().Equal("2.2.2.2")
}

func TestServerApi(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/serverapi").Expect().Body().Equal("")
	e.GET("/serverapi").WithHeader("X-Forwarded-For", "1.1.1.1").Expect().Body().Equal("1.1.1.1")
	e.GET("/serverapi").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2").Expect().Body().Equal("1.1.1.1")
	e.GET("/serverapi").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2,3.3.3.3").Expect().Body().Equal("1.1.1.1")
	e.GET("/serverapi").WithHeader("X-Forwarded-For", "1.1.1.1,2.2.2.2,3.3.3.3").WithHeader(RealIpKey, "4.4.4.4").Expect().Body().Equal("4.4.4.4")
}
