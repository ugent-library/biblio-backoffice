package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func detailsForm(ctx Context, b *BindDetails, errors validation.Errors) *form.Form {
	switch ctx.Publication.Type {
	case "book_chapter":
		return bookChapterDetailsForm(ctx, b, errors)
	case "book_editor":
		return bookEditorDetailsForm(ctx, b, errors)
	case "book":
		return bookDetailsForm(ctx, b, errors)
	case "conference":
		return conferenceDetailsForm(ctx, b, errors)
	case "dissertation":
		return dissertationDetailsForm(ctx, b, errors)
	case "issue_editor":
		return issueEditorDetailsForm(ctx, b, errors)
	case "journal_article":
		return journalArticleDetailsForm(ctx, b, errors)
	case "miscellaneous":
		return miscellaneousDetailsForm(ctx, b, errors)
	default:
		return form.New()
	}
}
