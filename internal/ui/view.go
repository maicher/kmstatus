package ui

import (
	"bytes"

	"github.com/maicher/kmst/internal/ui/out"
)

type statusSetter interface {
	SetStatus(*bytes.Buffer)
}

type View struct {
	display statusSetter
}

func NewView(isX bool) (*View, error) {
	var (
		v   View
		err error
	)

	if isX {
		v.display, err = out.NewWindow()
	} else {
		v.display, err = out.NewStd()
	}
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (v *View) Flush(buffer *bytes.Buffer) {
	v.display.SetStatus(buffer)
	buffer.Reset()
}
