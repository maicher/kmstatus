package audio

import "github.io/maicher/kmstatus/pkg/parsers"

type AudioParser struct {
}

func (p *AudioParser) Parse() (any, error) {
	a := Audio{}

	return a, nil
}

func NewAudioParser() (parsers.Parser, error) {
	parser := AudioParser{}

	return &parser, nil
}
