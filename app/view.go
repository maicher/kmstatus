package app

import (
	"bytes"
	"strings"
	"text/template"
)

type view struct {
	tmplName string
	tmpl     *template.Template
}

func (v *view) render(d *data) string {
	b := bytes.Buffer{}
	err := v.tmpl.ExecuteTemplate(&b, v.tmplName, d)
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(b.String(), "\n", "")
}

func newView(t string) *view {
	tmpl, err := template.New(t).Funcs(helpers).ParseFiles("app/templates/" + t)
	if err != nil {
		panic(err)
	}

	return &view{
		tmplName: t,
		tmpl:     tmpl,
	}
}
