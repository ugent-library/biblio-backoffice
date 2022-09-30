package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/bind"
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
	WOSType                 string   `form:"wos_type"`
	Year                    string   `form:"year"`
}

type YieldDetails struct {
	Context
	DisplayDetails *display.Display
}

type YieldEditDetails struct {
	Context
	Form     *form.Form
	Conflict bool
}

func (h *Handler) EditDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "show_modal", "publication/edit_details", YieldEditDetails{
		Context:  ctx,
		Form:     detailsForm(ctx.User, ctx.Locale, ctx.Publication, nil),
		Conflict: false,
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDetails
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("update publication details: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p := ctx.Publication

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

	if ctx.User.CanCurate() {
		p.WOSType = b.WOSType
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context:  ctx,
			Form:     detailsForm(ctx.User, ctx.Locale, ctx.Publication, validationErrs.(validation.Errors)),
			Conflict: false,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "publication/edit_details", YieldEditDetails{
			Context:  ctx,
			Form:     detailsForm(ctx.User, ctx.Locale, ctx.Publication, nil),
			Conflict: true,
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("update publication details: Could not save the publication:", "error", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_details", YieldDetails{
		Context:        ctx,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, ctx.Publication),
	})
}
