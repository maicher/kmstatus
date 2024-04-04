package ipc

import (
	"fmt"
	"net"
	"os"
)

type IPC struct {
	SocketPath string
	listener   net.Listener
}

func (i *IPC) Send(text string) error {
	conn, err := net.Dial("unix", i.SocketPath)
	if err != nil {
		return fmt.Errorf("Main process is not running: Error connecting to socket: %s", err)
	}

	_, err = conn.Write([]byte(text))
	if err != nil {
		return fmt.Errorf("Error sending message: %s", err)

	}

	return nil
}

func (i *IPC) Listen(f func(string)) error {
	var err error

	i.listener, err = net.Listen("unix", i.SocketPath)
	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
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
			return fmt.Errorf("Error reading: %s", err)
		}

		f(string(buffer[:n]))
	}
}

func (i *IPC) Close() {
	if i.listener != nil {
		i.listener.Close()
		os.Remove(i.SocketPath)
	}
}
