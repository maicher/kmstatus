package cpu

import "text/template"

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"bar": bar,
}

func bar(val float64) string {
	if val < 2 {
		return " "
	}

	if val < 9 {
		return "░"
	}

	if val < 30 {
		return "▒"
	}

	if val < 90 {
		return "▓"
	}

	return "█"
}
