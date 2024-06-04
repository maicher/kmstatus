package ipc

import (
	"fmt"
	"net"
)

type refreshFunc = func()
type setTextFunc = func(string)
type unsetTextFunc = func()
type errorFunc = func(error)

type IPC struct {
	SocketPath string
	listener   net.Listener

	RefreshFunc   refreshFunc
	SetTextFunc   setTextFunc
	UnsetTextFunc unsetTextFunc
	ErrorFunc     errorFunc
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

func (i *IPC) Listen() {
	var (
		cmd string
		err error
	)

	i.listener, err = net.Listen("unix", i.SocketPath)
	if err != nil {
		i.ErrorFunc(fmt.Errorf("error creating listener: %s", err))
		return
	}

	for {
		conn, err := i.listener.Accept()
		if err != nil {
			i.ErrorFunc(err)
			return
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			i.ErrorFunc(fmt.Errorf("error reading: %s", err))
			return
		}

		cmd = string(buffer[:n])
		switch cmd {
		case "cmd:refresh":
			if i.RefreshFunc != nil {
				i.RefreshFunc()
			}
		case "cmd:unsetText":
			if i.UnsetTextFunc != nil {
				i.UnsetTextFunc()
			}
		default:
			if i.SetTextFunc != nil {
				i.SetTextFunc(cmd)
			}
		}
	}
}

func (i *IPC) CloseListener() {
	if i.listener == nil {
		return
	}

	// Close also removes the socket file.
	i.listener.Close()
}
