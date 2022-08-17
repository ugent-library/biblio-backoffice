package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
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

type YieldDetails struct {
	Context
	DisplayDetails *display.Display
}

type YieldEditDetails struct {
	Context
	Form *form.Form
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    detailsForm(ctx.Locale, ctx.Publication, nil),
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := &BindDetails{}
	if err := bind.Request(r, b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	bindToPublication(b, ctx.Publication)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := detailsForm(ctx.Locale, ctx.Publication, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context: ctx,
			Form:    form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_details", YieldDetails{
		Context:        ctx,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, ctx.Publication),
	})
}
