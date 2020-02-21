/*
Package realip 用于获取用户真实的 ip。对于 client api 和 server api 我们采用不同
的方式来获取 ip。client api 是指用户直接调用的 api，比如供安卓 sdk、ios sdk 或
web 端直接使用的 api。server api 是指供其它 server 端服务调用的 api，比如 sdk
调用了 api1，api1 调用了 api2，那么 api2 就是 server api，api1 就是 client api。

client api 获取 ip 的方式：
    1. 确定 proxy 的数量，比如经过了 alb 和 nginx 转发，那么 proxy 的数量就是 2。
    2. 如果经过了 amazon cdn，proxy 的数量需要加 1，因为多了一次转发。目前
       amazon cdn 的行为是这样的，其它 cdn 需要视实际情况而定。
    3. 如果 proxy 的数量小于等于 X-Forwarded-For 的 ip 数量，取 X-Forwarded-For
       倒数第 proxy number 个 ip 作为用户真实的 ip。
    4. 如果经过以上过程没有提取到 ip，就用 remote addr 作为用户的 ip，remote
       addr 是指用户连接我们服务用的 ip 地址。

server api 获取 ip 的方式：
    1. 如果 http header 设置了 X-User-Real-Ip 字段，就把这个字段作为真实的 ip。
       这个字段是我们自定义的字段，需要调用方设置。
    2. 如果调用方未设置 X-User-Real-Ip，取 X-Forwarded-For 的第一个 ip。
    3. 如果以上过程都没有提取到 ip，就用 remote addr 作为用户的 ip。

examples：
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
        realIP := ""
        if value := ctx.Values().Get(RealIPKey); value != nil {
            realIP = value.(string)
        }
        ctx.Writef(realIP)
        ctx.Next()
    }
*/
package realip

import (
	"net"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// RealIpKey 是我们自定义的一个 HTTP Header，当用户调用 server api 时，需要设置
// 这个值。
const RealIPKey = "X-User-Real-Ip"

// NewForClientApi 接收 proxyNumber 和 proxyHandlers 参数，返回获取用户真实 ip
// 的 iris middleware，这个 middleware 适用于 client api。
// proxyNumber 表示 proxy 的数量，这个值应该根据项目的实际部署环境而定。
// proxyHandlers 允许用户自定义一些额外的确定 proxy 数量的规则，例如：
// func myProxyHandler(ctx iris.Context) uint {
//     extraProxyNumber := 0
//     // 一些确定额外的 proxy 数量的规则
//     return extraProxyNumber
// }
// middleware := NewForClientApi(1, myProxyHandler)
// 最终的 proxy 数量由 proxyNumber 加上所有 proxyHandlers 的返回值得到。
func NewForClientApi(proxyNumber uint8, proxyHandlers ...ProxyHandler) context.Handler {
	cfg := configForClientApi{
		proxyNumber:   proxyNumber,
		proxyHandlers: proxyHandlers,
	}
	cfg.addHandler(checkAmazonCDNProxyHandler)
	return cfg.serve
}

// NewForServerApi 返回获取用户真实 ip 的 iris middleware，这个 middleware
// 适用于 server api。
func NewForServerApi() context.Handler {
	return func(ctx iris.Context) {
		realIP := ctx.GetHeader(RealIPKey)
		if realIP == "" {
			realIP = getIpByXForwardedFor(ctx.GetHeader("X-Forwarded-For"), 0)
		}
		if realIP == "" {
			realIP = getIpByRemoteAddr(ctx.Request().RemoteAddr)
		}

		ctx.Values().Set(RealIPKey, realIP)
		ctx.Next()
	}
}

// ProxyHandler 接收 iris.Context 参数，返回 proxy 的数量，它主要用于让使用者
// 自定义一些确定 proxy 数量的规则。
type ProxyHandler func(iris.Context) uint

type configForClientApi struct {
	proxyNumber   uint8
	proxyHandlers []ProxyHandler
}

func (cfg *configForClientApi) serve(ctx iris.Context) {
	proxyNumber := uint(cfg.proxyNumber)

	for _, hdlr := range cfg.proxyHandlers {
		proxyNumber += hdlr(ctx)
	}

	realIP := ""
	if proxyNumber > 0 {
		realIP = getIpByXForwardedFor(ctx.GetHeader("X-Forwarded-For"), -int(proxyNumber))
	}
	if realIP == "" {
		realIP = getIpByRemoteAddr(ctx.Request().RemoteAddr)
	}

	ctx.Values().Set(RealIPKey, realIP)
	ctx.Next()
}

func (cfg *configForClientApi) addHandler(hdlr ProxyHandler) {
	cfg.proxyHandlers = append(cfg.proxyHandlers, hdlr)
}

func checkAmazonCDNProxyHandler(ctx iris.Context) uint {
	userAgent := ctx.GetHeader("User-Agent")
	if strings.Contains(userAgent, "Amazon CloudFront") {
		return 1
	}
	return 0
}

func getIpByXForwardedFor(xForwardedFor string, index int) string {
	var xForwardedForArray []string
	for _, v := range strings.Split(xForwardedFor, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			xForwardedForArray = append(xForwardedForArray, v)
		}
	}
	if index < 0 {
		index += len(xForwardedForArray)
	}
	if index >= 0 && index < len(xForwardedForArray) {
		return xForwardedForArray[index]
	}
	return ""
}

func getIpByRemoteAddr(remoteAddr string) string {
	realIP := strings.TrimSpace(remoteAddr)
	if realIP != "" {
		// if addr has port use the net.SplitHostPort otherwise(error occurs) take as it is
		if ip, _, err := net.SplitHostPort(realIP); err == nil {
			realIP = ip
		}
	}
	return realIP
}
