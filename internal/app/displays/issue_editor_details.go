package displays

import (
	"github.com/ugent-library/biblio-backoffice/internal/app/helpers"
	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/locale"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render/display"
)

func issueEditorDetails(user *models.User, l *locale.Locale, p *models.Publication) *display.Display {
	d := display.New().
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
			&display.List{
				Label:  l.T("builder.alternative_title"),
				Values: p.AlternativeTitle,
			},
			&display.Text{
				Label:    l.T("builder.journal_article.publication"),
				Value:    p.Publication,
				Required: p.ShowPublicationAsRequired(),
			},
			&display.Text{
				Label: l.T("builder.journal_article.publication_abbreviation"),
				Value: p.PublicationAbbreviation,
			},
		).
		AddSection(
			&display.List{
				Label:  l.T("builder.language"),
				Values: localize.LanguageNames(l, p.Language)},
			&display.Text{
				Label: l.T("builder.publication_status"),
				Value: l.TS("publication_publishing_statuses", p.PublicationStatus),
			},
			&display.Text{
				Label:         l.T("builder.extern"),
				Value:         helpers.FormatBool(p.Extern, "true", "false"),
				ValueTemplate: "format/boolean_string",
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
				Label: l.T("builder.page_count"),
				Value: p.PageCount,
			},
		).
		AddSection(
			&display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   p.WOSType,
				Tooltip: l.T("tooltip.publication.wos_type"),
			},
			&display.Text{
				Label:         l.T("builder.wos_id"),
				Value:         p.WOSID,
				ValueTemplate: "format/wos_id",
			},
			&display.List{
				Label:  l.T("builder.issn"),
				Values: p.ISSN,
			},
			&display.List{
				Label:  l.T("builder.eissn"),
				Values: p.EISSN,
			},
			&display.List{
				Label:  l.T("builder.isbn"),
				Values: p.ISBN,
			},
			&display.List{
				Label:  l.T("builder.eisbn"),
				Values: p.EISBN,
			},
		)

	if user.CanCurate() {
		d.Sections[0].Fields = append(d.Sections[0].Fields, &display.Text{
			Label:         l.T("builder.legacy"),
			Value:         helpers.FormatBool(p.Legacy, "true", "false"),
			ValueTemplate: "format/boolean_string",
		})
	}

	return d
}
