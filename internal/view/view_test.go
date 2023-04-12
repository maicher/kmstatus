package view

import (
	"strings"
	"testing"

	"github.io/maicher/kmstatus/pkg/parsers/cpu"
)

type Display struct {
	status string
}

func (d *Display) SetStatus(s string) {
	d.status = s
}

func Test_View(t *testing.T) {
	display := &Display{}
	v := View{
		templateName: "default.tmpl",
		templates:    mustCompileDefaultTemplate(),
		display:      display,
	}

	d := &Data{
		CPU: cpu.CPU{
			Freq: cpu.Freq(1000000),
		},
	}

	v.Render(d)

	if !strings.Contains(display.status, "1.0MHz") {
		t.Fatalf("Rendered status: %s does not contains phrase: '1MHz'", display.status)
	}
}
