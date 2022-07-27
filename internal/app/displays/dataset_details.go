package displays

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
)

func DatasetDetails(l *locale.Locale, d *models.Dataset) *display.Display {
	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label:    l.T("builder.title"),
				Value:    d.Title,
				Required: true,
			},
			&display.Text{
				Label:         l.T("builder.doi"),
				Value:         d.DOI,
				Required:      true,
				ValueTemplate: "format/doi",
			},
			&display.Text{
				Label:         l.T("builder.url"),
				Value:         d.URL,
				ValueTemplate: "format/link",
			},
		).
		AddSection(
			&display.Text{
				Label:    l.T("builder.publisher"),
				Value:    d.Publisher,
				Required: true,
			},
			&display.Text{
				Label:    l.T("builder.year"),
				Value:    d.Year,
				Required: true,
			},
		).
		AddSection(
			&display.List{
				Label:    l.T("builder.format"),
				Values:   d.Format,
				Required: true,
			},
			&display.List{
				Inline:        true,
				Label:         l.T("builder.keyword"),
				Values:        d.Keyword,
				ValueTemplate: "format/badge",
			},
		).
		AddSection(
			&display.Text{
				Label:    l.T("builder.license"),
				Value:    l.TS("licenses", d.License),
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.other_license"),
				Value: d.OtherLicense,
			},
			&display.Text{
				Label:    l.T("builder.access_level"),
				Value:    l.TS("access_levels", d.AccessLevel),
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.embargo"),
				Value: d.Embargo,
			},
			&display.Text{
				Label: l.T("builder.embargo_to"),
				Value: l.TS("access_levels", d.EmbargoTo),
			},
		)
}
