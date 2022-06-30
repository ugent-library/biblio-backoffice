package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func FormTypeDissertation(l *locale.Locale, b *BindDetails, errors validation.Errors) *form.Form {

	confirmationOptions := optionsForVocabulary(l, "confirmations")

	return form.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", b.Type),
			},
			&form.Text{
				Name:  "doi",
				Label: l.T("builder.doi"),
				Value: b.DOI,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/doi"),
			},
			&display.Text{
				Label:   l.T("builder.classification"),
				Value:   l.TS("publication_classifications", b.Classification),
				Tooltip: l.T("tooltip.publication.classification"),
			},
		).
		AddSection(
			&form.Text{
				Name:     "title",
				Label:    l.T("builder.title"),
				Value:    b.Title,
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/title"),
				Required: true,
			},
			&form.TextRepeat{
				Name:   "alternative_title",
				Label:  l.T("builder.alternative_title"),
				Values: b.AlternativeTitle,
				Cols:   9,
				Error:  localize.ValidationErrorAt(l, errors, "/alternative_title"),
			},
			&form.Text{
				Name:  "publication_abbreviation",
				Label: l.T("builder.publication_abbreviation"),
				Value: b.PublicationAbbreviation,
				Error: localize.ValidationErrorAt(l, errors, "/publication_abbreviation"),
				Cols:  3,
			},
		).
		AddSection(
			&form.SelectRepeat{
				Name:        "language",
				Label:       l.T("builder.language"),
				Options:     localize.LanguageSelectOptions(l),
				Values:      b.Language,
				EmptyOption: true,
				Cols:        9,
				Error:       localize.ValidationErrorAt(l, errors, "/language"),
			},
			&form.Select{
				Name:        "publication_status",
				Label:       l.T("builder.publication_status"),
				EmptyOption: true,
				Options:     optionsForVocabulary(l, "publication_publishing_statuses"),
				Value:       b.PublicationStatus,
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/publication_status"),
			},
			&form.Checkbox{
				Name:    "extern",
				Label:   l.T("builder.extern"),
				Value:   "true",
				Checked: b.Extern,
				Cols:    9,
				Error:   localize.ValidationErrorAt(l, errors, "/extern"),
			},
			&form.Text{
				Name:     "year",
				Label:    l.T("builder.year"),
				Value:    b.Year,
				Required: true,
				Cols:     3,
				Error:    localize.ValidationErrorAt(l, errors, "/year"),
			},
			&form.Text{
				Name:  "place_of_publication",
				Label: l.T("builder.place_of_publication"),
				Value: b.PlaceOfPublication,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/place_of_publication"),
			},
			&form.Text{
				Name:  "publisher",
				Label: l.T("builder.publisher"),
				Value: b.Publisher,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/publisher"),
			},
		).
		AddSection(
			&form.Text{
				Name:  "volume",
				Label: l.T("builder.volume"),
				Value: b.Volume,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/volume"),
			},
			&form.Text{
				Name:  "page_count",
				Label: l.T("builder.page_count"),
				Value: b.PageCount,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_count"),
			},
			&form.Text{
				Name:  "series_title",
				Label: l.T("builder.series_title"),
				Value: b.SeriesTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/series_title"),
			},
		).
		AddSection(
			&form.Text{
				Name:        "defense_date",
				Label:       l.T("builder.defense_date"),
				Value:       b.DefenseDate,
				Required:    true,
				Cols:        3,
				Placeholder: "e.g. 2022-04-30",
				Error:       localize.ValidationErrorAt(l, errors, "/defense_date"),
			},
			&form.Text{
				Name:        "defense_time",
				Label:       l.T("builder.defense_time"),
				Value:       b.DefenseTime,
				Required:    true,
				Cols:        3,
				Placeholder: "e.g. 11:00",
				Error:       localize.ValidationErrorAt(l, errors, "/defense_time"),
			},
			&form.Text{
				Name:     "defense_place",
				Label:    l.T("builder.defense_place"),
				Value:    b.DefensePlace,
				Required: true,
				Cols:     3,
				Error:    localize.ValidationErrorAt(l, errors, "/defense_place"),
			},
		).
		AddSection(
			&form.RadioButtonGroup{
				Name:    "has_confidential_data",
				Label:   l.T("builder.has_confidential_data"),
				Value:   b.HasConfidentialData,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_confidential_data"),
			},
			&form.RadioButtonGroup{
				Name:    "has_patent_application",
				Label:   l.T("builder.has_patent_application"),
				Value:   b.HasPatentApplication,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_patent_application"),
			},
			&form.RadioButtonGroup{
				Name:    "has_publications_planned",
				Label:   l.T("builder.has_publications_planned"),
				Value:   b.HasPublicationsPlanned,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_publications_planned"),
			},
			&form.RadioButtonGroup{
				Name:    "has_published_material",
				Label:   l.T("builder.has_published_material"),
				Value:   b.HasPublishedMaterial,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_published_material"),
			},
		).
		AddSection(
			&display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   l.TS("tooltip.publication", b.WOSType),
				Tooltip: l.T("tooltip.publication.wos_type"),
			},
			&form.Text{
				Name:        "wos_id",
				Label:       l.T("builder.wos_id"),
				Value:       b.WOSID,
				Cols:        3,
				Placeholder: "e.g. 000503382400004",
				Error:       localize.ValidationErrorAt(l, errors, "/wos_id"),
			},
			&form.TextRepeat{
				Name:        "issn",
				Label:       l.T("builder.issn"),
				Values:      b.ISSN,
				Cols:        3,
				Placeholder: "e.g. 2049-3630",
				Error:       localize.ValidationErrorAt(l, errors, "/issn"),
			},
			&form.TextRepeat{
				Name:        "eissn",
				Label:       l.T("builder.eissn"),
				Values:      b.EISSN,
				Cols:        3,
				Placeholder: "e.g. 2049-3630",
				Error:       localize.ValidationErrorAt(l, errors, "/eissn"),
			},
			&form.TextRepeat{
				Name:        "isbn",
				Label:       l.T("builder.isbn"),
				Values:      b.ISBN,
				Cols:        3,
				Placeholder: "e.g. 2049-3630",
				Error:       localize.ValidationErrorAt(l, errors, "/isbn"),
			},
			&form.TextRepeat{
				Name:        "eisbn",
				Label:       l.T("builder.eisbn"),
				Values:      b.EISBN,
				Cols:        3,
				Placeholder: "e.g. 2049-3630",
				Error:       localize.ValidationErrorAt(l, errors, "/eisbn"),
			},
		)
}
