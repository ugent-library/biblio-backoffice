package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func PublicationDetails(l *locale.Locale, p *models.Publication) *display.Display {
	switch p.Type {
	case "book_chapter":
		return bookChapterDetails(l, p)
	case "book_editor":
		return bookEditorDetails(l, p)
	case "book":
		return bookDetails(l, p)
	case "conference":
		return conferenceDetails(l, p)
	case "dissertation":
		return dissertationDetails(l, p)
	case "issue_editor":
		return issueEditorDetails(l, p)
	case "journal_article":
		return journalArticleDetails(l, p)
	case "miscellaneous":
		return miscellaneousDetails(l, p)
	default:
		return display.New()
	}
}
