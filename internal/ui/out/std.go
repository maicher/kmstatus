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
	buffer.WriteString("\n")
	os.Stdout.Write(buffer.Bytes())
}
