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
	config fluent.Config
	tag    string

	output interfaceFluent
}

func (o *FluentOutput) Level() Level {
	return o.level
}

func (o *FluentOutput) Init() {
	logger, error := fluent.New(o.config)
	if error != nil {
		panic(error)
	}
	o.output = logger
}

func (o *FluentOutput) Log(l Level, msg string, argPairs []interface{}) {
	args := util.AnyArrayToStrMap(argPairs)
	o.LogMap(l, msg, args)
}

func (o *FluentOutput) LogMap(l Level, msg string, values map[string]interface{}) {
	data := make(util.StrMap)
	data["message"] = msg
	if values != nil {
		data["data"] = values
	}
	err := o.output.Post(o.tag, data)
	if err != nil {
		fmt.Println("FluentOutputPostError:", err)
	}
}

func (o *FluentOutput) LogPlainMessage(l Level, args []interface{}) {
	msg := strings.Join(util.AnyArrayToStringArray(args), "")
	o.LogMap(l, msg, nil)
}

func (o *FluentOutput) LogFormatted(l Level, format string, args []interface{}) {
	msg := fmt.Sprintf(format, args...)
	o.LogMap(l, msg, nil)
}
