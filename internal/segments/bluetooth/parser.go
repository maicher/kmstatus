package bluetooth

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
	err := common.RunCommand(nil, "systemctl", "is-active", "--quiet", "bluetooth")
	data.IsServiceActive = err == nil

	// bluetoothctl waits for the bluetoothd indefinitely, do not run it when the service is down.
	if !data.IsServiceActive {
		data.IsControllerPowered = false
		data.DeviceType = ""
		return nil
	}

	var buf bytes.Buffer
	err = common.RunCommand(&buf, "bluetoothctl", "show")
	if err != nil {
		data.IsControllerPowered = false
		return nil
	}

	data.IsControllerPowered = strings.Contains(buf.String(), "Powered: yes")

	buf.Reset()
	err = common.RunCommand(&buf, "bluetoothctl", "info")
	if err != nil {
		data.DeviceType = ""
		return nil
	}

	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if strings.Contains(s.Text(), "Icon:") {
			var ignored string

			r := strings.NewReader(s.Text())
			fmt.Fscanf(r, "%s %s", &ignored, &(data.DeviceType))
		}
	}

	return nil
}
