package log

import (
	"fmt"
	"strings"

	"github.com/byte-power/rockgo/util"
	"github.com/fluent/fluent-logger-golang/fluent"
)

type interfaceFluent interface {
	Post(tag string, message interface{}) error
}

type FluentOutput struct {
	level  Level
	host   string
	port   int
	tag    string
	async  bool
	output interfaceFluent
}

func (o *FluentOutput) Level() Level {
	return o.level
}

func (o *FluentOutput) Init() {
	logger, error := fluent.New(fluent.Config{FluentPort: o.port, FluentHost: o.host, Async: o.async})
	if error != nil {
		panic(error)
	}
	o.output = logger
}

func (o *FluentOutput) Log(l Level, msg string, argPairs []interface{}) {
	args := util.AnyArrayToMap(argPairs)
	o.LogMap(l, msg, args)
}

func (o *FluentOutput) LogMap(l Level, msg string, values map[string]interface{}) {
	data := make(util.AnyMap)
	data["message"] = msg
	if values != nil {
		data["data"] = values
	}
	o.output.Post(o.tag, data)
}

func (o *FluentOutput) LogPlainMessage(l Level, args []interface{}) {
	msg := strings.Join(util.AnyArrayToStringArray(args), "")
	o.LogMap(l, msg, nil)
}

func (o *FluentOutput) LogFormatted(l Level, format string, args []interface{}) {
	msg := fmt.Sprintf(format, args...)
	o.LogMap(l, msg, nil)
}
