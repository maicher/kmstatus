package ipc

import (
	"encoding/json"
	"fmt"
	"net"
)

func Send(cmd *Cmd, socketPath string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("unable to send command to main process, because the main process is not running: %s", err)
	}

	jsonData, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("unable to send command to main process, because the command is invalid: %s", err)
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("unable to send command to main process: %s", err)

	}

	return nil
}
