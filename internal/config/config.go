package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/maicher/kmstatus/internal/segments"
)

//go:embed kmstatusrc.example.toml
var kmstatusrcExample string

type Config struct {
	Segments []segments.Config `toml:"segment"`
}

func (c *Config) MinInterval() time.Duration {
	interval := time.Minute

	for _, v := range c.Segments {
		if v.RefreshInterval > 0 {
			if v.RefreshInterval < interval {
				interval = v.RefreshInterval
			}
		}
	}

	return interval
}

func New(path string) (Config, error) {
	if path != "" {
		return parseConfig(path)
	}

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		path = filepath.Join(dir, "kmstatus/kmstatusrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}
	}

	if dir, ok := os.LookupEnv("HOME"); ok {
		path = filepath.Join(dir, ".config/kmstatus/kmstatusrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}

	}

	return parseDefaultConfig()
}

func parseDefaultConfig() (Config, error) {
	var c Config

	err := toml.Unmarshal([]byte(kmstatusrcExample), &c)
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
