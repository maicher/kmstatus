//go:build !X

package out

import (
	"bytes"
	"fmt"
)

type Window struct {
}

func NewWindow() (*Window, error) {
	return &Window{}, fmt.Errorf("XWindow is not supported. kmstatus needs to be built with -tags X to support the -x option")
}

func (w *Window) SetStatus(_ *bytes.Buffer) {
}
