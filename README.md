# RockGo Application Framework

RockGo is fast, simple application framework for Go.

RockGo is agent to iris, fluentd, statsd, zap and sentry too. It make easy to build perfect application or service.

## Features
* Application
	* [x] Service, ServiceGroup
	* [x] Config
	* [x] Basic middleware
		* Access log
		* recover & metric
	* [x] Logger integration
	* [x] Metric (Stats)
	* [x] Sentry
* Log
	* [x] Logger, Output, Format
	* [x] zap
	* [x] fluent
* Crypto
	* [x] AES
	* [x] Digest (MD5, SHA1/256/512)
	* [x] RSA
* Example
	* [x] Route - Application, Config
	* [ ] Metric & Sentry

## Example
Please visit [example](/tree/master/_example).

## Install
* go get
```
$ go get github.com/byte-power/rockgo
```
* Add import
```go
import "github.com/byte-power/rockgo/rock"
```
* Add config file named "rockgo.yaml" for internal modules on your settings directory
```yaml
app_name: myapp
log:
  LoggerName:
    console:
    fluent:
      level: info
      host: "myfluenthost.com"
      port: 24224
      async: true
metric:
  host: "127.0.0.1:8125"
sentry:
  dsn: "http://user@127.0.0.1/1"
  repanic: true
```

* Append routes & Run server
```go
func main() {
	// load each config file include rockgo.yaml in the directory to create Application
	app, err := rock.NewApplication("settings")
	if err != nil {
		panic(err)
	}
	// register route handler with Service
	app.NewService("root", "/").Get(func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	})
	app.Run(":8080")
}
```

## License
MIT
