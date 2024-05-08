package audio

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Parser struct {
}

func (p *Parser) Parse(data *data) error {
	var buf bytes.Buffer
	var r *strings.Reader

	cmd := exec.Command("pamixer", "--get-mute", "--get-volume")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err == nil {
		data.OutAvailable = true

		s := bufio.NewScanner(&buf)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			r = strings.NewReader(s.Text())
			fmt.Fscanf(r, "%t %d", &data.OutMuted, &data.OutVolume)
		}
	} else {
		data.OutAvailable = false
	}

	var mic string

	cmd = exec.Command("pamixer", "--list-sources")
	cmd.Stdout = &buf
	err = cmd.Run()
	if err == nil {
		s := bufio.NewScanner(&buf)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			if strings.Contains(s.Text(), "Microphone") {
				r = strings.NewReader(s.Text())
				fmt.Fscanf(r, "%s", &mic)
			}
		}
	}

	if mic == "" {
		data.InAvailable = false

		return nil
	}

	data.InAvailable = true
	cmd = exec.Command("pamixer", "--source", mic, "--get-mute", "--get-volume")
	cmd.Stdout = &buf
	err = cmd.Run()
	if err == nil {
		data.InAvailable = true

		s := bufio.NewScanner(&buf)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			r = strings.NewReader(s.Text())
			fmt.Fscanf(r, "%t %d", &data.InMuted, &data.InVolume)
		}
	}

	return nil
}
