package clock

import (
	"text/template"
	"time"
)

// Helper functions to be used in the templates.
var helpers = template.FuncMap{
	"format": format,
}

func format(layout string, t time.Time) string {
	return t.Format(layout)
}
