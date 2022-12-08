package app

import (
	"fmt"
	"strconv"
	"text/template"

	"github.io/maicher/stbar/pkg/parsers/cpu"
)

var helpers = template.FuncMap{
	"round":   round,
	"human":   human,
	"humanSI": humanSI,
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
