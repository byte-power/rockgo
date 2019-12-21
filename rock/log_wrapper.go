package rock

import (
	"fmt"

	"github.com/byte-power/rockgo/log"
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
)

type PanicHandler func(ctx iris.Context, err error)

var loggers = map[string]*log.Logger{}

// 用默认的logger防止调用Logger(name)时对应logger不存在而得到空指针
var defaultLogger log.Logger
var panicHandler PanicHandler

// 设置panic信息获取器，仅当sentry.repanic=true时生效
func SetPanicHandler(fn PanicHandler) {
	panicHandler = fn
}

// 取得在Application.Init初始化过的该名称对应的Logger
func Logger(name string) *log.Logger {
	if logger, ok := loggers[name]; ok {
		return logger
	}
	return &defaultLogger
}

func parseLogger(name string, m util.AnyMap) (*log.Logger, error) {
	var outputs []log.Output
	for k, v := range m {
		vs := util.AnyToAnyMap(v)
		if vs == nil && v != nil {
			return nil, fmt.Errorf("parse Logger失败：%v不是map", k)
		}
		switch k {
		case "console":
			it := parseConsoleLogger(name, vs)
			outputs = append(outputs, it)
		case "file":
			it := parseFileLogger(name, vs)
			outputs = append(outputs, it)
		case "fluent":
			it := parseFluentLogger(vs)
			outputs = append(outputs, it)
		}
	}
	return log.NewLogger(outputs...), nil
}

func parseConsoleLogger(name string, m util.AnyMap) log.Output {
	fmt := parseFormat(m)
	level := parseLevel(m["level"])
	stream := log.MakeConsoleStream(util.AnyToString(m["stream"]))
	return log.MakeConsoleOutput(name, fmt, level, stream)
}

func parseFileLogger(name string, m util.AnyMap) log.Output {
	fmt := parseFormat(m)
	level := parseLevel(m["level"])
	location := util.AnyToString(m["location"])
	rotation := parseFileRotation(util.AnyToAnyMap(m["rotation"]))
	return log.MakeFileOutput(name, fmt, level, location, rotation)
}

func parseFluentLogger(m util.AnyMap) log.Output {
	// TODO:
	level := parseLevel(m["level"])
	return log.MakeFluentOutput(level)
}

func parseFormat(m util.AnyMap) log.LocalFormat {
	msgFMT := parseMessageFormat(m["format"])
	fmt := log.MakeLocalFormat(msgFMT)
	if keys := util.AnyToAnyMap(m["keys"]); keys != nil {
		fmt.CallerKey = util.AnyToString(keys["caller"])
		fmt.TimeKey = util.AnyToString(keys["time"])
		fmt.MessageKey = util.AnyToString(keys["message"])
		fmt.LevelKey = util.AnyToString(keys["level"])
		fmt.NameKey = util.AnyToString(keys["name"])
	}
	if timeFMT, ok := m["time_format"].(string); ok {
		fmt.TimeFormat = log.MakeTimeFormat(timeFMT)
	}
	return fmt
}

func parseMessageFormat(v interface{}) log.MessageFormat {
	name, _ := v.(string)
	return log.MakeMessageFormat(name)
}

func parseLevel(v interface{}) log.Level {
	name, _ := v.(string)
	return log.MakeLevelWithName(name)
}

func parseFileRotation(m util.AnyMap) log.FileRotation {
	return log.FileRotation{
		MaxSize:      int(util.AnyToInt64(m["max_size"])),
		Compress:     util.AnyToBool(m["compress"]),
		MaxAge:       int(util.AnyToInt64(m["max_age"])),
		MaxBackups:   int(util.AnyToInt64(m["max_backups"])),
		LocalTime:    util.AnyToBool(m["localtime"]),
		RotateOnTime: util.AnyToBool(m["rotate_on_time"]),
		RotatePeriod: util.AnyToString(m["rotate_period"]),
		RotateAfter:  int(util.AnyToInt64(m["rotate_after"])),
	}
}
