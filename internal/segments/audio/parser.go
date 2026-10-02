package audio

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/maicher/kmstatus/internal/segments/common"
)

type Parser struct {
}

func (p *Parser) Parse(data *data) error {
	var buf bytes.Buffer
	var r *strings.Reader

	err := common.RunCommand(&buf, "pamixer", "--get-mute", "--get-volume")
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

	buf.Reset()
	err = common.RunCommand(&buf, "pamixer", "--list-sources")
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
	buf.Reset()
	err = common.RunCommand(&buf, "pamixer", "--source", mic, "--get-mute", "--get-volume")
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
