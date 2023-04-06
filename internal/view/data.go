package view

import (
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
	"github.io/maicher/kmstatus/pkg/parsers/processes"
)

type Data struct {
	PS  processes.PS
	CPU cpu.CPU
	Mem mem.Mem
	FS  filesystem.FS
}
