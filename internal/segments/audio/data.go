package audio

type data struct {
	InAvailable bool
	InVolume    int
	InMuted     bool

	OutAvailable bool
	OutVolume    int
	OutMuted     bool
}
