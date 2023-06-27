package view

import (
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/gpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
	"github.io/maicher/kmstatus/pkg/parsers/network"
	"github.io/maicher/kmstatus/pkg/parsers/processes"
)

type Data struct {
	PS  processes.PS
	CPU cpu.CPU
	GPU gpu.GPU
	Mem mem.Mem
	FS  filesystem.FS
	Net network.Net
}
