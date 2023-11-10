package main

import (
	"fmt"
	"github.com/worming004/pmid/pkg/log"
	"github.com/worming004/pmid/pkg/podman"
	"log/slog"
)

func main() {
	closer := log.SetupDefaultLogger(true)
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
