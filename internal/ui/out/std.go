package out

import "fmt"

type Std struct {
}

func NewStd() (*Std, error) {
	return &Std{}, nil
}

func (s *Std) SetStatus(name string) {
	fmt.Println(name)
}
