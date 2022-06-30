package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func FormPublicationDetails(ctx Context, b *BindDetails, errors validation.Errors) *form.Form {

	switch ctx.Publication.Type {
	case "book_chapter":
		return FormTypeBookChapter(ctx, b, errors)
	case "book_editor":
		return FormTypeBookEditor(ctx, b, errors)
	case "book":
		return FormTypeBook(ctx, b, errors)
	case "conference":
		return FormTypeConference(ctx, b, errors)
	case "dissertation":
		return FormTypeDissertation(ctx, b, errors)
	case "issue_editor":
		return FormTypeIssueEditor(ctx, b, errors)
	case "journal_article":
		return FormTypeJournalArticle(ctx, b, errors)
	case "miscellaneous":
		return FormTypeMiscellaneous(ctx, b, errors)
	default:
		return form.New()
	}
}
