package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	if path != "" {
		return parseConfig(path)
	}

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		path = filepath.Join(dir, "kmst/kmstrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}
	}

	if dir, ok := os.LookupEnv("HOME"); ok {
		path = filepath.Join(dir, ".config/kmst/kmstrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}

	}

	return parseDefaultConfig()
}

func parseDefaultConfig() (Config, error) {
	var c Config

	err := toml.Unmarshal([]byte(kmstrcExample), &c)
	if err != nil {
		return c, fmt.Errorf("Unable to parse default config file: %s", err)
	}

	return c, nil
}

func parseConfig(path string) (Config, error) {
	var c Config
	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, fmt.Errorf("Invalid path to the config file")
	} else if err != nil {
		return c, fmt.Errorf("Unable to read config file %s: %s", path, err)
	}

	err = toml.Unmarshal(bytes, &c)
	if err != nil {
		return c, fmt.Errorf("Unable to parse default config file: %s", err)
	}

	return c, nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
