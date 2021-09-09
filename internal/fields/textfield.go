package fields

import (
	"bytes"
	"html/template"
)

type TextField struct {
	Label string
	Value string
}

func (tf *TextField) Render() string {
	tpl := `
	 	<p><strong>{{ .Label }}</strong>: {{ .Value }}</p>
	`

	t := template.Must(template.New("field").Parse(tpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, struct {
		Label string
		Value string
	}{Label: tf.Label, Value: tf.Value})

	if err != nil {
		panic(err)
	}

	return buf.String()
}
