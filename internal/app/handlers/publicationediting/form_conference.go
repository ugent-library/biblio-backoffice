package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func FormConference(l *locale.Locale, b BindDetails, errors validation.Errors) *form.Form {
	return form.New().
		WithTheme("default")
}
