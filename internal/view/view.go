package view

import (
	"bytes"
	"log"
	"path"
	"strings"
	"text/template"

	"github.io/maicher/kmstatus/internal/config"
	"github.io/maicher/kmstatus/internal/view/out"
)

type StatusSetter interface {
	SetStatus(string)
}

type View struct {
	templateName string
	templates    *template.Template
	display      StatusSetter
}

type RenderView struct{}

func (v *View) Render(d *Data) {
	b := bytes.Buffer{}
	err := v.templates.ExecuteTemplate(&b, v.templateName, d)
	if err != nil {
		panic(err)
	}

	text := b.String()
	text = strings.ReplaceAll(text, ")\n\n", " ")
	text = strings.ReplaceAll(text, "\n", "")

	v.display.SetStatus(text)
}

func (v *View) RenderErr(e error) {
	log.Println(e)
}

func New(conf *config.Config) (*View, error) {
	t, err := compileTemplates(conf)
	if err != nil {
		return nil, err
	}

	v := &View{
		templateName: path.Base(conf.TemplateName),
		templates:    t,
	}

	if conf.XWindow {
		v.display, err = out.NewWindow()
	} else {
		v.display, err = out.NewStd()
	}
	if err != nil {
		return nil, err
	}

	return v, nil
}

func compileTemplates(c *config.Config) (*template.Template, error) {
	path, err := c.FindTemplatePath()
	if err == nil {
		return template.New("").Funcs(helpers).ParseFiles(path)
	}

	if c.HasDefaultTemplate() {
		return mustCompileDefaultTemplate(), nil
	}

	return nil, err
}

func mustCompileDefaultTemplate() *template.Template {
	return template.Must(
		template.New(config.DefaultTemplateName).Funcs(helpers).Parse(string(config.DefaultTemplate)),
	)
}
