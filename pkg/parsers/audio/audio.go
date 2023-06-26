package audio

type X struct {
	IsAvailable bool
	IsMuted     bool
	Volume      int
}

type Audio struct {
	Sink   X
	Source X
}
