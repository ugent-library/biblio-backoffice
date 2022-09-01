package publicationediting

import (
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func miscellaneousDetailsForm(user *models.User, l *locale.Locale, publication *models.Publication, errors validation.Errors) *form.Form {
	var wosTypeField form.Field
	if user.CanCuratePublications() {
		wosTypeField = &form.Text{
			Name:    "wos_type",
			Label:   l.T("builder.wos_type"),
			Value:   publication.WOSType,
			Tooltip: l.T("tooltip.publication.wos_type"),
		}
	} else {
		wosTypeField = &display.Text{
			Label:   l.T("builder.wos_type"),
			Value:   publication.WOSType,
			Tooltip: l.T("tooltip.publication.wos_type"),
		}
	}

	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors)).
		AddSection(
			&display.Text{
				Label: l.T("builder.type"),
				Value: l.TS("publication_types", publication.Type),
			},
			&form.Select{
				Name:    "miscellaneous_type",
				Label:   l.T("builder.miscellaneous_type"),
				Value:   publication.MiscellaneousType,
				Options: localize.VocabularySelectOptions(l, "miscellaneous_types"),
				Cols:    3,
				Error:   localize.ValidationErrorAt(l, errors, "/miscellaneous_type"),
			},
			&form.Text{
				Name:  "doi",
				Label: l.T("builder.doi"),
				Value: publication.DOI,
				Cols:  9,
				Help:  l.T("builder.doi.help"),
				Error: localize.ValidationErrorAt(l, errors, "/doi"),
			},
			&display.Text{
				Label:   l.T("builder.classification"),
				Value:   l.TS("publication_classifications", publication.Classification),
				Tooltip: l.T("tooltip.publication.classification"),
			},
		).
		AddSection(
			&form.Text{
				Name:     "title",
				Label:    l.T("builder.title"),
				Value:    publication.Title,
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/title"),
				Required: true,
			},
			&form.TextRepeat{
				Name:   "alternative_title",
				Label:  l.T("builder.alternative_title"),
				Values: publication.AlternativeTitle,
				Cols:   9,
				Error:  localize.ValidationErrorAt(l, errors, "/alternative_title"),
			},
			&form.Text{
				Name:     "publication",
				Label:    l.T("builder.miscellaneous.publication"),
				Value:    publication.Publication,
				Required: true,
				Cols:     9,
				Error:    localize.ValidationErrorAt(l, errors, "/publication"),
			},
			&form.Text{
				Name:  "publication_abbreviation",
				Label: l.T("builder.miscellaneous.publication_abbreviation"),
				Value: publication.PublicationAbbreviation,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/publication_abbreviation"),
			},
		).
		AddSection(
			&form.SelectRepeat{
				Name:        "language",
				Label:       l.T("builder.language"),
				Options:     localize.LanguageSelectOptions(l),
				Values:      publication.Language,
				EmptyOption: true,
				Cols:        9,
				Error:       localize.ValidationErrorAt(l, errors, "/language"),
			},
			&form.Select{
				Name:        "publication_status",
				Label:       l.T("builder.publication_status"),
				EmptyOption: true,
				Options:     localize.VocabularySelectOptions(l, "publication_publishing_statuses"),
				Value:       publication.PublicationStatus,
				Cols:        3,
				Error:       localize.ValidationErrorAt(l, errors, "/publication_status"),
			},
			&form.Checkbox{
				Name:    "extern",
				Label:   l.T("builder.extern"),
				Value:   "true",
				Checked: publication.Extern,
				Cols:    9,
				Error:   localize.ValidationErrorAt(l, errors, "/extern"),
			},
			&form.Text{
				Name:     "year",
				Label:    l.T("builder.year"),
				Value:    publication.Year,
				Required: true,
				Cols:     3,
				Help:     l.T("builder.year.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/year"),
			},
			&form.Text{
				Name:  "place_of_publication",
				Label: l.T("builder.place_of_publication"),
				Value: publication.PlaceOfPublication,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/place_of_publication"),
			},
			&form.Text{
				Name:  "publisher",
				Label: l.T("builder.publisher"),
				Value: publication.Publisher,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/publisher"),
			},
		).
		AddSection(
			&form.Text{
				Name:  "series_title",
				Label: l.T("builder.series_title"),
				Value: publication.SeriesTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/series_title"),
			},
			&form.Text{
				Name:  "volume",
				Label: l.T("builder.volume"),
				Value: publication.Volume,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/volume"),
			},
			&form.Text{
				Name:  "issue",
				Label: l.T("builder.issue"),
				Value: publication.Issue,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/issue"),
			},
			&form.Text{
				Name:  "edition",
				Label: l.T("builder.edition"),
				Value: publication.Edition,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/edition"),
			},
			&form.Text{
				Name:  "page_first",
				Label: l.T("builder.page_first"),
				Value: publication.PageFirst,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_first"),
			},
			&form.Text{
				Name:  "page_last",
				Label: l.T("builder.page_last"),
				Value: publication.PageLast,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_last"),
			},
			&form.Text{
				Name:  "page_count",
				Label: l.T("builder.page_count"),
				Value: publication.PageCount,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_count"),
			},
			&form.Text{
				Name:  "article_number",
				Label: l.T("builder.article_number"),
				Value: publication.ArticleNumber,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/article_number"),
			},
			&form.Text{
				Name:  "issue_title",
				Label: l.T("builder.issue_title"),
				Value: publication.IssueTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/issue_title"),
			},
			&form.Text{
				Name:  "report_number",
				Label: l.T("builder.report_number"),
				Value: publication.ReportNumber,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/report_number"),
			},
		).
		AddSection(
			wosTypeField,
			&form.Text{
				Name:  "wos_id",
				Label: l.T("builder.wos_id"),
				Value: publication.WOSID,
				Cols:  3,
				Help:  l.T("builder.wos_id.help"),
				Error: localize.ValidationErrorAt(l, errors, "/wos_id"),
			},
			&form.TextRepeat{
				Name:   "issn",
				Label:  l.T("builder.issn"),
				Values: publication.ISSN,
				Cols:   3,
				Help:   l.T("builder.issn.help"),
				Error:  localize.ValidationErrorAt(l, errors, "/issn"),
			},
			&form.TextRepeat{
				Name:   "eissn",
				Label:  l.T("builder.eissn"),
				Values: publication.EISSN,
				Cols:   3,
				Help:   l.T("builder.eissn.help"),
				Error:  localize.ValidationErrorAt(l, errors, "/eissn"),
			},
			&form.TextRepeat{
				Name:   "isbn",
				Label:  l.T("builder.isbn"),
				Values: publication.ISBN,
				Cols:   3,
				Help:   l.T("builder.isbn.help"),
				Error:  localize.ValidationErrorAt(l, errors, "/isbn"),
			},
			&form.TextRepeat{
				Name:   "eisbn",
				Label:  l.T("builder.eisbn"),
				Values: publication.EISBN,
				Cols:   3,
				Help:   l.T("builder.eisbn.help"),
				Error:  localize.ValidationErrorAt(l, errors, "/eisbn"),
			},
		)
}
