package segments

import "text/template"

type Template struct {
	Tmpl *template.Template
}

func (t *Template) NewTemplate(templateString string, helpers template.FuncMap) error {
	var err error

	t.Tmpl, err = template.New("").Funcs(helpers).Parse(templateString)

	return err
}
