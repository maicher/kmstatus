package network

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"human":     human,
	"ljust":     ljust,
	"hasPrefix": hasPrefix,
}

// Converts a number to a human-readable string.
// 1M == 1024k
func human(precision int, v int) string {
	return humanUnit(1024, precision, v, "kMGT")
}

func humanUnit(unit int, precision int, v int, x string) string {
	if v < unit {
		return strconv.Itoa(v)
	}

	div, exp := int(unit), 0
	for n := v / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.*f%c", precision, float64(v)/float64(div), x[exp])
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

func hasPrefix(prefix, s string) bool {
	return strings.HasPrefix(s, prefix)
}
