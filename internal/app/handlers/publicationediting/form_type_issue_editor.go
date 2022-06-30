package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func formTypeIssueEditor(ctx Context, b *BindDetails, errors validation.Errors) *form.Form {
	l := ctx.Locale
	p := ctx.Publication
	return form.New().
		WithTheme("default").
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", p.Type),
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
				Value:   l.TS("publication_classifications", p.Classification),
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
				Name:     "publication",
				Label:    l.T("builder.journal_article.publication"),
				Value:    b.Publication,
				Required: true,
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/publication"),
			},
			&form.Text{
				Name:  "publication_abbreviation",
				Label: l.T("builder.journal_article.publication_abbreviation"),
				Value: b.PublicationAbbreviation,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/publication_abbreviation"),
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
				Name:  "issue",
				Label: l.T("builder.issue"),
				Value: b.Issue,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/issue"),
			},
			&form.Text{
				Name:  "page_count",
				Label: l.T("builder.page_count"),
				Value: b.PageCount,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_count"),
			},
		).
		AddSection(
			&display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   l.TS("tooltip.publication", p.WOSType),
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
