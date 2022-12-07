package app

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

type System struct {
	CPU cpu.CPU
	Mem mem.Mem
}

type View struct {
	tmplName string
	tmpl     *template.Template
}

func (v *View) Render(s *System) string {
	b := bytes.Buffer{}
	err := v.tmpl.ExecuteTemplate(&b, v.tmplName, s)
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(b.String(), "\n", "")
}

func NewView(t string) *View {
	tmpl, err := template.New(t).Funcs(template.FuncMap{
		"round":   round,
		"human":   human,
		"humanSI": humanSI,
	}).ParseFiles("templates/" + t)
	if err != nil {
		panic(err)
	}

	return &View{
		tmplName: t,
		tmpl:     tmpl,
	}
}

// Converts a number to a human-readable string.
// 1M == 1024k
func human(precision int, v any) string {
	return humanUnit(1024, precision, v)
}

// Converts a number to a human-readable string.
// 1M == 1000k
func humanSI(precision int, v any) string {
	return humanUnit(1000, precision, v)
}

func humanUnit(unit int, precision int, v any) string {
	var b int
	switch v.(type) {
	case cpu.Freq:
		b = int(v.(cpu.Freq))
	case int:
		b = v.(int)
	default:
		panic("Unknown type")
	}

	if b < unit {
		return strconv.Itoa(b)
	}

	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.*f%c", precision, float64(b)/float64(div), "MGTk"[exp])
}

// Rounds a number according to given precision.
// round(0, 1.23)
// => 1
// round(1, 1.23)
// => 1.2
func round(precision int, v any) string {
	var f float64
	switch v.(type) {
	case cpu.Load:
		f = float64(v.(cpu.Load))
	case float64:
		f = v.(float64)
	default:
		panic("Unknown type")
	}
	return fmt.Sprintf("%.*f", precision, f)
}
