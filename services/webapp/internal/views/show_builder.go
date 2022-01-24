package views

import (
	"html/template"

	"github.com/ugent-library/go-locale/locale"
)

type showData struct {
	values        []string
	htmlValues    []template.HTML
	valueTemplate string
	Label         string
	Required      bool
	Tooltip       string
}

func (f *showData) Value() string {
	if len(f.values) > 0 {
		return f.values[0]
	}
	return ""
}

func (f *showData) HTMLValue() template.HTML {
	if len(f.htmlValues) > 0 {
		return f.htmlValues[0]
	}
	return ""
}

func (f *showData) Values() []string {
	return f.values
}

type showOption func(*showData)

type showLocaleOption func(string) string

type ShowBuilder struct {
	renderPartial func(string, interface{}) (template.HTML, error)
	locale        *locale.Locale
	localeScope   string
}

func NewShowBuilder(r func(string, interface{}) (template.HTML, error), l *locale.Locale) *ShowBuilder {
	return &ShowBuilder{
		renderPartial: r,
		locale:        l,
		localeScope:   "builder",
	}
}

func (b *ShowBuilder) newShowData(opts []showOption) *showData {
	d := &showData{}

	for _, opt := range opts {
		opt(d)
	}

	if d.valueTemplate != "" {
		d.htmlValues = make([]template.HTML, len(d.values))
		for i, v := range d.values {
			d.htmlValues[i], _ = b.renderPartial(d.valueTemplate, v)
		}
	}

	return d
}

func (b *ShowBuilder) Locale(scope string) showLocaleOption {
	return func(str string) string {
		return b.locale.Translate(scope, str)
	}
}

func (b *ShowBuilder) LanguageName() showLocaleOption {
	return func(str string) string {
		return b.locale.LanguageName(str)
	}
}

func (b *ShowBuilder) Label(v string, localeOpts ...showLocaleOption) showOption {
	return func(d *showData) {
		if len(localeOpts) > 0 {
			opt := localeOpts[0]
			if lbl := opt(v); lbl != "" {
				d.Label = lbl
			}
		}
		if d.Label == "" {
			d.Label = b.locale.Translate(b.localeScope, v)
		}
	}
}

func (b *ShowBuilder) Value(v string, localeOpts ...showLocaleOption) showOption {
	return func(d *showData) {
		d.values = []string{v}

		if v != "" && len(localeOpts) > 0 {
			opt := localeOpts[0]
			if lbl := opt(v); lbl != "" {
				d.values[0] = lbl
			}
		}
	}
}

func (b *ShowBuilder) Values(values []string, localeOpts ...showLocaleOption) showOption {
	return func(d *showData) {
		d.values = make([]string, len(values))
		copy(d.values, values)

		if len(localeOpts) > 0 {
			opt := localeOpts[0]
			for i, v := range values {
				if lbl := opt(v); lbl != "" {
					d.values[i] = lbl
				}
			}
		}
	}
}

func (b *ShowBuilder) Required() showOption {
	return func(d *showData) {
		d.Required = true
	}
}

func (b *ShowBuilder) Tooltip(v string) showOption {
	return func(d *showData) {
		d.Tooltip = v
	}
}

func (b *ShowBuilder) ValueTemplate(v string) showOption {
	return func(d *showData) {
		d.valueTemplate = v
	}
}

func (b *ShowBuilder) Text(opts ...showOption) (template.HTML, error) {
	return b.renderPartial("show_builder/_text", b.newShowData(opts))
}

func (b *ShowBuilder) List(opts ...showOption) (template.HTML, error) {
	return b.renderPartial("show_builder/_list", b.newShowData(opts))
}

func (b *ShowBuilder) BadgeList(opts ...showOption) (template.HTML, error) {
	return b.renderPartial("show_builder/_badge_list", b.newShowData(opts))
}
