package gpu

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.io/maicher/kmstatus/pkg/parsers"
)

type GPUParser struct {
	path string
}

func NewGPUParser() (parsers.Parser, error) {
	parser := &GPUParser{}

	path, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return parser, fmt.Errorf("Processes parser: %s\n", err)
	}

	parser.path = path

	return parser, nil
}

func (p *GPUParser) Parse() (any, error) {
	var buf bytes.Buffer
	gpu := NewGPU()

	cmd := exec.Command(p.path, "--format=csv,nounits,noheader", "--query-gpu=power.draw,temperature.gpu,utilization.gpu,clocks.current.graphics,memory.used,memory.total")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return gpu, fmt.Errorf("GPU parser: %s\n", err)
	}

	_, err = fmt.Fscanf(&buf, "%f, %d, %d, %d, %d, %d", &gpu.Power, &gpu.Temp, &gpu.Load, &gpu.Freq, &gpu.MemUsed, &gpu.MemTotal)

	return gpu, err
}
