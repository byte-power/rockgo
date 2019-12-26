# RockGo Application Framework

RockGo is fast, simple application framework for Go.

## Features
* Application
	* [x] Service, ServiceGroup
	* [x] Config
	* [x] Basic middleware
		* Access log
		* recover & metric
	* [80%] Logger integration
	* [x] Metric (Stats)
	* [x] Sentry
* Log
	* [x] Logger, Output, Format
	* [x] zap
	* [10%] fluent
* Crypto
	* [ ] AES
	* [ ] Digest (MD5, SHA1/256/512)
	* [ ] RSA
* Example
	* [x] Route - Application, Config
	* [ ] Metric & Sentry
* Test case
	* [ ] Log
	* [ ] Crypto
	* [x] utility

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
