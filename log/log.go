package log

type Logger struct {
	outpers []Output
}

type Output interface {
	Level() Level
	Log(l Level, msg string, argPairs []interface{})
	LogMap(l Level, msg string, values map[string]interface{})
	LogPlainMessage(l Level, args []interface{})
	LogFormatted(l Level, format string, args []interface{})
}

func NewLogger(outpers ...Output) *Logger {
	return &Logger{outpers: outpers}
}

func MakeConsoleOutput(name string, fmt LocalFormat, level Level, stream ConsoleStream) Output {
	return newZapConsoleLogger(name, fmt, level, stream)
}

func MakeFileOutput(name string, fmt LocalFormat, level Level, location string, rotation FileRotation) Output {
	return newZapFileLogger(name, fmt, level, location, rotation)
}

func MakeFluentOutput(level Level, host string, port int, tag string, async bool) Output {
	fluent_logger := FluentOutput{level: level, host: host, port: port, tag: tag, async: async}
	fluent_logger.Init()
	return &fluent_logger
}

// 发送可结构化的消息
// argPairs被视为键值对，键或值为nil的将不被记录
// 本地log：直接发给zap.w处理
// fluent log：将在内部拼成map {"msg": msg, arg0: arg1, arg2: arg3, ...}，msg仅当参数msg非空串才会被记录
func (l *Logger) Log(level Level, msg string, argPairs ...interface{}) {
	for _, it := range l.outpers {
		if level >= it.Level() {
			it.Log(level, msg, argPairs)
		}
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Log(LevelDebug, msg, args...)
}
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Log(LevelInfo, msg, args...)
}
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Log(LevelWarn, msg, args...)
}
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Log(LevelError, msg, args...)
}
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Log(LevelFatal, msg, args...)
}

// 发送可结构化的消息
// 本地log：解构后发给zap.w处理
// fluent log：直接发送
func (l *Logger) LogMap(level Level, msg string, values map[string]interface{}) {
	for _, it := range l.outpers {
		if it.Level() >= level {
			continue
		}
		it.LogMap(level, msg, values)
	}
}

func (l *Logger) Debugm(msg string, values map[string]interface{}) {
	l.LogMap(LevelDebug, msg, values)
}
func (l *Logger) Infom(msg string, values map[string]interface{}) {
	l.LogMap(LevelInfo, msg, values)
}
func (l *Logger) Warnm(msg string, values map[string]interface{}) {
	l.LogMap(LevelWarn, msg, values)
}
func (l *Logger) Errorm(msg string, values map[string]interface{}) {
	l.LogMap(LevelError, msg, values)
}
func (l *Logger) Fatalm(msg string, values map[string]interface{}) {
	l.LogMap(LevelFatal, msg, values)
}

// 发送简单消息
// 本地log：直接发给zap处理
// fluent log：args如仅包含一个map，会用它发给fluent，若有多个将合并为信息字符串，并将map[string]interface{}{"m":信息字符串}发给fluent
func (l *Logger) LogPlainMessage(level Level, args ...interface{}) {
	for _, it := range l.outpers {
		if it.Level() >= level {
			continue
		}
		it.LogPlainMessage(level, args)
	}
}

// 发送格式化后的简单消息
// 内部逻辑与LogPlainMessage基本相同
func (l *Logger) LogFormatted(level Level, format string, args ...interface{}) {
	for _, it := range l.outpers {
		if it.Level() > level {
			continue
		}
		it.LogFormatted(level, format, args)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.LogFormatted(LevelDebug, format, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.LogFormatted(LevelInfo, format, args...)
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.LogFormatted(LevelWarn, format, args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.LogFormatted(LevelError, format, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.LogFormatted(LevelFatal, format, args...)
}
