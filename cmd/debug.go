package main

import (
	"bubblepod/pkg/podman"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		slog.Error("cannot open log file: %v", err)
		panic(err)
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))
	defer f.Close()

	cmd := podman.PodmanCommands{}
	img, err := cmd.GetImages()
	if err != nil {
		slog.Error("cannot get podman images: %v", err)
		panic(err)
	}

	for _, i := range img {
		fmt.Printf("%v\n", i)
	}
}
