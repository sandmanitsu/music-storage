package logger

import (
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type Logger struct {
	Sl *slog.Logger
}

// Create logger
func NewLogger(env string) *Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("logger error: env is empty or incorrect value")
	}

	return &Logger{Sl: logger}
}

// Wrap error to slog.Attr struct
func (l *Logger) Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func (l *Logger) Info(msg string, arg ...any) {
	l.Sl.Info(msg, arg...)
}

func (l *Logger) Debug(msg string, arg ...any) {
	l.Sl.Debug(msg, arg...)
}

func (l *Logger) Error(msg string, arg ...any) {
	l.Sl.Debug(msg, arg...)
}

func (l *Logger) InfoAPI(msg string, code int, path string, err string) {
	l.Sl.Info(
		msg,
		slog.Attr{Key: "code", Value: slog.IntValue(code)},
		slog.Attr{Key: "path", Value: slog.StringValue(msg)},
		slog.Attr{Key: "error", Value: slog.StringValue(err)},
	)
}
