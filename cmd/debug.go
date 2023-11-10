package main

import (
	"bubblepod/pkg/log"
	"bubblepod/pkg/podman"
	"fmt"
	"log/slog"
)

func main() {
	closer := log.SetupDefaultLogger()
	defer closer.Close()

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
