package views

import (
	"bytes"
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/presenters"
	r "github.com/ugent-library/biblio-backend/internal/renderer"
)

type View interface {
	NewView(presenter presenters.Presenter) *View
	Render() template.HTML
}

type DescriptionView struct {
	Presenter presenters.Presenter
}

func (gv *DescriptionView) Render() string {
	tree := gv.Presenter.Process()

	r.DoRender(tree)

	content := template.HTML(r.DoRender(tree))

	// Wrap the rendered tree in a description template
	tpl := `
	<div class="description">
		<h2>Title: {{.Title}}</h2>
		<p>Show some content below</p>
		{{ .Content }}
	</div>
	`

	t := template.Must(template.New("description").Parse(tpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, struct {
		Title   string
		Content template.HTML
	}{"Description", content})

	if err != nil {
		panic(err)
	}

	return buf.String()
}
