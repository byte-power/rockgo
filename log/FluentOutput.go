package log

import (
	"fmt"
	"strings"

	"github.com/byte-power/rockgo/util"
	"github.com/fluent/fluent-logger-golang/fluent"
)

type FluentOutput struct {
	level  Level
	host   string
	port   int
	tag    string
	async  bool
	output *fluent.Fluent
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
	data := make(util.AnyMap)
	data["message"] = msg
	data["data"] = util.AnyArrayToMap(argPairs)
	o.output.Post(o.tag, data)
}

func (o *FluentOutput) LogMap(l Level, msg string, values map[string]interface{}) {
	data := make(util.AnyMap)
	data["message"] = msg
	data["data"] = values
	o.output.Post(o.tag, data)
}

func (o *FluentOutput) LogPlainMessage(l Level, args []interface{}) {
	data := make(util.AnyMap)
	data["message"] = strings.Join(util.AnyArrayToStringArray(args), "")
	o.output.Post(o.tag, data)
}

func (o *FluentOutput) LogFormatted(l Level, format string, args []interface{}) {
	data := make(util.AnyMap)
	data["message"] = fmt.Sprintf(format, args...)
	o.output.Post(o.tag, data)
}
