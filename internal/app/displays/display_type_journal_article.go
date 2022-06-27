package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/helpers"
)

func DisplayTypeJournalArticle(l *locale.Locale, p *models.Publication) *display.Display {
	trLangs := []string{}
	for _, lang := range p.Language {
		trLangs = append(trLangs, l.LanguageName(lang))
	}
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", p.Type),
			},
			&display.Text{
				Label: l.T("builder.journal_article_type"),
				Value: l.TS("journal_article_types", p.JournalArticleType),
			},
			&display.Text{
				Label:         l.T("builder.doi"),
				Value:         p.DOI,
				ValueTemplate: "format/doi",
			},
			&display.Text{
				Label: l.T("builder.classification"),
				Value: l.TS("publication_classifications", p.Classification),
			},
		).
		AddSection(
			&display.Text{
				Label:    l.T("builder.title"),
				Value:    p.Title,
				Required: true,
			},
			&display.Text{
				Label:  l.T("builder.alternative_title"),
				List:   true,
				Values: p.AlternativeTitle,
			},
			&display.Text{
				Label:    l.T("builder.journal_article.publication"),
				Value:    p.Publication,
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.journal_article.publication_abbreviation"),
				Value: p.PublicationAbbreviation,
			},
		).
		AddSection(
			&display.Text{
				Label:  l.T("builder.language"),
				List:   true,
				Values: trLangs,
			},
			&display.Text{
				Label: l.T("builder.publication_status"),
				Value: l.TS("publication_publishing_statuses", p.PublicationStatus),
			},
			&display.Text{
				Label: l.T("builder.extern"),
				Value: helpers.FormatBool(p.Extern, "✓", "-"),
			},
			&display.Text{
				Label:    l.T("builder.year"),
				Value:    p.Year,
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.place_of_publication"),
				Value: p.PlaceOfPublication,
			},
			&display.Text{
				Label: l.T("builder.publisher"),
				Value: p.Publisher,
			},
		).
		AddSection(
			&display.Text{
				Label: l.T("builder.volume"),
				Value: p.Volume,
			},
			&display.Text{
				Label: l.T("builder.issue"),
				Value: p.Issue,
			},
			&display.Text{
				Label: l.T("builder.pages"),
				Value: helpers.FormatRange(p.PageFirst, p.PageLast),
			},
			&display.Text{
				Label: l.T("builder.page_count"),
				Value: p.PageCount,
			},
			&display.Text{
				Label: l.T("builder.article_number"),
				Value: p.ArticleNumber,
			},
			&display.Text{
				Label: l.T("builder.issue_title"),
				Value: p.IssueTitle,
			},
		).
		AddSection(
			&display.Text{
				Label: l.T("builder.wos_type"),
				Value: p.WOSType,
			},
			&display.Text{
				Label: l.T("builder.wos_id"),
				Value: p.WOSID,
			},
			&display.Text{
				Label:  l.T("builder.issn"),
				List:   true,
				Values: p.ISSN,
			},
			&display.Text{
				Label:  l.T("builder.eissn"),
				List:   true,
				Values: p.EISSN,
			},
			&display.Text{
				Label:  l.T("builder.isbn"),
				List:   true,
				Values: p.ISBN,
			},
			&display.Text{
				Label:  l.T("builder.eisbn"),
				List:   true,
				Values: p.EISBN,
			},
			&display.Text{
				Label: l.T("builder.pubmed_id"),
				Value: p.PubMedID,
			},
			&display.Text{
				Label: l.T("builder.arxiv_id"),
				Value: p.ArxivID,
			},
			&display.Text{
				Label: l.T("builder.esci_id"),
				Value: p.ESCIID,
			},
		)
}
