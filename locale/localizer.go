package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

type Localizer struct {
	locales []*Locale
	matcher language.Matcher
}

func NewLocalizer(langs ...string) *Localizer {
	tags := make([]language.Tag, len(langs))
	locs := make([]*Locale, len(langs))

	for i, lang := range langs {
		tag := language.MustParse(lang)
		tags[i] = tag
		locs[i] = &Locale{
			Tag:           tag,
			languageNamer: display.Languages(tag),
			printer:       message.NewPrinter(tag),
		}
	}

	return &Localizer{
		locales: locs,
		matcher: language.NewMatcher(tags),
	}

}

func (l *Localizer) Locales() []*Locale {
	return l.locales
}

func (l *Localizer) DefaultLocale() *Locale {
	return l.locales[0]
}

func (l *Localizer) GetLocale(langs ...string) *Locale {
	_, match := language.MatchStrings(l.matcher, langs...)
	return l.locales[match]
}
