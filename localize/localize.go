package localize

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// TODO hardcoded to English for now
var languageNamer = display.Languages(language.MustParse("en"))

func ValidationErrors(loc *gotext.Locale, errs *okay.Errors) []string {
	msgs := make([]string, len(errs.Errors))
	for i, e := range errs.Errors {
		msgs[i] = loc.Get("validation." + e.Rule)
	}
	return msgs
}

func ValidationErrorAt(loc *gotext.Locale, errs *okay.Errors, key string) string {
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

func LanguageSelectOptions() []form.SelectOption {
	vals, ok := vocabularies.Map["language_codes"]
	if !ok {
		return nil
	}

	opts := make([]form.SelectOption, len(vals))

	for i, v := range vals {
		opts[i] = form.SelectOption{
			Value: v,
			Label: LanguageName(v),
		}
	}

	return opts
}

func VocabularyTerms(loc *gotext.Locale, key string) map[string]string {
	vals, ok := vocabularies.Map[key]
	if !ok {
		return nil
	}

	translatedTerms := make(map[string]string, len(vals))

	for _, v := range vals {
		translatedTerms[v] = loc.Get(key + "." + v)
	}

	return translatedTerms
}

func VocabularySelectOptions(loc *gotext.Locale, key string) []form.SelectOption {
	vals, ok := vocabularies.Map[key]
	if !ok {
		return nil
	}

	opts := make([]form.SelectOption, len(vals))

	for i, v := range vals {
		opts[i] = form.SelectOption{
			Value: v,
			Label: loc.Get(key + "." + v),
		}
	}

	return opts
}
