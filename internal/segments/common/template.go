package common

import (
	"strings"
	"text/template"
)

type Template struct {
	Tmpl *template.Template
}

func (t *Template) NewTemplate(templateString string, helpers template.FuncMap) error {
	var err error

	templateString = strings.ReplaceAll(templateString, "\n", "")
	t.Tmpl, err = template.New("").Funcs(helpers).Parse(templateString)

	return err
}
