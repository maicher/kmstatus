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
	data.OutAvailable, data.OutMuted, data.OutVolume = getVolume()

	mic := findMicrophone()
	if mic == "" {
		data.InAvailable, data.InMuted, data.InVolume = false, false, 0
		return nil
	}
	data.InAvailable, data.InMuted, data.InVolume = getVolume("--source", mic)

	return nil
}

// findMicrophone returns the index of the source with "Microphone" in its description.
// It is not necessarily the default source (which may be e.g. a webcam).
func findMicrophone() (mic string) {
	var buf bytes.Buffer

	err := common.RunCommand(&buf, "pamixer", "--list-sources")
	if err != nil {
		return ""
	}

	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if strings.Contains(s.Text(), "Microphone") {
			mic, _, _ = strings.Cut(s.Text(), " ")
		}
	}

	return mic
}

// getVolume reads the mute state and the volume of the default sink,
// or of a device selected with the additional pamixer args.
func getVolume(args ...string) (available, muted bool, volume int) {
	var buf bytes.Buffer

	args = append(args, "--get-mute", "--get-volume")
	err := common.RunCommand(&buf, "pamixer", args...)
	if err != nil {
		return false, false, 0
	}

	_, err = fmt.Fscanf(&buf, "%t %d", &muted, &volume)
	if err != nil {
		return false, false, 0
	}

	return true, muted, volume
}
