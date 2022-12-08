package app

import (
	"bytes"
	"strings"
	"text/template"
)

type View struct {
	tmplName string
	tmpl     *template.Template
}

func (v *View) Render(d *Data) string {
	b := bytes.Buffer{}
	err := v.tmpl.ExecuteTemplate(&b, v.tmplName, d)
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(b.String(), "\n", "")
}

func NewView(t string) *View {
	tmpl, err := template.New(t).Funcs(helpers).ParseFiles("app/templates/" + t)
	if err != nil {
		panic(err)
	}

	return &View{
		tmplName: t,
		tmpl:     tmpl,
	}
}
