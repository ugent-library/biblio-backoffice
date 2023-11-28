package displays

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func PublicationDetails(user *models.Person, loc *gotext.Locale, p *models.Publication) *display.Display {
	switch p.Type {
	case "book_chapter":
		return bookChapterDetails(user, loc, p)
	case "book_editor":
		return bookEditorDetails(user, loc, p)
	case "book":
		return bookDetails(user, loc, p)
	case "conference":
		return conferenceDetails(user, loc, p)
	case "dissertation":
		return dissertationDetails(user, loc, p)
	case "issue_editor":
		return issueEditorDetails(user, loc, p)
	case "journal_article":
		return journalArticleDetails(user, loc, p)
	case "miscellaneous":
		return miscellaneousDetails(user, loc, p)
	default:
		return display.New()
	}
}
