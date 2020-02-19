package realip

import (
	"net"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const RealIpKey = "X-User-Real-Ip"

func NewForClientApi(proxy_number int) context.Handler {
	return func(ctx iris.Context) {
		proxy_number := proxy_number

		userAgent := ctx.GetHeader("User-Agent")
		if strings.Contains(userAgent, "Amazon CloudFront") {
			proxy_number += 1
		}

		realIp := ""
		if proxy_number > 0 {
			realIp = getIpByXForwardedFor(ctx.GetHeader("X-Forwarded-For"), -proxy_number)
		}
		if realIp == "" {
			realIp = getIpByRemoteAddr(ctx.Request().RemoteAddr)
		}

		ctx.Values().Set(RealIpKey, realIp)
		ctx.Next()
	}
}

func NewForServerApi() context.Handler {
	return func(ctx iris.Context) {
		realIp := ctx.GetHeader(RealIpKey)
		if realIp == "" {
			realIp = getIpByXForwardedFor(ctx.GetHeader("X-Forwarded-For"), 0)
		}
		if realIp == "" {
			realIp = getIpByRemoteAddr(ctx.Request().RemoteAddr)
		}

		ctx.Values().Set(RealIpKey, realIp)
		ctx.Next()
	}
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
	realIp := strings.TrimSpace(remoteAddr)
	if realIp != "" {
		// if addr has port use the net.SplitHostPort otherwise(error occurs) take as it is
		if ip, _, err := net.SplitHostPort(realIp); err == nil {
			realIp = ip
		}
	}
	return realIp
}
