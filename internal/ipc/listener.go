package ipc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// Upper limit for a single command message.
const maxMsgSize = 64 * 1024

// Time given to a client to send its command.
const readTimeout = 2 * time.Second

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

// Open creates the socket. A stale socket file left by a process that
// did not exit cleanly is removed. If another instance is listening on
// the socket, an error is returned.
func (l *Listener) Open() error {
	err := removeStaleSocket(l.SocketPath)
	if err != nil {
		return err
	}

	l.listener, err = net.Listen("unix", l.SocketPath)
	if err != nil {
		return fmt.Errorf("error creating listener: %s", err)
	}

	// Only the owner is allowed to send commands.
	err = os.Chmod(l.SocketPath, 0600)
	if err != nil {
		l.listener.Close()
		return fmt.Errorf("error setting socket permissions: %s", err)
	}

	return nil
}

// Listen accepts connections until the listener is closed.
// Invalid commands are reported to stderr and do not stop listening.
func (l *Listener) Listen() {
	for {
		conn, err := l.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			l.ErrorHandler(err)
			return
		}

		err = l.handle(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kmstatus: invalid control command: %s\n", err)
		}
	}
}

func (l *Listener) handle(conn net.Conn) error {
	var cmd Cmd

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(readTimeout))
	data, err := io.ReadAll(io.LimitReader(conn, maxMsgSize+1))
	if err != nil {
		return fmt.Errorf("error reading: %s", err)
	}
	if len(data) > maxMsgSize {
		return fmt.Errorf("message exceeds %d bytes", maxMsgSize)
	}

	// A connection without data is a probe whether the instance is running.
	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, &cmd)
	if err != nil {
		return err
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
		return fmt.Errorf("unknown command %q", cmd.Name)
	}

	return nil
}

func (i *Listener) Close() {
	if i.listener == nil {
		return
	}

	// Close also removes the socket file.
	i.listener.Close()
}

func removeStaleSocket(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to check socket %s: %s", path, err)
	}
	if info.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("%s exists and is not a socket", path)
	}

	conn, err := net.Dial("unix", path)
	if err == nil {
		conn.Close()
		return fmt.Errorf("another kmstatus instance is already listening on %s", path)
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("unable to remove stale socket %s: %s", path, err)
	}

	return nil
}
