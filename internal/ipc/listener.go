package ipc

import (
	"encoding/json"
	"fmt"
	"net"
)

type refreshHandlerFunc = func()
type setTextHandlerFunc = func(string)
type unsetTextHandlerFunc = func()
type errorHandlerFunc = func(error)

type Listener struct {
	SocketPath string
	listener   net.Listener

	RefreshHandler   refreshHandlerFunc
	SetTextHandler   setTextHandlerFunc
	UnsetTextHandler unsetTextHandlerFunc
	ErrorHandler     errorHandlerFunc
}

func (l *Listener) Listen() {
	var (
		cmd Cmd
		err error
	)

	l.listener, err = net.Listen("unix", l.SocketPath)
	if err != nil {
		l.ErrorHandler(fmt.Errorf("error creating listener: %s", err))
		return
	}

	for {
		conn, err := l.listener.Accept()
		if err != nil {
			l.ErrorHandler(err)
			return
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			l.ErrorHandler(fmt.Errorf("error reading: %s", err))
			return
		}

		err = json.Unmarshal(buffer[:n], &cmd)
		if err != nil {
			l.ErrorHandler(err)
			return
		}
		switch cmd.Name {
		case Refresh:
			if l.RefreshHandler != nil {
				l.RefreshHandler()
			}
		case SetText:
			if l.SetTextHandler != nil {
				l.SetTextHandler(cmd.Payload)
			}
		case UnsetText:
			if l.UnsetTextHandler != nil {
				l.UnsetTextHandler()
			}
		default:
		}
	}
}

func (i *Listener) Close() {
	if i.listener == nil {
		return
	}

	// Close also removes the socket file.
	i.listener.Close()
}
