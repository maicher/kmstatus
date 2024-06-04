package ipc

import (
	"fmt"
	"net"
)

func Send(cmd, socketPath string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("unable to send command to main process, because the main process is not running: %s", err)
	}

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("unable to send command to main process: %s", err)

	}

	return nil
}
