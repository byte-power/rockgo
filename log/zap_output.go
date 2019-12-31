package log

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const callerSkip = 3

var _ Output = (*zapOutput)(nil)

type zapOutput struct {
	level  Level
	output *zap.SugaredLogger
}

func (o *zapOutput) Level() Level {
	return o.level
}

func (o *zapOutput) Log(l Level, msg string, argPairs []interface{}) {
	switch l {
	case LevelDebug:
		o.output.Debugw(msg, argPairs...)
	case LevelInfo:
		o.output.Infow(msg, argPairs...)
	case LevelWarn:
		o.output.Warnw(msg, argPairs...)
	case LevelError:
		o.output.Errorw(msg, argPairs...)
	case LevelFatal:
		o.output.Fatalw(msg, argPairs...)
	}
}

func (o *zapOutput) LogMap(l Level, msg string, values map[string]interface{}) {
	var args []interface{}
	if count := len(values); count > 0 {
		args = make([]interface{}, 0, count*2)
		for k, v := range values {
			args = append(args, k, v)
		}
	}
	o.Log(l, msg, args)
}

func (o *zapOutput) LogPlainMessage(l Level, args []interface{}) {
	switch l {
	case LevelDebug:
		o.output.Debug(args...)
	case LevelInfo:
		o.output.Info(args...)
	case LevelWarn:
		o.output.Warn(args...)
	case LevelError:
		o.output.Error(args...)
	case LevelFatal:
		o.output.Fatal(args...)
	}
}

func (o *zapOutput) LogFormatted(l Level, format string, args []interface{}) {
	switch l {
	case LevelDebug:
		o.output.Debugf(format, args...)
	case LevelInfo:
		o.output.Infof(format, args...)
	case LevelWarn:
		o.output.Warnf(format, args...)
	case LevelError:
		o.output.Errorf(format, args...)
	case LevelFatal:
		o.output.Fatalf(format, args...)
	}
}

func newZapConsoleLogger(name string, fmt LocalFormat, level Level, stream zapcore.WriteSyncer) *zapOutput {
	writer := zapcore.Lock(stream)
	encoder := makeZapEncoder(fmt.Format.isJSON(), makeZapEncoderConfig(fmt))
	core := zapcore.NewCore(encoder, writer, makeZapLevel(level))
	output := zap.New(core,
		zap.AddCallerSkip(callerSkip),
		zap.AddCaller(),
	).Sugar()
	if name != "" {
		output = output.Named(name)
	}
	return &zapOutput{level: level, output: output}
}

func newZapFileLogger(name string, fmt LocalFormat, level Level, location string, rotation FileRotation) *zapOutput {
	fileLogger := lumberjack.Logger{
		Filename:   location,
		MaxSize:    rotation.MaxSize,
		Compress:   rotation.Compress,
		MaxAge:     rotation.MaxAge,
		MaxBackups: rotation.MaxBackups,
	}
	if rotation.RotateOnTime {
		go func() {
			interval := parseRotationPeriod(rotation.RotatePeriod, rotation.RotateAfter)
			for {
				<-time.After(interval)
				fileLogger.Rotate()
			}
		}()
	}
	writer := zapcore.AddSync(&fileLogger)
	encoder := makeZapEncoder(fmt.Format.isJSON(), makeZapEncoderConfig(fmt))
	core := zapcore.NewCore(encoder, writer, makeZapLevel(level))
	output := zap.New(core,
		zap.AddCallerSkip(callerSkip),
		zap.AddCaller(),
	).Sugar()
	if name != "" {
		output = output.Named(name)
	}
	return &zapOutput{level: level, output: output}
}

func makeZapEncoderConfig(f LocalFormat) zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.CallerKey = f.CallerKey
	cfg.LevelKey = f.LevelKey
	cfg.MessageKey = f.MessageKey
	cfg.NameKey = f.NameKey
	cfg.TimeKey = f.TimeKey
	if f.TimeFormat == "" {
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg.EncodeTime.UnmarshalText([]byte(f.TimeFormat))
	}
	return cfg
}

func makeZapEncoder(isJSON bool, encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func makeZapLevel(l Level) zapcore.Level {
	switch l {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func parseRotationPeriod(period string, n int) time.Duration {
	switch strings.ToLower(period) {
	case "day", "daily", "d":
		return time.Hour * 24 * time.Duration(n)
	case "hour", "hourly", "h":
		return time.Hour * time.Duration(n)
	case "minute", "m":
		return time.Minute * time.Duration(n)
	case "second", "s":
		return time.Second * time.Duration(n)
	default:
		return time.Hour * 24
	}
}
