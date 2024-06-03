package ipc

import (
	"fmt"
	"net"
)

type IPC struct {
	SocketPath string
	listener   net.Listener
}

func (i *IPC) Send(text string) error {
	conn, err := net.Dial("unix", i.SocketPath)
	if err != nil {
		return fmt.Errorf("unable to send command to main process, because the main process is not running: %s", err)
	}

	_, err = conn.Write([]byte(text))
	if err != nil {
		return fmt.Errorf("unable to send command to main process: %s", err)

	}

	return nil
}

func (i *IPC) Listen(f func(string)) error {
	var err error

	i.listener, err = net.Listen("unix", i.SocketPath)
	if err != nil {
		return fmt.Errorf("error creating listener: %s", err)
	}

	for {
		conn, err := i.listener.Accept()
		if err != nil {
			return nil
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return fmt.Errorf("error reading: %s", err)
		}

		f(string(buffer[:n]))
	}
}

func (i *IPC) CloseListener() {
	if i.listener == nil {
		return
	}

	// Close also removes the socket file.
	i.listener.Close()
}
