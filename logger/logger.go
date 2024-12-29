package logger

import (
	"log/slog"
	"os"
)

var LogLevel slog.LevelVar
var Log *slog.Logger

func init() {

	SetLogLevel("info") //info is default log level

	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &LogLevel}))
	slog.SetDefault(Log)
}

func SetLogLevel(level string) {
	switch level {
	case "debug":
		LogLevel.Set(slog.LevelDebug)
	case "info":
		LogLevel.Set(slog.LevelInfo)
	case "warn":
		LogLevel.Set(slog.LevelWarn)
	case "error":
		LogLevel.Set(slog.LevelError)
	default:
		panic("invalid log level : " + level)
	}
}
