package bluetooth

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type BluetoothParser struct {
}

func NewBluetoothParser() (*BluetoothParser, error) {
	var p BluetoothParser

	return &p, nil
}

func (p *BluetoothParser) Parse(data *Data) error {
	err := exec.Command("sh", "-c", "systemctl is-active --quiet bluetooth").Run()
	if err == nil {
		data.IsServiceActive = true
	} else {
		data.IsServiceActive = false
	}

	var buf bytes.Buffer
	cmd := exec.Command("bluetoothctl", "show")
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		data.IsControllerPowered = false
		return nil
	}

	data.IsControllerPowered = strings.Contains(buf.String(), "Powered: yes")

	buf.Reset()
	cmd = exec.Command("bluetoothctl", "info")
	cmd.Stdout = &buf
	err = cmd.Run()
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
