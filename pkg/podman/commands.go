package podman

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type Image struct {
	Name, Id string
}
type PodmanCommands struct {
}

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
			Id:   it[2],
		})

	}

	return result, nil
}
