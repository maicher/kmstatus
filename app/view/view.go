package view

import (
	"bytes"
	"path"
	"strings"
	"text/template"

	"github.io/maicher/kmstatus/app/config"
)

const DefaultTemplate = `{{.CPU.Load | round 1}}% {{.CPU.Freq | humanSI 1}}Hz 
{{range $t := .CPU.Temp}}{{$t}}°{{end}}
 {{.Mem.MemUsed | human 0}}({{.Mem.MemTotal | human 0}}) Swap: {{.Mem.SwapUsed | human 0}}({{.Mem.SwapTotal | human 0}})`

type View struct {
	templateName string
	templates    *template.Template
}

type RenderView struct{}

func (v *View) Render(d *Data) string {
	b := bytes.Buffer{}
	err := v.templates.ExecuteTemplate(&b, v.templateName, d)
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(b.String(), "\n", "")
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
		template.New(config.DefaultTemplateName).Funcs(helpers).Parse(DefaultTemplate),
	)
}
