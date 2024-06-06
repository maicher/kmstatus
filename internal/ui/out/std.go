package out

import (
	"bytes"
	"os"
)

type Std struct {
}

func NewStd() (*Std, error) {
	return &Std{}, nil
}

func (s *Std) SetStatus(buffer *bytes.Buffer) {
	os.Stdout.Write(append(buffer.Bytes(), '\n'))
}
