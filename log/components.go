package log

import "os"

import "strings"

type Level int

const (
	LevelDebug Level = 1 << iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func MakeLevelWithName(name string) Level {
	switch strings.ToLower(name) {
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelDebug
	}
}

// -------------------------------

type MessageFormat string

const (
	MessageFormatJSON MessageFormat = "json"
	MessageFormatText MessageFormat = "text"
)

// MakeMessageFormat would product MessageFormat with raw string.
// MessageFormatText would be default returning if no matched.
func MakeMessageFormat(raw string) MessageFormat {
	switch strings.ToLower(raw) {
	case string(MessageFormatJSON):
		return MessageFormatJSON
	default:
		return MessageFormatText
	}
}

func (f MessageFormat) isJSON() bool {
	return f == MessageFormatJSON
}

// -------------------------------

type TimeFormat string

const (
	TimeFormatISO8601 TimeFormat = "iso8601"
	TimeFormatSeconds TimeFormat = "seconds"
	TimeFormatMillis  TimeFormat = "millis"
	TimeFormatNanos   TimeFormat = "nanos"
)

func MakeTimeFormat(raw string) TimeFormat {
	switch strings.ToLower(raw) {
	case string(TimeFormatSeconds):
		return TimeFormatSeconds
	case string(TimeFormatMillis):
		return TimeFormatMillis
	case string(TimeFormatNanos):
		return TimeFormatNanos
	default:
		return TimeFormatISO8601
	}
}

// -------------------------------

type LocalFormat struct {
	Format MessageFormat

	MessageKey string // 默认为M
	TimeKey    string // 默认为T
	LevelKey   string // 默认为L
	NameKey    string // 默认为N
	CallerKey  string // 默认为C
	// 时间格式
	TimeFormat TimeFormat // 默认为TimeFormatISO8601
}

func MakeLocalFormat(msg MessageFormat) LocalFormat {
	return LocalFormat{
		Format:     msg,
		MessageKey: "M",
		TimeKey:    "T",
		LevelKey:   "L",
		NameKey:    "N",
		CallerKey:  "C",
		TimeFormat: TimeFormatISO8601,
	}
}

// -------------------------------

type ConsoleStream string

const (
	ConsoleStreamStdout ConsoleStream = "stdout"
	ConsoleStreamStderr ConsoleStream = "stderr"
)

func MakeConsoleStream(raw string) ConsoleStream {
	switch strings.ToLower(raw) {
	case string(ConsoleStreamStderr):
		return ConsoleStreamStderr
	default:
		return ConsoleStreamStdout
	}
}

func (s ConsoleStream) stream() *os.File {
	switch s {
	case ConsoleStreamStderr:
		return os.Stderr
	default:
		return os.Stdout
	}
}

// -------------------------------

// Rotation stores configs for the log rotation.
// See more in https://github.com/natefinch/lumberjack/tree/v2.0
type FileRotation struct {
	MaxSize    int
	Compress   bool
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	// RotateOnTime enables log rotation based on time.
	RotateOnTime bool
	// RotatePeriod is the period for log rotation.
	// Supports daily(d), hourly(h), minute(m) and second(s).
	RotatePeriod string
	// RotateAfter sets a value for time based rotation.
	// Log file rotates every RotateAfter * RotatePeriod.
	RotateAfter int
}
