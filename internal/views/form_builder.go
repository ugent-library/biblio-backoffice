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

type formOption func(*formData) error

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

func (b *FormBuilder) newFormData(opts []formOption) (*formData, error) {
	d := &formData{}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return d, err
		}
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

	return d, nil
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
	return func(d *formData) error {
		d.Name = str
		return nil
	}
}

func (b *FormBuilder) Value(v interface{}) formOption {
	return func(d *formData) error {
		d.Value = v
		return nil
	}
}

func (b *FormBuilder) Label(str string) formOption {
	return func(d *formData) error {
		d.Label = str
		return nil
	}
}

func (b *FormBuilder) Tooltip(str string) formOption {
	return func(d *formData) error {
		d.Tooltip = str
		return nil
	}
}

func (b *FormBuilder) Placeholder(str string) formOption {
	return func(d *formData) error {
		d.Placeholder = str
		return nil
	}
}

func (b *FormBuilder) Required() formOption {
	return func(d *formData) error {
		d.Required = true
		return nil
	}
}

func (b *FormBuilder) Checked() formOption {
	return func(d *formData) error {
		d.Checked = true
		return nil
	}
}

// TODO use functional options here too
func (b *FormBuilder) Choices(choices []string, scopes ...string) formOption {
	return func(d *formData) error {
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
		return nil
	}
}

func (b *FormBuilder) EmptyChoice() formOption {
	return func(d *formData) error {
		d.EmptyChoice = true
		return nil
	}
}

func (b *FormBuilder) Cols(num int) formOption {
	return func(d *formData) error {
		d.Cols = num
		return nil
	}
}

func (b *FormBuilder) Rows(num int) formOption {
	return func(d *formData) error {
		d.Rows = num
		return nil
	}
}

func (b *FormBuilder) ErrorPointer(ptr string) formOption {
	return func(d *formData) error {
		d.errorPointer = ptr
		return nil
	}
}

func (b *FormBuilder) Text(opts ...formOption) (template.HTML, error) {
	d, err := b.newFormData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "form_builder/_text", d)
}

func (b *FormBuilder) TextMultiple(opts ...formOption) (template.HTML, error) {
	d, err := b.newFormData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "form_builder/_text_multiple", d)
}

func (b *FormBuilder) Checkbox(opts ...formOption) (template.HTML, error) {
	d, err := b.newFormData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "form_builder/_checkbox", d)
}

func (b *FormBuilder) List(opts ...formOption) (template.HTML, error) {
	d, err := b.newFormData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "form_builder/_list", d)
}

func (b *FormBuilder) ListMultiple(opts ...formOption) (template.HTML, error) {
	d, err := b.newFormData(opts)
	if err != nil {
		return template.HTML(""), err
	}
	return RenderPartial(b.renderer, "form_builder/_list_multiple", d)
}
