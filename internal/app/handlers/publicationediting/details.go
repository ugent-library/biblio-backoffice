package publicationediting

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type BindDetails struct {
	AlternativeTitle        []string `form:"alternative_title"`
	ArticleNumber           string   `form:"article_number"`
	ArxivID                 string   `form:"arxiv_id"`
	ConferenceType          string   `form:"conference_type"`
	DefenseDate             string   `form:"defense_date"`
	DefensePlace            string   `form:"defense_place"`
	DefenseTime             string   `form:"defense_time"`
	DOI                     string   `form:"doi"`
	Edition                 string   `form:"edition"`
	EISBN                   []string `form:"eisbn"`
	EISSN                   []string `form:"eissn"`
	ESCIID                  string   `form:"esci_id"`
	Extern                  bool     `form:"extern"`
	HasConfidentialData     string   `form:"has_confidential_data"`
	HasPatentApplication    string   `form:"has_patent_application"`
	HasPublicationsPlanned  string   `form:"has_publications_planned"`
	HasPublishedMaterial    string   `form:"has_published_material"`
	ISBN                    []string `form:"isbn"`
	ISSN                    []string `form:"issn"`
	Issue                   string   `form:"issue"`
	IssueTitle              string   `form:"issue_title"`
	JournalArticleType      string   `form:"journal_article_type"`
	Language                []string `form:"language"`
	MiscellaneousType       string   `form:"miscellaneous_type"`
	PageCount               string   `form:"page_count"`
	PageFirst               string   `form:"page_first"`
	PageLast                string   `form:"page_last"`
	PlaceOfPublication      string   `form:"place_of_publication"`
	Publication             string   `form:"publication"`
	PublicationAbbreviation string   `form:"publication_abbreviation"`
	PublicationStatus       string   `form:"publication_status"`
	Publisher               string   `form:"publisher"`
	PubMedID                string   `form:"pubmed_id"`
	ReportNumber            string   `form:"report_number"`
	SeriesTitle             string   `form:"series_title"`
	Title                   string   `form:"title"`
	Volume                  string   `form:"volume"`
	WOSID                   string   `form:"wos_id"`
	Year                    string   `form:"year"`
}

func bindToPublication(b *BindDetails, p *models.Publication) {
	p.AlternativeTitle = b.AlternativeTitle
	p.ArticleNumber = b.ArticleNumber
	p.ArxivID = b.ArxivID
	p.ConferenceType = b.ConferenceType
	p.DefenseDate = b.DefenseDate
	p.DefensePlace = b.DefensePlace
	p.DefenseTime = b.DefenseTime
	p.DOI = b.DOI
	p.Edition = b.Edition
	p.EISBN = b.EISBN
	p.EISSN = b.EISSN
	p.ESCIID = b.ESCIID
	p.Extern = b.Extern
	p.HasConfidentialData = b.HasConfidentialData
	p.HasPatentApplication = b.HasPatentApplication
	p.HasPublicationsPlanned = b.HasPublicationsPlanned
	p.HasPublishedMaterial = b.HasPublishedMaterial
	p.ISBN = b.ISBN
	p.ISSN = b.ISSN
	p.Issue = b.Issue
	p.IssueTitle = b.IssueTitle
	p.JournalArticleType = b.JournalArticleType
	p.Language = b.Language
	p.MiscellaneousType = b.MiscellaneousType
	p.PageCount = b.PageCount
	p.PageFirst = b.PageFirst
	p.PageLast = b.PageLast
	p.PlaceOfPublication = b.PlaceOfPublication
	p.Publication = b.Publication
	p.PublicationAbbreviation = b.PublicationAbbreviation
	p.PublicationStatus = b.PublicationStatus
	p.Publisher = b.Publisher
	p.PubMedID = b.PubMedID
	p.ReportNumber = b.ReportNumber
	p.SeriesTitle = b.SeriesTitle
	p.Title = b.Title
	p.Volume = b.Volume
	p.WOSID = b.WOSID
	p.Year = b.Year
}

func publicationToBind(p *models.Publication, b *BindDetails) {
	b.AlternativeTitle = p.AlternativeTitle
	b.ArticleNumber = p.ArticleNumber
	b.ArxivID = p.ArxivID
	b.ConferenceType = p.ConferenceType
	b.DefenseDate = p.DefenseDate
	b.DefensePlace = p.DefensePlace
	b.DefenseTime = p.DefenseTime
	b.DOI = p.DOI
	b.Edition = p.Edition
	b.EISBN = p.EISBN
	b.EISSN = p.EISSN
	b.ESCIID = p.ESCIID
	b.Extern = p.Extern
	b.HasConfidentialData = p.HasConfidentialData
	b.HasPatentApplication = p.HasPatentApplication
	b.HasPublicationsPlanned = p.HasPublicationsPlanned
	b.HasPublishedMaterial = p.HasPublishedMaterial
	b.ISBN = p.ISBN
	b.ISSN = p.ISSN
	b.Issue = p.Issue
	b.IssueTitle = p.IssueTitle
	b.JournalArticleType = p.JournalArticleType
	b.Language = p.Language
	b.MiscellaneousType = p.MiscellaneousType
	b.PageCount = p.PageCount
	b.PageFirst = p.PageFirst
	b.PageLast = p.PageLast
	b.PlaceOfPublication = p.PlaceOfPublication
	b.Publication = p.Publication
	b.PublicationAbbreviation = p.PublicationAbbreviation
	b.PublicationStatus = p.PublicationStatus
	b.Publisher = p.Publisher
	b.PubMedID = p.PubMedID
	b.ReportNumber = p.ReportNumber
	b.SeriesTitle = p.SeriesTitle
	b.Title = p.Title
	b.Volume = p.Volume
	b.WOSID = p.WOSID
	b.Year = p.Year
}

func (b *BindDetails) cleanValues() {
	/*
		Remove empty values introduced by GUI
		added to the end
	*/
	b.AlternativeTitle = cleanStringSlice(b.AlternativeTitle)
	b.EISBN = cleanStringSlice(b.EISBN)
	b.EISSN = cleanStringSlice(b.EISSN)
	b.ISBN = cleanStringSlice(b.ISBN)
	b.ISSN = cleanStringSlice(b.ISSN)
	b.Language = cleanStringSlice(b.Language)

	// trim spaces (TODO: other?)
	b.ArxivID = strings.TrimSpace(b.ArxivID)
	b.DOI = strings.TrimSpace(b.DOI)
	b.DefenseDate = strings.TrimSpace(b.DefenseDate)
	b.DefensePlace = strings.TrimSpace(b.DefensePlace)
	b.DefenseTime = strings.TrimSpace(b.DefenseTime)
	b.ESCIID = strings.TrimSpace(b.ESCIID)
	b.PubMedID = strings.TrimSpace(b.PubMedID)
	b.Year = strings.TrimSpace(b.Year)
}

type YieldDetails struct {
	Context
	DisplayDetails *display.Display
}

type YieldEditDetails struct {
	Context
	Form *form.Form
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {

	b := &BindDetails{}
	// copy attributes from (current) publication to bind
	publicationToBind(ctx.Publication, b)

	render.Render(w, "publication/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    formPublicationDetails(ctx, b, nil),
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := &BindDetails{}
	if err := bind.Request(r, b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// cleanup values from form
	b.cleanValues()

	/*
		copy attributes from bind to (current) publication
		in order to validate publication
	*/
	bindToPublication(b, ctx.Publication)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := formPublicationDetails(ctx, b, validationErrs.(validation.Errors))

		render.Render(w, "publication/refresh_edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Render(w, "publication/refresh_details", YieldDetails{
		Context:        ctx,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, ctx.Publication),
	})
}

func cleanStringSlice(vals []string) []string {
	var tmp []string
	for _, str := range vals {
		str = strings.TrimSpace(str)
		if str != "" {
			tmp = append(tmp, str)
		}
	}
	return tmp
}

func optionsForVocabulary(locale *locale.Locale, key string) []form.SelectOption {

	options := []form.SelectOption{}
	values, ok := vocabularies.Map[key]

	if !ok {
		return options
	}

	for _, v := range values {
		o := form.SelectOption{}
		o.Label = locale.TS(key, v)
		o.Value = v
		options = append(options, o)
	}

	return options
}
