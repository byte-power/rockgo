package rock

import (
	"errors"
	"fmt"

	"github.com/byte-power/rockgo/log"
	"github.com/byte-power/rockgo/util"
)

var loggers = map[string]*log.Logger{}

// 用默认的logger防止调用Logger(name)时对应logger不存在而得到空指针
var defaultLogger log.Logger

// 取得在Application.Init初始化过的该名称对应的Logger
// 若无对应的Logger，返回默认Logger
func Logger(name string) *log.Logger {
	if logger, ok := loggers[name]; ok {
		return logger
	}
	return &defaultLogger
}

func parseLogger(appName, name string, cfg util.AnyMap) (*log.Logger, error) {
	var outputs []log.Output
	for k, v := range cfg {
		vs := util.AnyToAnyMap(v)
		if vs == nil && v != nil {
			return nil, fmt.Errorf("'log.%v' should be map", k)
		}
		switch k {
		case "console":
			it := parseConsoleLogger(name, vs)
			outputs = append(outputs, it)
		case "file":
			it := parseFileLogger(name, vs)
			outputs = append(outputs, it)
		case "fluent":
			it, _ := parseFluentLogger(vs)
			outputs = append(outputs, it)
		}
	}
	return log.NewLogger(outputs...), nil
}

func parseConsoleLogger(name string, cfg util.AnyMap) log.Output {
	fmt := parseFormat(cfg)
	level := parseLevel(cfg["level"])
	stream := log.MakeConsoleStream(util.AnyToString(cfg["stream"]))
	return log.MakeConsoleOutput(name, fmt, level, stream)
}

func parseFileLogger(name string, cfg util.AnyMap) log.Output {
	fmt := parseFormat(cfg)
	level := parseLevel(cfg["level"])
	location := util.AnyToString(cfg["location"])
	rotation := parseFileRotation(util.AnyToAnyMap(cfg["rotation"]))
	return log.MakeFileOutput(name, fmt, level, location, rotation)
}

func parseFluentLogger(cfg util.AnyMap) (log.Output, error) {
	level := parseLevel(cfg["level"])
	if cfg["host"] == nil {
		return nil, errors.New("fluent host must be specified in config.")
	}
	host := util.AnyToString(cfg["host"])
	port := int(24224)
	if cfg["port"] != nil {
		port = int(util.AnyToInt64(cfg["port"]))
	}
	if cfg["tag"] == nil {
		return nil, errors.New("fluent tag must be sepecified in config.")
	}
	tag := util.AnyToString(cfg["tag"])
	async := false
	if cfg["async"] != nil {
		async = util.AnyToBool(cfg["async"])
	}
	return log.MakeFluentOutput(level, host, port, tag, async), nil
}

func parseFormat(cfg util.AnyMap) log.LocalFormat {
	msgFMT := parseMessageFormat(cfg["format"])
	fmt := log.MakeLocalFormat(msgFMT)
	if keys := util.AnyToAnyMap(cfg["keys"]); keys != nil {
		fmt.CallerKey = util.AnyToString(keys["caller"])
		fmt.TimeKey = util.AnyToString(keys["time"])
		fmt.MessageKey = util.AnyToString(keys["message"])
		fmt.LevelKey = util.AnyToString(keys["level"])
		fmt.NameKey = util.AnyToString(keys["name"])
	}
	if timeFMT, ok := cfg["time_format"].(string); ok {
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

func parseFileRotation(cfg util.AnyMap) log.FileRotation {
	return log.FileRotation{
		MaxSize:      int(util.AnyToInt64(cfg["max_size"])),
		Compress:     util.AnyToBool(cfg["compress"]),
		MaxAge:       int(util.AnyToInt64(cfg["max_age"])),
		MaxBackups:   int(util.AnyToInt64(cfg["max_backups"])),
		LocalTime:    util.AnyToBool(cfg["localtime"]),
		RotateOnTime: util.AnyToBool(cfg["rotate_on_time"]),
		RotatePeriod: util.AnyToString(cfg["rotate_period"]),
		RotateAfter:  int(util.AnyToInt64(cfg["rotate_after"])),
	}
}
