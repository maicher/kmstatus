package view

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers/cpu"
)

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"round":   round,
	"human":   human,
	"humanSI": humanSI,
	"ljust":   ljust,
	"clock":   clock,
}

// Prepends a string with spaces so that the length of the output string is num.
func ljust(num int, s string) string {
	l := num - len(s)
	if l < 0 {
		return s
	}

	return strings.Repeat(" ", l) + s
}

func clock(format string) string {
	return time.Now().Format(format)
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

	return fmt.Sprintf("%.*f%c", precision, float64(b)/float64(div), "kMGT"[exp])
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
