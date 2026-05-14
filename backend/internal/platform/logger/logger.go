package logger

import (
	"log/slog"
	"os"
)

// New creates a structured logger based on environment.
func New(env string) *slog.Logger {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	return slog.New(handler)
}

// Redacted returns an attribute with its value hidden.
func Redacted(key string) slog.Attr {
	return slog.String(key, "[REDACTED]")
}
