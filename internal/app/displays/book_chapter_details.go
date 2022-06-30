package displays

import (
	"github.com/ugent-library/biblio-backend/internal/app/helpers"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func bookChapterDetails(l *locale.Locale, p *models.Publication) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", p.Type),
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
				Label:    l.T("builder.book_chapter.publication"),
				Value:    p.Publication,
				Required: true,
			},
		).
		AddSection(
			&display.Text{
				Label:  l.T("builder.language"),
				List:   true,
				Values: localize.LanguageNames(l, p.Language)},
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
				Label: l.T("builder.pages"),
				Value: helpers.FormatRange(p.PageFirst, p.PageLast),
			},
			&display.Text{
				Label: l.T("builder.page_count"),
				Value: p.PageCount,
			},
		).
		AddSection(
			&display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   l.TS("tooltip.publication", p.WOSType),
				Tooltip: l.T("tooltip.publication.wos_type"),
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
		)
}
