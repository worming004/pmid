package podman

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type Image struct {
	Name, Id, Tag string
}
type PodmanCommands struct {
}

var Default = PodmanCommands{}

func (pc PodmanCommands) GetImages() ([]Image, error) {
	cmd := exec.Command("podman", "images")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(output)
	return parseGetImage(reader)
}

func parseGetImage(reader io.Reader) ([]Image, error) {

	outputScanner := bufio.NewScanner(reader)

	outputScanner.Split(bufio.ScanLines)
	// Ignore header line
	if ok := outputScanner.Scan(); !ok {
		return nil, errors.New("scanning first line failed")
	}

	var result []Image
	for outputScanner.Scan() {
		t := outputScanner.Text()
		it := strings.Fields(t)
		slog.Debug(t, "length", len(it))

		result = append(result, Image{
			Name: it[0],
			Tag:  it[1],
			Id:   it[2],
		})

	}

	return result, nil
}

func (pm PodmanCommands) DeleteImageById(id string) error {
	slog.Debug("deleting image with id", "id", id)
	cmd := exec.Command("podman", "image", "rm", id)
	return cmd.Run()
}

func (pm PodmanCommands) DeleteImageByTag(repo, tag string) error {
	fullTag := fmt.Sprintf("%s:%s", repo, tag)
	slog.Debug("deleting image with tag", "tag", fullTag)
	cmd := exec.Command("podman", "image", "rm", fullTag)
	return cmd.Run()
}
