package fields

import (
	"bytes"
	"html/template"
)

type FieldSet struct {
	Label string
	Items []Field
}

func (fs *FieldSet) Render(children []string) string {
	tpl := `
		<div class="fieldset">
		    <h2>{{ .Label }}</h2>
			<ul class="items">
			{{ range .Items }}
				<li>{{ . }}</li>
			{{ end }}
			</div>
		</div>
	`

	// Haxx conversion to template.HTML
	var items []template.HTML
	for _, v := range children {
		items = append(items, template.HTML(v))
	}

	t := template.Must(template.New("fieldset").Parse(tpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, struct {
		Label string
		Items []template.HTML
	}{Label: fs.Label, Items: items})

	if err != nil {
		panic(err)
	}

	return buf.String()
}
