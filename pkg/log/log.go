package log

import (
	"io"
	"log/slog"
	"os"
)

func SetupDefaultLogger(tofile bool) io.Closer {
	if tofile {

		f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			slog.Error("cannot open log file: %v", err)
			panic(err)
		}

		handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))
		return f
	}

	return noClose{}
}

type noClose struct{}

func (nc noClose) Close() error {
	return nil
}
