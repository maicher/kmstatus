package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

const help = `NAME
  kmstatus

SYNOPSIS
  kmstatus [OPTION...]

OPTIONS
  --print-config, -C     print config and exit
  --print-template, -T   print the default template and exit
  --template-name, -t    name of the template to use (default: ` + DefaultTemplateName + `)

SIGNALS
  By default kmstatus will refresh every n seconds,
  where n is the lowest interval value from -template option.
  A signal can be sent to kmstatus to trigger an additional refresh:
    pkill -USR1 kmstatus$

ENVIRONMENT VARIABLES
  XDG_CONFIG_HOME
  HOME
`

const DefaultTemplateName = "default.tmpl"

type NewParserFunc func() (parsers.Parser, error)

type ParserConfig struct {
	Name          string
	NewParserFunc NewParserFunc
	Interval      time.Duration
	OnSig         bool
}

type Config struct {
	PrintConfig   bool
	PrintTemplate bool

	TemplateName  string
	ParserConfigs []ParserConfig
}

func (c *Config) setFlags(pc *ParserConfig, f NewParserFunc, name string, interval time.Duration, onSig bool) {
	pc.NewParserFunc = f
	pc.Name = name
	flag.DurationVar(&(pc.Interval), name, interval, "")
	flag.BoolVar(&(pc.OnSig), name+"sig", onSig, "")
}

// Parse parses the command line options to the Config struct.
func Parse() *Config {
	c := &Config{
		ParserConfigs: make([]ParserConfig, 4),
	}

	flag.BoolVar(&c.PrintConfig, "print-config", false, "")
	flag.BoolVar(&c.PrintConfig, "C", false, "")

	flag.BoolVar(&c.PrintTemplate, "print-template", false, "")
	flag.BoolVar(&c.PrintTemplate, "T", false, "")

	flag.StringVar(&c.TemplateName, "template-name", "default.tmpl", "")
	flag.StringVar(&c.TemplateName, "t", DefaultTemplateName, "")

	c.setFlags(&(c.ParserConfigs[0]), cpu.NewFreqParser, "cpu-freq", time.Second, false)
	c.setFlags(&(c.ParserConfigs[1]), cpu.NewLoadParser, "cpu-load", time.Second, false)
	c.setFlags(&(c.ParserConfigs[2]), cpu.NewTempParser, "cpu-temp", time.Second, false)
	c.setFlags(&(c.ParserConfigs[3]), mem.NewMemParser, "mem", time.Second, false)

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprintf(f, help) }
	flag.Parse()

	return c
}

func (c *Config) HasDefaultTemplate() bool {
	return c.TemplateName == DefaultTemplateName
}

func (c *Config) FindTemplatePath() (path string, err error) {
	if strings.HasPrefix(c.TemplateName, "/") {
		path = c.TemplateName
	} else {
		dir, present := os.LookupEnv("XDG_CONFIG_HOME")
		if !present {
			dir = os.Getenv("HOME")
		}

		path = filepath.Join(dir, "kmstatus", c.TemplateName)
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, fmt.Errorf("A path to template does not exist: %w", err)
	}

	return
}
