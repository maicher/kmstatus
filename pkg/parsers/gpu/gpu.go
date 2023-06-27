package gpu

// Freq holds CPU frequency in kHz.
type Freq int

type SpaceMB int

type GPU struct {
	Power    float64
	Temp     int
	Load     int
	Freq     Freq
	MemUsed  SpaceMB
	MemTotal SpaceMB
}

func NewGPU() GPU {
	return GPU{}
}

func (gpu *GPU) MemUsedPercentage() float64 {
	return 100 * float64(gpu.MemUsed) / float64(gpu.MemTotal)
}
