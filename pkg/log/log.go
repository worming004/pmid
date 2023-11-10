package log

import (
	"io"
	"log/slog"
	"os"
)

func SetupDefaultLogger() io.Closer {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		slog.Error("cannot open log file: %v", err)
		panic(err)
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))
	return f
}
