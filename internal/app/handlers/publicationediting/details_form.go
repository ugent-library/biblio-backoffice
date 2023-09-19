package publicationediting

import (
	"fmt"
	"html/template"

	"github.com/ugent-library/biblio-backoffice/internal/app/localize"
	"github.com/ugent-library/biblio-backoffice/internal/locale"
	"github.com/ugent-library/biblio-backoffice/internal/render/display"
	"github.com/ugent-library/biblio-backoffice/internal/render/form"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/models"
)

func detailsForm(user *models.User, l *locale.Locale, p *models.Publication, errors validation.Errors) *form.Form {
	f := form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(l, errors))

	section1 := []form.Field{}

	if user.CanChangeType(p) {
		section1 = append(section1, &form.Select{
			Template: "publication/type",
			Name:     "type",
			Label:    l.T("builder.type"),
			Options:  localize.VocabularySelectOptions(l, "publication_types"),
			Value:    p.Type,
			Cols:     3,
			Help:     template.HTML(l.T("builder.type.help")),
			Vars: struct {
				Publication *models.Publication
			}{
				Publication: p,
			},
		})
	} else {
		section1 = append(section1, &display.Text{
			Label:   l.T("builder.type"),
			Value:   l.TS("publication_types", p.Type),
			Tooltip: l.T("tooltip.publication.type"),
		})
	}

	if p.UsesJournalArticleType() {
		section1 = append(section1, &form.Select{
			Name:        "journal_article_type",
			Label:       l.T("builder.journal_article_type"),
			Options:     localize.VocabularySelectOptions(l, "journal_article_types"),
			EmptyOption: true,
			Value:       p.JournalArticleType,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/journal_article_type"),
		})
	}

	if p.UsesConferenceType() {
		section1 = append(section1, &form.Select{
			Name:        "conference_type",
			Label:       l.T("builder.conference_type"),
			Value:       p.ConferenceType,
			Options:     localize.VocabularySelectOptions(l, "conference_types"),
			EmptyOption: true,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/conference_type"),
		})
	}

	if p.UsesMiscellaneousType() {
		section1 = append(section1, &form.Select{
			Name:        "miscellaneous_type",
			Label:       l.T("builder.miscellaneous_type"),
			Value:       p.MiscellaneousType,
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "miscellaneous_types"),
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/miscellaneous_type"),
		})
	}

	if p.UsesDOI() {
		section1 = append(section1, &form.Text{
			Name:  "doi",
			Label: l.T("builder.doi"),
			Value: p.DOI,
			Cols:  9,
			Help:  template.HTML(l.T("builder.doi.help")),
			Error: localize.ValidationErrorAt(l, errors, "/doi"),
		})
	}

	if user.CanCurate() {
		vals := p.ClassificationChoices()
		opts := make([]form.SelectOption, len(vals))
		for i, v := range vals {
			opts[i] = form.SelectOption{
				Value: v,
				Label: l.TS("publication_classifications", v),
			}
		}

		section1 = append(section1, &form.Select{
			Name:    "classification",
			Label:   l.T("builder.classification"),
			Options: opts,
			Value:   p.Classification,
			Cols:    3,
			Error:   localize.ValidationErrorAt(l, errors, "/classification"),
		})
	} else {
		section1 = append(section1, &display.Text{
			Label:   l.T("builder.classification"),
			Value:   l.TS("publication_classifications", p.Classification),
			Tooltip: l.T("tooltip.publication.classification"),
		})
	}

	if user.CanCurate() {
		section1 = append(section1, &form.Checkbox{
			Name:    "legacy",
			Label:   l.T("builder.legacy"),
			Value:   "true",
			Checked: p.Legacy,
			Cols:    9,
			Error:   localize.ValidationErrorAt(l, errors, "/legacy"),
		})
	}

	if len(section1) > 0 {
		f.AddSection(section1...)
	}

	section2 := []form.Field{}

	if p.UsesTitle() {
		section2 = append(section2, &form.Text{
			Name:     "title",
			Label:    l.T("builder.title"),
			Value:    p.Title,
			Cols:     9,
			Error:    localize.ValidationErrorAt(l, errors, "/title"),
			Required: true,
		})
	}

	if p.UsesAlternativeTitle() {
		section2 = append(section2, &form.TextRepeat{
			Name:   "alternative_title",
			Label:  l.T("builder.alternative_title"),
			Values: p.AlternativeTitle,
			Cols:   9,
			Error:  localize.ValidationErrorAt(l, errors, "/alternative_title"),
		})
	}

	if p.UsesPublication() {
		section2 = append(section2, &form.Text{
			Name:     "publication",
			Label:    l.T(fmt.Sprintf("builder.%s.publication", p.Type)),
			Value:    p.Publication,
			Cols:     9,
			Required: p.ShowPublicationAsRequired(),
			Error:    localize.ValidationErrorAt(l, errors, "/publication"),
		})
	}

	if p.UsesPublicationAbbreviation() {
		section2 = append(section2, &form.Text{
			Name:  "publication_abbreviation",
			Label: l.T(fmt.Sprintf("builder.%s.publication_abbreviation", p.Type)),
			Value: p.PublicationAbbreviation,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/publication_abbreviation"),
		})
	}

	if len(section2) > 0 {
		f.AddSection(section2...)
	}

	section3 := []form.Field{}

	if p.UsesLanguage() {
		section3 = append(section3, &form.SelectRepeat{
			Name:        "language",
			Label:       l.T("builder.language"),
			Options:     localize.LanguageSelectOptions(l),
			Values:      p.Language,
			EmptyOption: true,
			Cols:        9,
			Error:       localize.ValidationErrorAt(l, errors, "/language"),
		})
	}

	if p.UsesPublicationStatus() {
		section3 = append(section3, &form.Select{
			Name:        "publication_status",
			Label:       l.T("builder.publication_status"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(l, "publication_publishing_statuses"),
			Value:       p.PublicationStatus,
			Cols:        3,
			Error:       localize.ValidationErrorAt(l, errors, "/publication_status"),
		})
	}

	section3 = append(section3, &form.Checkbox{
		Name:    "extern",
		Label:   l.T("builder.extern"),
		Value:   "true",
		Checked: p.Extern,
		Cols:    9,
		Error:   localize.ValidationErrorAt(l, errors, "/extern"),
	})

	if p.UsesYear() {
		section3 = append(section3, &form.Text{
			Name:     "year",
			Label:    l.T("builder.year"),
			Value:    p.Year,
			Required: true,
			Cols:     3,
			Help:     template.HTML(l.T("builder.year.help")),
			Error:    localize.ValidationErrorAt(l, errors, "/year"),
		})
	}

	if p.UsesPublisher() {
		section3 = append(section3,
			&form.Text{
				Name:  "place_of_publication",
				Label: l.T("builder.place_of_publication"),
				Value: p.PlaceOfPublication,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/place_of_publication"),
			},
			&form.Text{
				Name:  "publisher",
				Label: l.T("builder.publisher"),
				Value: p.Publisher,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/publisher"),
			})
	}

	if len(section3) > 0 {
		f.AddSection(section3...)
	}

	section4 := []form.Field{}

	if p.UsesSeriesTitle() {
		section4 = append(section4, &form.Text{
			Name:  "series_title",
			Label: l.T("builder.series_title"),
			Value: p.SeriesTitle,
			Cols:  9,
			Error: localize.ValidationErrorAt(l, errors, "/series_title"),
		})
	}

	if p.UsesVolume() {
		section4 = append(section4, &form.Text{
			Name:  "volume",
			Label: l.T("builder.volume"),
			Value: p.Volume,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/volume"),
		})
	}

	if p.UsesIssue() {
		section4 = append(section4,
			&form.Text{
				Name:  "issue",
				Label: l.T("builder.issue"),
				Value: p.Issue,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/issue"),
			}, &form.Text{
				Name:  "issue_title",
				Label: l.T("builder.issue_title"),
				Value: p.IssueTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(l, errors, "/issue_title"),
			})
	}

	if p.UsesEdition() {
		section4 = append(section4, &form.Text{
			Name:  "edition",
			Label: l.T("builder.edition"),
			Value: p.Edition,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/edition"),
		})
	}

	if p.UsesPage() {
		section4 = append(section4,
			&form.Text{
				Name:  "page_first",
				Label: l.T("builder.page_first"),
				Value: p.PageFirst,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_first"),
			},
			&form.Text{
				Name:  "page_last",
				Label: l.T("builder.page_last"),
				Value: p.PageLast,
				Cols:  3,
				Error: localize.ValidationErrorAt(l, errors, "/page_last"),
			},
		)
	}

	if p.UsesPageCount() {
		section4 = append(section4, &form.Text{
			Name:  "page_count",
			Label: l.T("builder.page_count"),
			Value: p.PageCount,
			Cols:  3,
			Help:  template.HTML(l.T("builder.page_count.help")),
			Error: localize.ValidationErrorAt(l, errors, "/page_count"),
		})
	}

	if p.UsesArticleNumber() {
		section4 = append(section4, &form.Text{
			Name:  "article_number",
			Label: l.T("builder.article_number"),
			Value: p.ArticleNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/article_number"),
		})
	}

	if p.UsesReportNumber() {
		section4 = append(section4, &form.Text{
			Name:  "report_number",
			Label: l.T("builder.report_number"),
			Value: p.ReportNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(l, errors, "/report_number"),
		})
	}

	if len(section4) > 0 {
		f.AddSection(section4...)
	}

	if p.UsesDefense() {
		f.AddSection(
			&form.Text{
				Name:     "defense_date",
				Label:    l.T("builder.defense_date"),
				Value:    p.DefenseDate,
				Required: p.ShowDefenseAsRequired(),
				Cols:     3,
				Help:     template.HTML(l.T("builder.defense_date.help")),
				Error:    localize.ValidationErrorAt(l, errors, "/defense_date"),
			},
			&form.Text{
				Name:     "defense_place",
				Label:    l.T("builder.defense_place"),
				Value:    p.DefensePlace,
				Required: p.ShowDefenseAsRequired(),
				Cols:     3,
				Error:    localize.ValidationErrorAt(l, errors, "/defense_place"),
			},
		)
	}

	if p.UsesConfirmations() {
		confirmationOptions := localize.VocabularySelectOptions(l, "confirmations")

		f.AddSection(
			&form.RadioButtonGroup{
				Name:    "has_confidential_data",
				Label:   l.T("builder.has_confidential_data"),
				Value:   p.HasConfidentialData,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_confidential_data"),
			},
			&form.RadioButtonGroup{
				Name:    "has_patent_application",
				Label:   l.T("builder.has_patent_application"),
				Value:   p.HasPatentApplication,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_patent_application"),
			},
			&form.RadioButtonGroup{
				Name:    "has_publications_planned",
				Label:   l.T("builder.has_publications_planned"),
				Value:   p.HasPublicationsPlanned,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_publications_planned"),
			},
			&form.RadioButtonGroup{
				Name:    "has_published_material",
				Label:   l.T("builder.has_published_material"),
				Value:   p.HasPublishedMaterial,
				Options: confirmationOptions,
				Error:   localize.ValidationErrorAt(l, errors, "/has_published_material"),
			},
		)
	}

	section5 := []form.Field{}

	if p.UsesWOS() {
		if user.CanCurate() {
			section5 = append(section5, &form.Text{
				Name:    "wos_type",
				Label:   l.T("builder.wos_type"),
				Value:   p.WOSType,
				Tooltip: l.T("tooltip.publication.wos_type"),
			})
		} else {
			section5 = append(section5, &display.Text{
				Label:   l.T("builder.wos_type"),
				Value:   p.WOSType,
				Tooltip: l.T("tooltip.publication.wos_type"),
			})
		}

		section5 = append(section5, &form.Text{
			Name:  "wos_id",
			Label: l.T("builder.wos_id"),
			Value: p.WOSID,
			Cols:  3,
			Help:  template.HTML(l.T("builder.wos_id.help")),
			Error: localize.ValidationErrorAt(l, errors, "/wos_id"),
		})
	}

	if p.UsesISSN() {
		section5 = append(section5,
			&form.TextRepeat{
				Name:   "issn",
				Label:  l.T("builder.issn"),
				Values: p.ISSN,
				Cols:   3,
				Help:   template.HTML(l.T("builder.issn.help")),
				Error:  localize.ValidationErrorAt(l, errors, "/issn"),
			},
			&form.TextRepeat{
				Name:   "eissn",
				Label:  l.T("builder.eissn"),
				Values: p.EISSN,
				Cols:   3,
				Help:   template.HTML(l.T("builder.eissn.help")),
				Error:  localize.ValidationErrorAt(l, errors, "/eissn"),
			})
	}

	if p.UsesISBN() {
		section5 = append(section5,
			&form.TextRepeat{
				Name:   "isbn",
				Label:  l.T("builder.isbn"),
				Values: p.ISBN,
				Cols:   3,
				Help:   template.HTML(l.T("builder.isbn.help")),
				Error:  localize.ValidationErrorAt(l, errors, "/isbn"),
			},
			&form.TextRepeat{
				Name:   "eisbn",
				Label:  l.T("builder.eisbn"),
				Values: p.EISBN,
				Cols:   3,
				Help:   template.HTML(l.T("builder.eisbn.help")),
				Error:  localize.ValidationErrorAt(l, errors, "/eisbn"),
			})
	}

	if p.UsesPubMedID() {
		section5 = append(section5, &form.Text{
			Name:  "pubmed_id",
			Label: l.T("builder.pubmed_id"),
			Value: p.PubMedID,
			Cols:  3,
			Help:  template.HTML(l.T("builder.pubmed_id.help")),
			Error: localize.ValidationErrorAt(l, errors, "/pubmed_id"),
		})
	}

	if p.UsesArxivID() {
		section5 = append(section5, &form.Text{
			Name:  "arxiv_id",
			Label: l.T("builder.arxiv_id"),
			Value: p.ArxivID,
			Cols:  3,
			Help:  template.HTML(l.T("builder.arxiv_id.help")),
			Error: localize.ValidationErrorAt(l, errors, "/arxiv_id"),
		})
	}

	if p.UsesESCIID() {
		section5 = append(section5, &form.Text{
			Name:  "esci_id",
			Label: l.T("builder.esci_id"),
			Value: p.ESCIID,
			Cols:  3,
			Help:  template.HTML(l.T("builder.esci_id.help")),
			Error: localize.ValidationErrorAt(l, errors, "/esci_id"),
		})
	}

	if len(section5) > 0 {
		f.AddSection(section5...)
	}

	return f
}
