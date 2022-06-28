package localize

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render/form"
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

func LanguageNames(loc *locale.Locale, codes []string) []string {
	names := make([]string, len(codes))
	for i, code := range codes {
		names[i] = loc.LanguageName(code)
	}
	return names
}

func LanguageSelectOptions(loc *locale.Locale) []form.SelectOption {
	codes := vocabularies.Map["language_codes"]
	opts := make([]form.SelectOption, len(codes))
	for i, code := range codes {
		opts[i] = form.SelectOption{
			Value: code,
			Label: loc.LanguageName(code),
		}
	}
	return opts
}

func LicenseSelectOptions(loc *locale.Locale) []form.SelectOption {
	codes := vocabularies.Map["cc_licenses"]
	opts := make([]form.SelectOption, len(codes))
	for i, code := range codes {
		opts[i] = form.SelectOption{
			Value: code,
			Label: loc.TranslateScope("cc_licenses", code),
		}
	}
	return opts
}

func AccessLevelSelectOptions(loc *locale.Locale) []form.SelectOption {
	codes := vocabularies.Map["access_levels"]
	opts := make([]form.SelectOption, len(codes))
	for i, code := range codes {
		opts[i] = form.SelectOption{
			Value: code,
			Label: loc.TranslateScope("access_levels", code),
		}
	}
	return opts
}
