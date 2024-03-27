package main

import (
	"fmt"
	"time"

	_ "embed"

	"github.com/BurntSushi/toml"
)

//go:embed doc.txt
var doc string

//go:embed kmstrc.example.toml
var kmstrcExample string

type Parser struct {
	Name     string
	Interval time.Duration
	OnSig    bool
	Template string
}

type Config struct {
	Timefmt string
	Parsers []Parser `toml:"parser"`
}

func NewConfig(path string) (Config, error) {
	var c Config
	err := toml.Unmarshal([]byte(kmstrcExample), &c)
	if err != nil {
		return c, fmt.Errorf("Unable to parse config %s: %s", path, err)
	}

	return c, nil
}
