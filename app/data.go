package app

import (
	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

type Data struct {
	CPU cpu.CPU
	Mem mem.Mem
}
