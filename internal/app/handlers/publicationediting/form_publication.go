package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func formPublicationDetails(ctx Context, b *BindDetails, errors validation.Errors) *form.Form {

	switch ctx.Publication.Type {
	case "book_chapter":
		return formTypeBookChapter(ctx, b, errors)
	case "book_editor":
		return formTypeBookEditor(ctx, b, errors)
	case "book":
		return formTypeBook(ctx, b, errors)
	case "conference":
		return formTypeConference(ctx, b, errors)
	case "dissertation":
		return formTypeDissertation(ctx, b, errors)
	case "issue_editor":
		return formTypeIssueEditor(ctx, b, errors)
	case "journal_article":
		return formTypeJournalArticle(ctx, b, errors)
	case "miscellaneous":
		return formTypeMiscellaneous(ctx, b, errors)
	default:
		return form.New()
	}
}
