package localize

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/views/form"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// TODO hardcoded to English for now
var languageNamer = display.Languages(language.MustParse("en"))

func ValidationErrors(loc *gotext.Locale, errs *okay.Errors) []string {
	if errs == nil || len(errs.Errors) == 0 {
		return nil
	}
	msgs := make([]string, len(errs.Errors))
	for i, e := range errs.Errors {
		msgs[i] = loc.Get("validation." + e.Rule)
	}
	return msgs
}

func ValidationErrorAt(loc *gotext.Locale, errs *okay.Errors, key string) string {
	if errs == nil {
		return ""
	}
	err := errs.Get(key)
	if err == nil {
		return ""
	}
	return loc.Get("validation." + err.Rule)
}

// TODO memoize this
func LanguageName(code string) string {
	tag := language.Make(code)
	if name := languageNamer.Name(tag); name != "" {
		return name
	}
	return code

}

func LanguageNames(codes []string) []string {
	names := make([]string, len(codes))
	for i, code := range codes {
		names[i] = LanguageName(code)
	}
	return names
}

func ClassificationSelectOptions(loc *gotext.Locale, vals []string) []form.Option {
	opts := make([]form.Option, len(vals))
	for i, v := range vals {
		opts[i] = form.Option{
			Value: v,
			Label: loc.Get("publication_classifications." + v),
		}
	}
	return opts
}

func ResearchFieldOptions(loc *gotext.Locale) []form.Option {
	opts := make([]form.Option, len(vocabularies.Map["research_fields"]))
	for i, v := range vocabularies.Map["research_fields"] {
		opts[i].Label = v
		opts[i].Value = v
	}
	return opts
}

func LanguageSelectOptions() []form.Option {
	vals, ok := vocabularies.Map["language_codes"]
	if !ok {
		return nil
	}

	opts := make([]form.Option, len(vals))

	for i, v := range vals {
		opts[i] = form.Option{
			Value: v,
			Label: LanguageName(v),
		}
	}

	return opts
}

func VocabularySelectOptions(loc *gotext.Locale, key string) []form.Option {
	vals, ok := vocabularies.Map[key]
	if !ok {
		return nil
	}

	opts := make([]form.Option, len(vals))

	for i, v := range vals {
		opts[i] = form.Option{
			Value: v,
			Label: loc.Get(key + "." + v),
		}
	}

	return opts
}
