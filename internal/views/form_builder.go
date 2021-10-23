package views

import (
	"html/template"

	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type formData struct {
	Name          string
	Value         interface{}
	Label         string
	Tooltip       string
	Placeholder   string
	Required      bool
	Checked       bool
	Choices       []string
	ChoicesLabels []string
	EmptyChoice   bool
	Cols          int
	Rows          int
	Error         string
	errorPointer  string
}

type formOption func(*formData)

type FormBuilder struct {
	renderer *render.Render
	Locale   *locale.Locale
	Errors   jsonapi.Errors
}

func NewFormBuilder(r *render.Render, l *locale.Locale, e jsonapi.Errors) *FormBuilder {
	return &FormBuilder{
		renderer: r,
		Locale:   l,
		Errors:   e,
	}
}

func (b *FormBuilder) newFormData(opts []formOption) *formData {
	d := &formData{}

	for _, opt := range opts {
		opt(d)
	}

	if d.Label == "" {
		d.Label = b.Locale.Translate("form_builder", d.Name)
	}

	if d.errorPointer == "" {
		d.errorPointer = "/data/" + d.Name
	}

	if formErr := b.errorFor(d.errorPointer); formErr != nil {
		d.Error = formErr.Title
	}

	return d
}

func (b *FormBuilder) errorFor(pointer string) *jsonapi.Error {
	for _, err := range b.Errors {
		if err.Source.Pointer == pointer {
			return &err
		}
	}
	return nil
}

func (b *FormBuilder) Name(str string) formOption {
	return func(d *formData) {
		d.Name = str
	}
}

func (b *FormBuilder) Value(v interface{}) formOption {
	return func(d *formData) {
		d.Value = v
	}
}

func (b *FormBuilder) Label(str string) formOption {
	return func(d *formData) {
		d.Label = str
	}
}

func (b *FormBuilder) Tooltip(str string) formOption {
	return func(d *formData) {
		d.Tooltip = str
	}
}

func (b *FormBuilder) Placeholder(str string) formOption {
	return func(d *formData) {
		d.Placeholder = str
	}
}

func (b *FormBuilder) Required() formOption {
	return func(d *formData) {
		d.Required = true
	}
}

func (b *FormBuilder) Checked() formOption {
	return func(d *formData) {
		d.Checked = true
	}
}

// TODO use functional options here too
func (b *FormBuilder) Choices(choices []string, scopes ...string) formOption {
	return func(d *formData) {
		d.Choices = choices
		if len(scopes) > 0 {
			d.ChoicesLabels = make([]string, len(choices))
			scope := scopes[0]
			// pseudo locale scopes
			if scope == ":language_name" {
				for i, c := range choices {
					d.ChoicesLabels[i] = b.Locale.LanguageName(c)
				}
			} else {
				for i, c := range choices {
					d.ChoicesLabels[i] = b.Locale.Translate(scope, c)
				}
			}
		}
	}
}

func (b *FormBuilder) EmptyChoice() formOption {
	return func(d *formData) {
		d.EmptyChoice = true
	}
}

func (b *FormBuilder) Cols(num int) formOption {
	return func(d *formData) {
		d.Cols = num
	}
}

func (b *FormBuilder) Rows(num int) formOption {
	return func(d *formData) {
		d.Rows = num
	}
}

func (b *FormBuilder) ErrorPointer(ptr string) formOption {
	return func(d *formData) {
		d.errorPointer = ptr
	}
}

func (b *FormBuilder) Text(opts ...formOption) (template.HTML, error) {
	return RenderPartial(b.renderer, "form_builder/_text", b.newFormData(opts))
}

func (b *FormBuilder) TextMultiple(opts ...formOption) (template.HTML, error) {
	return RenderPartial(b.renderer, "form_builder/_text_multiple", b.newFormData(opts))
}

func (b *FormBuilder) Checkbox(opts ...formOption) (template.HTML, error) {
	return RenderPartial(b.renderer, "form_builder/_checkbox", b.newFormData(opts))
}

func (b *FormBuilder) List(opts ...formOption) (template.HTML, error) {
	return RenderPartial(b.renderer, "form_builder/_list", b.newFormData(opts))
}

func (b *FormBuilder) ListMultiple(opts ...formOption) (template.HTML, error) {
	return RenderPartial(b.renderer, "form_builder/_list_multiple", b.newFormData(opts))
}
