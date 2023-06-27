package view

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/gpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"round":       round,
	"human":       human,
	"humanK":      humanK,
	"humanSI":     humanSI,
	"humanKSI":    humanKSI,
	"ljust":       ljust,
	"clock":       clock,
	"lastSegment": lastSegment,
}

// Prepends a string of int with spaces so that the length of the output string is num.
func ljust(num int, x any) string {
	var s string

	switch x.(type) {
	case string:
		s = x.(string)
	case int:
		s = strconv.Itoa(x.(int))
	}

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
	return humanUnit(1024, precision, v, "kMGT")
}

// Converts a number to a human-readable string.
// 1M == 1024k
func humanK(precision int, v any) string {
	return humanUnit(1024, precision, v, "MGT")
}

// Converts a number to a human-readable string.
// 1M == 1000k
func humanSI(precision int, v any) string {
	return humanUnit(1000, precision, v, "kMGT")
}

// Converts a number to a human-readable string.
// 1M == 1000k
func humanKSI(precision int, v any) string {
	return humanUnit(1000, precision, v, "MGT")
}

func humanUnit(unit int, precision int, v any, x string) string {
	var b int

	switch v.(type) {
	case cpu.Freq:
		b = int(v.(cpu.Freq))
	case gpu.Freq:
		b = int(v.(gpu.Freq)) * 1000000
	case mem.SpaceKB:
		b = int(v.(mem.SpaceKB))
	case gpu.SpaceMB:
		b = int(v.(gpu.SpaceMB)) * 1024 * 1024
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

	return fmt.Sprintf("%.*f%c", precision, float64(b)/float64(div), x[exp])
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

func lastSegment(name string) string {
	index := strings.LastIndex(name, "/")
	if index <= 0 {
		return name
	}

	return name[index:]
}
