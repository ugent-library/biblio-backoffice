package localize

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

func ValidationErrors(loc *locale.Locale, errs validation.Errors) []string {
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = loc.TranslateScope("validation", e.Code)
	}
	return msgs
}

func ValidationErrorAt(loc *locale.Locale, errs validation.Errors, ptr string) string {
	err := errs.At(ptr)
	if err == nil {
		return ""
	}
	return loc.TranslateScope("validation", err.Code)
}

func LanguageSelectOptions(loc *locale.Locale) []render.SelectOption {
	codes := vocabularies.Map["language_codes"]
	opts := make([]render.SelectOption, len(codes))
	for i, code := range codes {
		opts[i] = render.SelectOption{
			Value: code,
			Label: loc.LanguageName(code),
		}
	}
	return opts
}
