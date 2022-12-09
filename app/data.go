package app

import (
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

type data struct {
	CPU cpu.CPU
	Mem mem.Mem
}
