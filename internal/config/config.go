package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/maicher/kmstatus/internal/segments"
)

// Used for segments without the refreshinterval option.
const DefaultRefreshInterval = time.Second

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

func New(path, kmstatusrcExample string) (Config, error) {
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

	return parseDefaultConfig(kmstatusrcExample)
}

func parseDefaultConfig(kmstatusrcExample string) (Config, error) {
	c, err := decode(kmstatusrcExample)
	if err != nil {
		return c, fmt.Errorf("unable to parse default config file: %s", err)
	}

	return c, nil
}

func parseConfig(path string) (Config, error) {
	var c Config
	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, fmt.Errorf("invalid path to the config file")
	} else if err != nil {
		return c, fmt.Errorf("unable to read config file %s: %s", path, err)
	}

	c, err = decode(string(bytes))
	if err != nil {
		return c, fmt.Errorf("unable to parse config file %s: %s", path, err)
	}

	return c, nil
}

func decode(data string) (Config, error) {
	var c Config

	md, err := toml.Decode(data, &c)
	if err != nil {
		return c, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return c, fmt.Errorf("unknown option %q", undecoded[0].String())
	}

	for i := range c.Segments {
		if c.Segments[i].RefreshInterval < 0 {
			return c, fmt.Errorf("segment %d (%s): refreshinterval must be positive", i+1, c.Segments[i].ParserName)
		}

		if c.Segments[i].RefreshInterval == 0 {
			c.Segments[i].RefreshInterval = DefaultRefreshInterval
		}
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
