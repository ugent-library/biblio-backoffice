package publicationediting

import (
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func detailsForm(user *models.User, l *locale.Locale, publication *models.Publication, errors validation.Errors) *form.Form {
	f := form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors))

	section1 := []form.Field{}

	section1 = append(section1, &form.Select{
		Template: "publication/type",
		Name:     "type",
		Label:    l.T("builder.type"),
		Options:  localize.VocabularySelectOptions(l, "publication_types"),
		Value:    publication.Type,
		Cols:     3,
		Help:     l.T("builder.type.help"),
		Vars: struct {
			Publication *models.Publication
		}{
			Publication: publication,
		},
	})

	if publication.UsesJournalArticleType() {
		section1 = append(section1, &form.Select{
			Name:        "journal_article_type",
			Label:       l.T("builder.journal_article_type"),
			Options:     localize.VocabularySelectOptions(l, "journal_article_types"),
			EmptyOption: true,
			Value:       publication.JournalArticleType,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/journal_article_type"),
		})
	}

	if publication.UsesConferenceType() {
		section1 = append(section1, &form.Select{
			Name:        "conference_type",
			Label:       l.T("builder.conference_type"),
			Value:       publication.ConferenceType,
			Options:     localize.VocabularySelectOptions(l, "conference_types"),
			EmptyOption: true,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/conference_type"),
		})
	}

	if publication.UsesMiscellaneousType() {
		section1 = append(section1, &form.Select{
			Name:        "miscellaneous_type",
			Label:       l.T("builder.miscellaneous_type"),
			Value:       publication.MiscellaneousType,
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "miscellaneous_types"),
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/miscellaneous_type"),
		})
	}

	if publication.UsesDOI() {
		section1 = append(section1, &form.Text{
			Name:  "doi",
			Label: l.T("builder.doi"),
			Value: publication.DOI,
			Cols:  9,
			Help:  l.T("builder.doi.help"),
			Error: localize.ValidationErrorAt(l, errors, "/doi"),
		})
	}

	section1 = append(section1, &display.Text{
		Label:   l.T("builder.classification"),
		Value:   l.TS("publication_classifications", publication.Classification),
		Tooltip: l.T("tooltip.publication.classification"),
	})

	if len(section1) > 0 {
		f.AddSection(section1...)
	}

	section2 := []form.Field{}

	if publication.UsesTitle() {
		section2 = append(section2, &form.Text{
			Name:     "title",
			Label:    l.T("builder.title"),
			Value:    publication.Title,
			Cols:     9,
			Error:    localize.ValidationErrorAt(l, errors, "/title"),
			Required: true,
		})
	}

	if publication.UsesAlternativeTitle() {
		section2 = append(section2, &form.TextRepeat{
			Name:   "alternative_title",
			Label:  l.T("builder.alternative_title"),
			Values: publication.AlternativeTitle,
			Cols:   9,
			Error:  localize.ValidationErrorAt(l, errors, "/alternative_title"),
		})
	}

	if publication.UsesPublication() {
		section2 = append(section2, &form.Text{
			Name:     "publication",
			Label:    l.T(fmt.Sprintf("builder.%s.publication", publication.Type)),
			Value:    publication.Publication,
			Cols:     9,
			Required: true,
			Error:    localize.ValidationErrorAt(l, errors, "/publication"),
		})
	}

	if publication.UsesPublicationAbbreviation() {
		section2 = append(section2, &form.Text{
			Name:  "publication_abbreviation",
			Label: l.T(fmt.Sprintf("builder.%s.publication_abbreviation", publication.Type)),
			Value: publication.PublicationAbbreviation,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/publication_abbreviation"),
		})
	}

	if len(section2) > 0 {
		f.AddSection(section2...)
	}

	section3 := []form.Field{}

	if publication.UsesLanguage() {
		section3 = append(section3, &form.SelectRepeat{
			Name:        "language",
			Label:       l.T("builder.language"),
			Options:     localize.LanguageSelectOptions(l),
			Values:      publication.Language,
			EmptyOption: true,
			Cols:        9,
			Error:       localize.ValidationErrorAt(l, errors, "/language"),
		})
	}

	if publication.UsesPublicationStatus() {
		section3 = append(section3, &form.Select{
			Name:        "publication_status",
			Label:       l.T("builder.publication_status"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_publishing_statuses"),
			Value:       publication.PublicationStatus,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/publication_status"),
		})
	}

	section3 = append(section3, &form.Checkbox{
		Name:    "extern",
		Label:   l.T("builder.extern"),
		Value:   "true",
		Checked: publication.Extern,
		Cols:    9,
		Error:   localize.ValidationErrorAt(l, errors, "/extern"),
	})

	if publication.UsesYear() {
		section3 = append(section3, &form.Text{
			Name:     "year",
			Label:    l.T("builder.year"),
			Value:    publication.Year,
			Required: true,
			Cols:     3,
			Help:     l.T("builder.year.help"),
			Error:    localize.ValidationErrorAt(l, errors, "/year"),
		})
	}

	if publication.UsesPublisher() {
		section3 = append(section3,
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
			})
	}

	if len(section3) > 0 {
		f.AddSection(section3...)
	}

	section4 := []form.Field{}

	if publication.UsesSeriesTitle() {
		section4 = append(section4, &form.Text{
			Name:  "series_title",
			Label: l.T("builder.series_title"),
			Value: publication.SeriesTitle,
			Cols:  9,
			Error: localize.ValidationErrorAt(l, errors, "/series_title"),
		})
	}

	if publication.UsesVolume() {
		section4 = append(section4, &form.Text{
			Name:  "volume",
			Label: l.T("builder.volume"),
			Value: publication.Volume,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/volume"),
		})
	}

	if publication.UsesIssue() {
		section4 = append(section4,
			&form.Text{
				Name:  "issue",
				Label: l.T("builder.issue"),
				Value: publication.Issue,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/issue"),
			}, &form.Text{
				Name:  "issue_title",
				Label: l.T("builder.issue_title"),
				Value: publication.IssueTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/issue_title"),
			})
	}

	if publication.UsesEdition() {
		section4 = append(section4, &form.Text{
			Name:  "edition",
			Label: l.T("builder.edition"),
			Value: publication.Edition,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/edition"),
		})
	}

	if publication.UsesPage() {
		section4 = append(section4,
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
		)
	}

	if publication.UsesPageCount() {
		section4 = append(section4, &form.Text{
			Name:  "page_count",
			Label: l.T("builder.page_count"),
			Value: publication.PageCount,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/page_count"),
		})
	}

	if publication.UsesArticleNumber() {
		section4 = append(section4, &form.Text{
			Name:  "article_number",
			Label: l.T("builder.article_number"),
			Value: publication.ArticleNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/article_number"),
		})
	}

	if publication.UsesReportNumber() {
		section4 = append(section4, &form.Text{
			Name:  "report_number",
			Label: l.T("builder.report_number"),
			Value: publication.ReportNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/report_number"),
		})
	}

	if len(section4) > 0 {
		f.AddSection(section4...)
	}

	if publication.UsesDefense() {
		f.AddSection(
			&form.Text{
				Name:     "defense_date",
				Label:    l.T("builder.defense_date"),
				Value:    publication.DefenseDate,
				Required: true,
				Cols:     3,
				Help:     l.T("builder.defense_date.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/defense_date"),
			},
			&form.Text{
				Name:     "defense_time",
				Label:    l.T("builder.defense_time"),
				Value:    publication.DefenseTime,
				Required: true,
				Cols:     3,
				Help:     l.T("builder.defense_time.help"),
				Error:    localize.ValidationErrorAt(l, errors, "/defense_time"),
			},
			&form.Text{
				Name:     "defense_place",
				Label:    l.T("builder.defense_place"),
				Value:    publication.DefensePlace,
				Required: true,
				Cols:     3,
				Error:    localize.ValidationErrorAt(l, errors, "/defense_place"),
			},
		)
	}

	if publication.UsesConfirmations() {
		confirmationOptions := localize.VocabularySelectOptions(l, "confirmations")

		f.AddSection(
			&form.RadioButtonGroup{
				Name:    "has_confidential_data",
				Label:   l.T("builder.has_confidential_data"),
				Value:   publication.HasConfidentialData,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_confidential_data"),
			},
			&form.RadioButtonGroup{
				Name:    "has_patent_application",
				Label:   l.T("builder.has_patent_application"),
				Value:   publication.HasPatentApplication,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_patent_application"),
			},
			&form.RadioButtonGroup{
				Name:    "has_publications_planned",
				Label:   l.T("builder.has_publications_planned"),
				Value:   publication.HasPublicationsPlanned,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_publications_planned"),
			},
			&form.RadioButtonGroup{
				Name:    "has_published_material",
				Label:   l.T("builder.has_published_material"),
				Value:   publication.HasPublishedMaterial,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_published_material"),
			},
		)
	}

	section5 := []form.Field{}

	if publication.UsesWOS() {
		if user.CanCuratePublications() {
			section5 = append(section5, &form.Text{
				Name:    "wos_type",
				Label:   l.T("builder.wos_type"),
				Value:   publication.WOSType,
				Tooltip: l.T("tooltip.publication.wos_type"),
			})
		} else {
			section5 = append(section5, &display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   publication.WOSType,
				Tooltip: l.T("tooltip.publication.wos_type"),
			})
		}

		section5 = append(section5, &form.Text{
			Name:  "wos_id",
			Label: l.T("builder.wos_id"),
			Value: publication.WOSID,
			Cols:  3,
			Help:  l.T("builder.wos_id.help"),
			Error: localize.ValidationErrorAt(l, errors, "/wos_id"),
		})
	}

	if publication.UsesISSN() {
		section5 = append(section5,
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
			})
	}

	if publication.UsesISBN() {
		section5 = append(section5,
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
			})
	}

	if publication.UsesPubMedID() {
		section5 = append(section5, &form.Text{
			Name:  "pubmed_id",
			Label: l.T("builder.pubmed_id"),
			Value: publication.PubMedID,
			Cols:  3,
			Help:  l.T("builder.pubmed_id.help"),
			Error: localize.ValidationErrorAt(l, errors, "/pubmed_id"),
		})
	}

	if publication.UsesArxivID() {
		section5 = append(section5, &form.Text{
			Name:  "arxiv_id",
			Label: l.T("builder.arxiv_id"),
			Value: publication.ArxivID,
			Cols:  3,
			Help:  l.T("builder.arxiv_id.help"),
			Error: localize.ValidationErrorAt(l, errors, "/arxiv_id"),
		})
	}

	if publication.UsesESCIID() {
		section5 = append(section5, &form.Text{
			Name:  "esci_id",
			Label: l.T("builder.esci_id"),
			Value: publication.ESCIID,
			Cols:  3,
			Help:  l.T("builder.esci_id.help"),
			Error: localize.ValidationErrorAt(l, errors, "/esci_id"),
		})
	}

	if len(section5) > 0 {
		f.AddSection(section5...)
	}

	return f
}
