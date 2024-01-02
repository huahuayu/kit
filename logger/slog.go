package logger

import (
	"log/slog"
	"os"
	"time"
)

func DefaultSlog() *slog.Logger {
	var timeFormatter = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			// Replace time format to time.RFC3339
			return slog.Attr{Key: slog.TimeKey, Value: slog.StringValue(a.Value.Time().Format(time.RFC3339))}
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: timeFormatter}))
	slog.SetDefault(logger)
	return logger
}
