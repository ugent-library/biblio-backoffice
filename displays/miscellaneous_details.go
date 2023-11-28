package displays

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/helpers"
	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
)

func miscellaneousDetails(user *models.Person, loc *gotext.Locale, p *models.Publication) *display.Display {
	d := display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: loc.Get("builder.type"),
				Value: loc.Get("publication_types." + p.Type),
			},
			&display.Text{
				Label: loc.Get("builder.miscellaneous_type"),
				Value: loc.Get("miscellaneous_types." + p.MiscellaneousType),
			},
			&display.Link{
				Label: loc.Get("builder.doi"),
				Value: p.DOI,
				URL:   identifiers.DOI.Resolve(p.DOI),
			},
			&display.Text{
				Label: loc.Get("builder.classification"),
				Value: loc.Get("publication_classifications." + p.Classification),
			},
		).
		AddSection(
			&display.Text{
				Label:    loc.Get("builder.title"),
				Value:    p.Title,
				Required: true,
			},
			&display.List{
				Label:  loc.Get("builder.alternative_title"),
				Values: p.AlternativeTitle,
			},
			&display.Text{
				Label:    loc.Get("builder.publication"),
				Value:    p.Publication,
				Required: p.ShowPublicationAsRequired(),
			},
			&display.Text{
				Label: loc.Get("builder.publication_abbreviation"),
				Value: p.PublicationAbbreviation,
			},
		).
		AddSection(
			&display.List{
				Label:  loc.Get("builder.language"),
				Values: localize.LanguageNames(p.Language)},
			&display.Text{
				Label: loc.Get("builder.publication_status"),
				Value: loc.Get("publication_publishing_statuses." + p.PublicationStatus),
			},
			&display.Text{
				Label:         loc.Get("builder.extern"),
				Value:         helpers.FormatBool(p.Extern, "true", "false"),
				ValueTemplate: "format/boolean_string",
			},
			&display.Text{
				Label:    loc.Get("builder.year"),
				Value:    p.Year,
				Required: true,
			},
			&display.Text{
				Label: loc.Get("builder.place_of_publication"),
				Value: p.PlaceOfPublication,
			},
			&display.Text{
				Label: loc.Get("builder.publisher"),
				Value: p.Publisher,
			},
		).
		AddSection(
			&display.Text{
				Label: loc.Get("builder.series_title"),
				Value: p.SeriesTitle,
			},
			&display.Text{
				Label: loc.Get("builder.volume"),
				Value: p.Volume,
			},
			&display.Text{
				Label: loc.Get("builder.issue"),
				Value: p.Issue,
			},
			&display.Text{
				Label: loc.Get("builder.edition"),
				Value: p.Edition,
			},
			&display.Text{
				Label: loc.Get("builder.pages"),
				Value: helpers.FormatRange(p.PageFirst, p.PageLast),
			},
			&display.Text{
				Label: loc.Get("builder.page_count"),
				Value: p.PageCount,
			},
			&display.Text{
				Label: loc.Get("builder.article_number"),
				Value: p.ArticleNumber,
			},
			&display.Text{
				Label: loc.Get("builder.issue_title"),
				Value: p.IssueTitle,
			},
			&display.Text{
				Label: loc.Get("builder.report_number"),
				Value: p.ReportNumber,
			},
		).
		AddSection(
			&display.Text{
				Label:   loc.Get("builder.wos_type"),
				Value:   p.WOSType,
				Tooltip: loc.Get("tooltip.publication.wos_type"),
			},
			&display.Link{
				Label: loc.Get("builder.wos_id"),
				Value: p.WOSID,
				URL:   identifiers.WebOfScience.Resolve(p.WOSID),
			},
			&display.List{
				Label:  loc.Get("builder.issn"),
				Values: p.ISSN,
			},
			&display.List{
				Label:  loc.Get("builder.eissn"),
				Values: p.EISSN,
			},
			&display.List{
				Label:  loc.Get("builder.isbn"),
				Values: p.ISBN,
			},
			&display.List{
				Label:  loc.Get("builder.eisbn"),
				Values: p.EISBN,
			},
			&display.Text{
				Label: loc.Get("builder.pubmed_id"),
				Value: p.PubMedID,
			},
			&display.Text{
				Label: loc.Get("builder.arxiv_id"),
				Value: p.ArxivID,
			},
			&display.Text{
				Label: loc.Get("builder.esci_id"),
				Value: p.ESCIID,
			},
		)

	if user.CanCurate() {
		d.Sections[0].Fields = append(d.Sections[0].Fields, &display.Text{
			Label:         loc.Get("builder.legacy"),
			Value:         helpers.FormatBool(p.Legacy, "true", "false"),
			ValueTemplate: "format/boolean_string",
		})
	}

	return d
}
