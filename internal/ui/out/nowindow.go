//go:build !X

package out

import "fmt"

type Window struct {
}

func NewWindow() (*Window, error) {
	return &Window{}, fmt.Errorf("XWindow is not supported. kmst needs to be build with -tag X to support the -x option")
}

func (w *Window) SetStatus(name string) {
}
