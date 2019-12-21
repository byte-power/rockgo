package log

import "github.com/fluent/fluent-logger-golang/fluent"

type FluentOutput struct {
	level  Level
	output *fluent.Fluent
}

func (o *FluentOutput) Level() Level {
	return o.level
}

func (o *FluentOutput) Log(l Level, msg string, argPairs []interface{}) {

}

func (o *FluentOutput) LogMap(l Level, msg string, values map[string]interface{}) {

}

func (o *FluentOutput) LogPlainMessage(l Level, args []interface{}) {

}

func (o *FluentOutput) LogFormatted(l Level, format string, args []interface{}) {

}
