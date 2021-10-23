package views

import (
	"html/template"

	"github.com/unrolled/render"
)

type showData struct {
	Value    interface{}
	Label    string
	Required bool
}

type showOption func(*showData) error

type ShowBuilder struct {
	renderer *render.Render
}

func NewShowBuilder(r *render.Render) *ShowBuilder {
	return &ShowBuilder{
		renderer: r,
	}
}

func (b *ShowBuilder) newShowData(opts []showOption) (*showData, error) {
	d := &showData{}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return d, err
		}
	}

	return d, nil
}

func (b *ShowBuilder) Value(str string) showOption {
	return func(d *showData) error {
		d.Value = str
		return nil
	}
}

func (b *ShowBuilder) Label(str string) showOption {
	return func(d *showData) error {
		d.Label = str
		return nil
	}
}

func (b *ShowBuilder) Required() showOption {
	return func(d *showData) error {
		d.Required = true
		return nil
	}
}

func (b *ShowBuilder) Text(opts ...showOption) (template.HTML, error) {
	d, err := b.newShowData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "show_builder/_text", d)
}
