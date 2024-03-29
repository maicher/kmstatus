package mem

import (
	"fmt"
	"strconv"
	"text/template"
)

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"human": human,
}

// Converts a number to a human-readable string.
// 1M == 1024k
func human(precision int, v int) string {
	return humanUnit(1024, precision, v, "MGT")
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
