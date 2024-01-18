package displays

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

func DatasetDetails(user *models.Person, loc *gotext.Locale, d *models.Dataset) *display.Display {
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
				Label:    loc.Get("builder.title"),
				Value:    d.Title,
				Required: true,
			},
			&display.Text{
				Label:    loc.Get("builder.identifier_type"),
				Value:    loc.Get("identifier." + identifierType),
				Required: true,
			},
			&display.Link{
				Label:    loc.Get("builder.identifier"),
				Value:    identifier,
				URL:      identifiers.Resolve(identifierType, identifier),
				Required: true,
			},
		).
		AddSection(
			&display.List{
				Label:  loc.Get("builder.language"),
				Values: localize.LanguageNames(d.Language)},
			&display.Text{
				Label:    loc.Get("builder.year"),
				Value:    d.Year,
				Required: true,
			},
			&display.Text{
				Label:    loc.Get("builder.publisher"),
				Value:    d.Publisher,
				Required: true,
			},
		).
		AddSection(
			&display.List{
				Label:    loc.Get("builder.format"),
				Values:   d.Format,
				Required: true,
			},
			&display.List{
				Inline:        true,
				Label:         loc.Get("builder.keyword"),
				Values:        d.Keyword,
				ValueTemplate: "format/badge",
			},
		).
		AddSection(
			&display.Text{
				Label:    loc.Get("builder.license"),
				Value:    loc.Get("dataset_licenses." + d.License),
				Required: true,
			},
			&display.Text{
				Label: loc.Get("builder.other_license"),
				Value: d.OtherLicense,
			},
			&display.Text{
				Label:    loc.Get("builder.access_level"),
				Value:    loc.Get("dataset_access_levels." + d.AccessLevel),
				Required: true,
			},
			&display.Text{
				Label: loc.Get("builder.embargo_date"),
				Value: d.EmbargoDate,
			},
			&display.Text{
				Label: loc.Get("builder.access_level_after_embargo"),
				Value: loc.Get("dataset_access_levels_after_embargo." + d.AccessLevelAfterEmbargo),
			},
		)
}
