package publicationediting

import (
	"fmt"
	"html/template"

	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render/display"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/validation"
)

func detailsForm(user *models.Person, loc *gotext.Locale, p *models.Publication, errors validation.Errors) *form.Form {
	f := form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(loc, errors))

	section1 := []form.Field{}

	if user.CanChangeType(p) {
		section1 = append(section1, &form.Select{
			Template: "publication/type",
			Name:     "type",
			Label:    loc.Get("builder.type"),
			Options:  localize.VocabularySelectOptions(loc, "publication_types"),
			Value:    p.Type,
			Cols:     3,
			Help:     template.HTML(loc.Get("builder.type.help")),
			Vars: struct {
				Publication *models.Publication
			}{
				Publication: p,
			},
		})
	} else {
		section1 = append(section1, &display.Text{
			Label:   loc.Get("builder.type"),
			Value:   loc.Get("publication_types." + p.Type),
			Tooltip: loc.Get("tooltip.publication.type"),
		})
	}

	if p.UsesJournalArticleType() {
		section1 = append(section1, &form.Select{
			Name:        "journal_article_type",
			Label:       loc.Get("builder.journal_article_type"),
			Options:     localize.VocabularySelectOptions(loc, "journal_article_types"),
			EmptyOption: true,
			Value:       p.JournalArticleType,
			Cols:        3,
			Error:       localize.ValidationErrorAt(loc, errors, "/journal_article_type"),
		})
	}

	if p.UsesConferenceType() {
		section1 = append(section1, &form.Select{
			Name:        "conference_type",
			Label:       loc.Get("builder.conference_type"),
			Value:       p.ConferenceType,
			Options:     localize.VocabularySelectOptions(loc, "conference_types"),
			EmptyOption: true,
			Cols:        3,
			Error:       localize.ValidationErrorAt(loc, errors, "/conference_type"),
		})
	}

	if p.UsesMiscellaneousType() {
		section1 = append(section1, &form.Select{
			Name:        "miscellaneous_type",
			Label:       loc.Get("builder.miscellaneous_type"),
			Value:       p.MiscellaneousType,
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(loc, "miscellaneous_types"),
			Cols:        3,
			Error:       localize.ValidationErrorAt(loc, errors, "/miscellaneous_type"),
		})
	}

	if p.UsesDOI() {
		section1 = append(section1, &form.Text{
			Name:  "doi",
			Label: loc.Get("builder.doi"),
			Value: p.DOI,
			Cols:  9,
			Help:  template.HTML(loc.Get("builder.doi.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/doi"),
		})
	}

	if user.CanCurate() {
		vals := p.ClassificationChoices()
		opts := make([]form.SelectOption, len(vals))
		for i, v := range vals {
			opts[i] = form.SelectOption{
				Value: v,
				Label: loc.Get("publication_classifications." + v),
			}
		}

		section1 = append(section1, &form.Select{
			Name:    "classification",
			Label:   loc.Get("builder.classification"),
			Options: opts,
			Value:   p.Classification,
			Cols:    3,
			Error:   localize.ValidationErrorAt(loc, errors, "/classification"),
		})
	} else {
		section1 = append(section1, &display.Text{
			Label:   loc.Get("builder.classification"),
			Value:   loc.Get("publication_classifications." + p.Classification),
			Tooltip: loc.Get("tooltip.publication.classification"),
		})
	}

	if user.CanCurate() {
		section1 = append(section1, &form.Checkbox{
			Name:    "legacy",
			Label:   loc.Get("builder.legacy"),
			Value:   "true",
			Checked: p.Legacy,
			Cols:    9,
			Error:   localize.ValidationErrorAt(loc, errors, "/legacy"),
		})
	}

	if len(section1) > 0 {
		f.AddSection(section1...)
	}

	section2 := []form.Field{}

	if p.UsesTitle() {
		section2 = append(section2, &form.Text{
			Name:     "title",
			Label:    loc.Get("builder.title"),
			Value:    p.Title,
			Cols:     9,
			Error:    localize.ValidationErrorAt(loc, errors, "/title"),
			Required: true,
		})
	}

	if p.UsesAlternativeTitle() {
		section2 = append(section2, &form.TextRepeat{
			Name:   "alternative_title",
			Label:  loc.Get("builder.alternative_title"),
			Values: p.AlternativeTitle,
			Cols:   9,
			Error:  localize.ValidationErrorAt(loc, errors, "/alternative_title"),
		})
	}

	if p.UsesPublication() {
		section2 = append(section2, &form.Text{
			Name:     "publication",
			Label:    loc.Get(fmt.Sprintf("builder.%s.publication", p.Type)),
			Value:    p.Publication,
			Cols:     9,
			Required: p.ShowPublicationAsRequired(),
			Error:    localize.ValidationErrorAt(loc, errors, "/publication"),
		})
	}

	if p.UsesPublicationAbbreviation() {
		section2 = append(section2, &form.Text{
			Name:  "publication_abbreviation",
			Label: loc.Get(fmt.Sprintf("builder.%s.publication_abbreviation", p.Type)),
			Value: p.PublicationAbbreviation,
			Cols:  3,
			Error: localize.ValidationErrorAt(loc, errors, "/publication_abbreviation"),
		})
	}

	if len(section2) > 0 {
		f.AddSection(section2...)
	}

	section3 := []form.Field{}

	if p.UsesLanguage() {
		section3 = append(section3, &form.SelectRepeat{
			Name:        "language",
			Label:       loc.Get("builder.language"),
			Options:     localize.LanguageSelectOptions(),
			Values:      p.Language,
			EmptyOption: true,
			Cols:        9,
			Error:       localize.ValidationErrorAt(loc, errors, "/language"),
		})
	}

	if p.UsesPublicationStatus() {
		section3 = append(section3, &form.Select{
			Name:        "publication_status",
			Label:       loc.Get("builder.publication_status"),
			EmptyOption: true,
			Options:     localize.VocabularySelectOptions(loc, "publication_publishing_statuses"),
			Value:       p.PublicationStatus,
			Cols:        3,
			Error:       localize.ValidationErrorAt(loc, errors, "/publication_status"),
		})
	}

	section3 = append(section3, &form.Checkbox{
		Name:    "extern",
		Label:   loc.Get("builder.extern"),
		Value:   "true",
		Checked: p.Extern,
		Cols:    9,
		Error:   localize.ValidationErrorAt(loc, errors, "/extern"),
	})

	if p.UsesYear() {
		section3 = append(section3, &form.Text{
			Name:     "year",
			Label:    loc.Get("builder.year"),
			Value:    p.Year,
			Required: true,
			Cols:     3,
			Help:     template.HTML(loc.Get("builder.year.help")),
			Error:    localize.ValidationErrorAt(loc, errors, "/year"),
		})
	}

	if p.UsesPublisher() {
		section3 = append(section3,
			&form.Text{
				Name:  "place_of_publication",
				Label: loc.Get("builder.place_of_publication"),
				Value: p.PlaceOfPublication,
				Cols:  9,
				Error: localize.ValidationErrorAt(loc, errors, "/place_of_publication"),
			},
			&form.Text{
				Name:  "publisher",
				Label: loc.Get("builder.publisher"),
				Value: p.Publisher,
				Cols:  9,
				Error: localize.ValidationErrorAt(loc, errors, "/publisher"),
			})
	}

	if len(section3) > 0 {
		f.AddSection(section3...)
	}

	section4 := []form.Field{}

	if p.UsesSeriesTitle() {
		section4 = append(section4, &form.Text{
			Name:  "series_title",
			Label: loc.Get("builder.series_title"),
			Value: p.SeriesTitle,
			Cols:  9,
			Error: localize.ValidationErrorAt(loc, errors, "/series_title"),
		})
	}

	if p.UsesVolume() {
		section4 = append(section4, &form.Text{
			Name:  "volume",
			Label: loc.Get("builder.volume"),
			Value: p.Volume,
			Cols:  3,
			Error: localize.ValidationErrorAt(loc, errors, "/volume"),
		})
	}

	if p.UsesIssue() {
		section4 = append(section4,
			&form.Text{
				Name:  "issue",
				Label: loc.Get("builder.issue"),
				Value: p.Issue,
				Cols:  3,
				Error: localize.ValidationErrorAt(loc, errors, "/issue"),
			}, &form.Text{
				Name:  "issue_title",
				Label: loc.Get("builder.issue_title"),
				Value: p.IssueTitle,
				Cols:  9,
				Error: localize.ValidationErrorAt(loc, errors, "/issue_title"),
			})
	}

	if p.UsesEdition() {
		section4 = append(section4, &form.Text{
			Name:  "edition",
			Label: loc.Get("builder.edition"),
			Value: p.Edition,
			Cols:  3,
			Error: localize.ValidationErrorAt(loc, errors, "/edition"),
		})
	}

	if p.UsesPage() {
		section4 = append(section4,
			&form.Text{
				Name:  "page_first",
				Label: loc.Get("builder.page_first"),
				Value: p.PageFirst,
				Cols:  3,
				Error: localize.ValidationErrorAt(loc, errors, "/page_first"),
			},
			&form.Text{
				Name:  "page_last",
				Label: loc.Get("builder.page_last"),
				Value: p.PageLast,
				Cols:  3,
				Error: localize.ValidationErrorAt(loc, errors, "/page_last"),
			},
		)
	}

	if p.UsesPageCount() {
		section4 = append(section4, &form.Text{
			Name:  "page_count",
			Label: loc.Get("builder.page_count"),
			Value: p.PageCount,
			Cols:  3,
			Help:  template.HTML(loc.Get("builder.page_count.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/page_count"),
		})
	}

	if p.UsesArticleNumber() {
		section4 = append(section4, &form.Text{
			Name:  "article_number",
			Label: loc.Get("builder.article_number"),
			Value: p.ArticleNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(loc, errors, "/article_number"),
		})
	}

	if p.UsesReportNumber() {
		section4 = append(section4, &form.Text{
			Name:  "report_number",
			Label: loc.Get("builder.report_number"),
			Value: p.ReportNumber,
			Cols:  3,
			Error: localize.ValidationErrorAt(loc, errors, "/report_number"),
		})
	}

	if len(section4) > 0 {
		f.AddSection(section4...)
	}

	if p.UsesDefense() {
		f.AddSection(
			&form.Text{
				Name:     "defense_date",
				Label:    loc.Get("builder.defense_date"),
				Value:    p.DefenseDate,
				Required: p.ShowDefenseAsRequired(),
				Cols:     3,
				Help:     template.HTML(loc.Get("builder.defense_date.help")),
				Error:    localize.ValidationErrorAt(loc, errors, "/defense_date"),
			},
			&form.Text{
				Name:     "defense_place",
				Label:    loc.Get("builder.defense_place"),
				Value:    p.DefensePlace,
				Required: p.ShowDefenseAsRequired(),
				Cols:     3,
				Error:    localize.ValidationErrorAt(loc, errors, "/defense_place"),
			},
		)
	}

	if p.UsesConfirmations() {
		confirmationOptions := localize.VocabularySelectOptions(loc, "confirmations")

		f.AddSection(
			&form.RadioButtonGroup{
				Name:     "has_confidential_data",
				Label:    loc.Get("builder.has_confidential_data"),
				Value:    p.HasConfidentialData,
				Cols:     9,
				Options:  confirmationOptions,
				Error:    localize.ValidationErrorAt(loc, errors, "/has_confidential_data"),
				Required: true,
			},
			&form.RadioButtonGroup{
				Name:     "has_patent_application",
				Label:    loc.Get("builder.has_patent_application"),
				Value:    p.HasPatentApplication,
				Cols:     9,
				Options:  confirmationOptions,
				Error:    localize.ValidationErrorAt(loc, errors, "/has_patent_application"),
				Required: true,
			},
			&form.RadioButtonGroup{
				Name:     "has_publications_planned",
				Label:    loc.Get("builder.has_publications_planned"),
				Value:    p.HasPublicationsPlanned,
				Cols:     9,
				Options:  confirmationOptions,
				Error:    localize.ValidationErrorAt(loc, errors, "/has_publications_planned"),
				Required: true,
			},
			&form.RadioButtonGroup{
				Name:     "has_published_material",
				Label:    loc.Get("builder.has_published_material"),
				Value:    p.HasPublishedMaterial,
				Cols:     9,
				Options:  confirmationOptions,
				Error:    localize.ValidationErrorAt(loc, errors, "/has_published_material"),
				Required: true,
			},
		)
	}

	section5 := []form.Field{}

	if p.UsesWOS() {
		if user.CanCurate() {
			section5 = append(section5, &form.Text{
				Name:    "wos_type",
				Label:   loc.Get("builder.wos_type"),
				Value:   p.WOSType,
				Cols:    3,
				Tooltip: loc.Get("tooltip.publication.wos_type"),
			})
		} else {
			section5 = append(section5, &display.Text{
				Label:   loc.Get("builder.wos_type"),
				Value:   p.WOSType,
				Tooltip: loc.Get("tooltip.publication.wos_type"),
			})
		}

		section5 = append(section5, &form.Text{
			Name:  "wos_id",
			Label: loc.Get("builder.wos_id"),
			Value: p.WOSID,
			Cols:  3,
			Help:  template.HTML(loc.Get("builder.wos_id.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/wos_id"),
		})
	}

	if p.UsesISSN() {
		section5 = append(section5,
			&form.TextRepeat{
				Name:   "issn",
				Label:  loc.Get("builder.issn"),
				Values: p.ISSN,
				Cols:   3,
				Help:   template.HTML(loc.Get("builder.issn.help")),
				Error:  localize.ValidationErrorAt(loc, errors, "/issn"),
			},
			&form.TextRepeat{
				Name:   "eissn",
				Label:  loc.Get("builder.eissn"),
				Values: p.EISSN,
				Cols:   3,
				Help:   template.HTML(loc.Get("builder.eissn.help")),
				Error:  localize.ValidationErrorAt(loc, errors, "/eissn"),
			})
	}

	if p.UsesISBN() {
		section5 = append(section5,
			&form.TextRepeat{
				Name:   "isbn",
				Label:  loc.Get("builder.isbn"),
				Values: p.ISBN,
				Cols:   3,
				Help:   template.HTML(loc.Get("builder.isbn.help")),
				Error:  localize.ValidationErrorAt(loc, errors, "/isbn"),
			},
			&form.TextRepeat{
				Name:   "eisbn",
				Label:  loc.Get("builder.eisbn"),
				Values: p.EISBN,
				Cols:   3,
				Help:   template.HTML(loc.Get("builder.eisbn.help")),
				Error:  localize.ValidationErrorAt(loc, errors, "/eisbn"),
			})
	}

	if p.UsesPubMedID() {
		section5 = append(section5, &form.Text{
			Name:  "pubmed_id",
			Label: loc.Get("builder.pubmed_id"),
			Value: p.PubMedID,
			Cols:  3,
			Help:  template.HTML(loc.Get("builder.pubmed_id.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/pubmed_id"),
		})
	}

	if p.UsesArxivID() {
		section5 = append(section5, &form.Text{
			Name:  "arxiv_id",
			Label: loc.Get("builder.arxiv_id"),
			Value: p.ArxivID,
			Cols:  3,
			Help:  template.HTML(loc.Get("builder.arxiv_id.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/arxiv_id"),
		})
	}

	if p.UsesESCIID() {
		section5 = append(section5, &form.Text{
			Name:  "esci_id",
			Label: loc.Get("builder.esci_id"),
			Value: p.ESCIID,
			Cols:  3,
			Help:  template.HTML(loc.Get("builder.esci_id.help")),
			Error: localize.ValidationErrorAt(loc, errors, "/esci_id"),
		})
	}

	if len(section5) > 0 {
		f.AddSection(section5...)
	}

	return f
}
