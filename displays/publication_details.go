package displays

import (
	"github.com/ugent-library/biblio-backoffice/locale"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func PublicationDetails(user *models.User, l *locale.Locale, p *models.Publication) *display.Display {
	switch p.Type {
	case "book_chapter":
		return bookChapterDetails(user, l, p)
	case "book_editor":
		return bookEditorDetails(user, l, p)
	case "book":
		return bookDetails(user, l, p)
	case "conference":
		return conferenceDetails(user, l, p)
	case "dissertation":
		return dissertationDetails(user, l, p)
	case "issue_editor":
		return issueEditorDetails(user, l, p)
	case "journal_article":
		return journalArticleDetails(user, l, p)
	case "miscellaneous":
		return miscellaneousDetails(user, l, p)
	default:
		return display.New()
	}
}
