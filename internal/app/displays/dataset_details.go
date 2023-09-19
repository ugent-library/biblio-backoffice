package displays

import (
	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/locale"
	"github.com/ugent-library/biblio-backoffice/internal/render/display"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
	"github.com/ugent-library/biblio-backoffice/models"
)

func DatasetDetails(user *models.User, l *locale.Locale, d *models.Dataset) *display.Display {
	var identifierType, identifier string
	for _, key := range vocabularies.Map["dataset_identifier_types"] {
		if val := d.Identifiers.Get(key); val != "" {
			identifierType = key
			identifier = val
			break
		}
	}

	return display.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label:    l.T("builder.title"),
				Value:    d.Title,
				Required: true,
			},
			&display.Text{
				Label:    l.T("builder.identifier_type"),
				Value:    l.TS("identifier", identifierType),
				Required: true,
			},
			&display.Link{
				Label:    l.T("builder.identifier"),
				Value:    identifier,
				URL:      identifiers.Resolve(identifierType, identifier),
				Required: true,
			},
		).
		AddSection(
			&display.List{
				Label:  l.T("builder.language"),
				Values: localize.LanguageNames(l, d.Language)},
			&display.Text{
				Label:    l.T("builder.year"),
				Value:    d.Year,
				Required: true,
			},
			&display.Text{
				Label:    l.T("builder.publisher"),
				Value:    d.Publisher,
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
				Value:    l.TS("dataset_licenses", d.License),
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.other_license"),
				Value: d.OtherLicense,
			},
			&display.Text{
				Label:    l.T("builder.access_level"),
				Value:    l.TS("dataset_access_levels", d.AccessLevel),
				Required: true,
			},
			&display.Text{
				Label: l.T("builder.embargo_date"),
				Value: d.EmbargoDate,
			},
			&display.Text{
				Label: l.T("builder.access_level_after_embargo"),
				Value: l.TS("dataset_access_levels_after_embargo", d.AccessLevelAfterEmbargo),
			},
		)
}
