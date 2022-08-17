package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func detailsForm(l *locale.Locale, publication *models.Publication, errors validation.Errors) *form.Form {
	switch publication.Type {
	case "book_chapter":
		return bookChapterDetailsForm(l, publication, errors)
	case "book_editor":
		return bookEditorDetailsForm(l, publication, errors)
	case "book":
		return bookDetailsForm(l, publication, errors)
	case "conference":
		return conferenceDetailsForm(l, publication, errors)
	case "dissertation":
		return dissertationDetailsForm(l, publication, errors)
	case "issue_editor":
		return issueEditorDetailsForm(l, publication, errors)
	case "journal_article":
		return journalArticleDetailsForm(l, publication, errors)
	case "miscellaneous":
		return miscellaneousDetailsForm(l, publication, errors)
	default:
		return form.New()
	}
}
