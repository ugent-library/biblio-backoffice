package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func FormPublicationDetails(l *locale.Locale, b BindDetails, errors validation.Errors) *form.Form {
	switch b.Type {
	case "book_chapter":
		return FormTypeBookChapter(l, b, errors)
	case "book_editor":
		return FormTypeBookEditor(l, b, errors)
	case "book":
		return FormTypeBook(l, b, errors)
	case "conference":
		return FormTypeConference(l, b, errors)
	case "dissertation":
		return FormTypeDissertation(l, b, errors)
	case "issue_editor":
		return FormTypeIssueEditor(l, b, errors)
	case "journal_article":
		return FormTypeJournalArticle(l, b, errors)
	case "miscellaneous":
		return FormTypeMiscellaneous(l, b, errors)
	default:
		return form.New()
	}
}
