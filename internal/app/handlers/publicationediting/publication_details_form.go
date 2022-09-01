package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func detailsForm(user *models.User, l *locale.Locale, publication *models.Publication, errors validation.Errors) *form.Form {
	switch publication.Type {
	case "book_chapter":
		return bookChapterDetailsForm(user, l, publication, errors)
	case "book_editor":
		return bookEditorDetailsForm(user, l, publication, errors)
	case "book":
		return bookDetailsForm(user, l, publication, errors)
	case "conference":
		return conferenceDetailsForm(user, l, publication, errors)
	case "dissertation":
		return dissertationDetailsForm(user, l, publication, errors)
	case "issue_editor":
		return issueEditorDetailsForm(user, l, publication, errors)
	case "journal_article":
		return journalArticleDetailsForm(user, l, publication, errors)
	case "miscellaneous":
		return miscellaneousDetailsForm(user, l, publication, errors)
	default:
		return form.New()
	}
}
