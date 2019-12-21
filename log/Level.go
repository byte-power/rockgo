package log

type Level int

const (
	LevelDebug = 1 << iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func MakeLevelWithName(name string) Level {
	switch name {
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
