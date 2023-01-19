package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

const help = `NAME
  kmstatus

SYNOPSIS
  kmstatus [OPTION...]

OPTIONS
  --print-config, -C     print config and exit
  --print-template, -T   print the default template and exit

  --xwindow, -x          print output to default's window WM_NAME (instead stdout)
                         (to use this option kmstatus needs to be build with -tag X)
  --template-name, -t    name of the template to use (default: ` + DefaultTemplateName + `)

  --[parser-name] INTERVAL
  --[parser-name]-sig

PARSER NAMES
  cpu-freq
  cpu-load
  cpu-temp
  mem
  fs

CONFIG
  Config options can be put in the $HOME/.config/kmstatus/` + ConfigFileName + ` file.
  Put each option in a new line as a space-separated pair.

SIGNALS
  By default kmstatus will refresh every n seconds,
  where n is the lowest interval value from parser options.
  A signal can be sent to kmstatus to trigger an additional refresh:
    pkill -USR1 kmstatus$

ENVIRONMENT VARIABLES
  XDG_CONFIG_HOME
  HOME
`

const (
	DefaultTemplateName = "default.tmpl"
	ConfigFileName      = "kmstatusrc"
)

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

	XWindow bool
}

// Parse parses the command line options to the Config struct.
func Parse() *Config {
	c := &Config{
		ParserConfigs: make([]ParserConfig, 5),
	}

	flag.BoolVar(&c.PrintConfig, "print-config", false, "")
	flag.BoolVar(&c.PrintConfig, "C", false, "")

	flag.BoolVar(&c.PrintTemplate, "print-template", false, "")
	flag.BoolVar(&c.PrintTemplate, "T", false, "")

	flag.StringVar(&c.TemplateName, "template-name", "default.tmpl", "")
	flag.StringVar(&c.TemplateName, "t", DefaultTemplateName, "")

	flag.BoolVar(&c.XWindow, "xwindow", false, "")
	flag.BoolVar(&c.XWindow, "x", false, "")

	c.setFlags(&(c.ParserConfigs[0]), cpu.NewFreqParser, "cpu-freq", time.Second, false)
	c.setFlags(&(c.ParserConfigs[1]), cpu.NewLoadParser, "cpu-load", time.Second, false)
	c.setFlags(&(c.ParserConfigs[2]), cpu.NewTempParser, "cpu-temp", time.Second, false)
	c.setFlags(&(c.ParserConfigs[3]), mem.NewMemParser, "mem", time.Second, false)
	c.setFlags(&(c.ParserConfigs[4]), filesystem.NewFSParser, "fs", time.Second, false)

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprintf(f, help) }

	// Parse options from the config file.
	tmp := os.Args
	options, ok := optionsFromFile()
	if ok {
		os.Args = options
		flag.Parse()
	}

	// Then parse options from the command line.
	os.Args = tmp
	flag.Parse()

	return c
}

func (c *Config) setFlags(pc *ParserConfig, f NewParserFunc, name string, interval time.Duration, onSig bool) {
	pc.NewParserFunc = f
	pc.Name = name
	flag.DurationVar(&(pc.Interval), name, interval, "")
	flag.BoolVar(&(pc.OnSig), name+"-sig", onSig, "")
}

func (c *Config) HasDefaultTemplate() bool {
	return c.TemplateName == DefaultTemplateName
}

func (c *Config) FindTemplatePath() (path string, err error) {
	if strings.HasPrefix(c.TemplateName, "/") {
		path = c.TemplateName
	} else {
		path = filepath.Join(configDir(), c.TemplateName)
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, fmt.Errorf("A path to template does not exist: %w", err)
	}

	return
}

func optionsFromFile() ([]string, bool) {
	ret := []string{}
	path := filepath.Join(configDir(), ConfigFileName)

	file, err := os.Open(path)
	if err != nil {
		return ret, false
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		ret = append(ret, "-"+strings.Replace(s.Text(), " ", "=", 1))
	}

	return ret, true
}

func configDir() string {
	dir, present := os.LookupEnv("XDG_CONFIG_HOME")
	if !present {
		dir = os.Getenv("HOME")
	}

	return filepath.Join(dir, "kmstatus")
}
