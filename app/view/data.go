package view

import (
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

type Data struct {
	CPU cpu.CPU
	Mem mem.Mem
	FS  filesystem.FS
}
