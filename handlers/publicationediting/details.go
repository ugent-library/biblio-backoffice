package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindDetails struct {
	AlternativeTitle        []string `form:"alternative_title"`
	ArticleNumber           string   `form:"article_number"`
	ArxivID                 string   `form:"arxiv_id"`
	Classification          string   `form:"classification"`
	ConferenceType          string   `form:"conference_type"`
	DefenseDate             string   `form:"defense_date"`
	DefensePlace            string   `form:"defense_place"`
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
	Legacy                  bool     `form:"legacy"`
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

func EditDetails(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)
	views.ShowModal(publicationviews.EditDetailsDialog(c, p, false, nil)).Render(r.Context(), w)
}

func UpdateDetails(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	var b BindDetails
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p.AlternativeTitle = b.AlternativeTitle
	p.ArticleNumber = b.ArticleNumber
	p.ArxivID = b.ArxivID
	p.ConferenceType = b.ConferenceType
	p.DefenseDate = b.DefenseDate
	p.DefensePlace = b.DefensePlace
	// see https://github.com/ugent-library/biblio-backoffice/issues/1058
	//p.DefenseTime = b.DefenseTime
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

	if c.User.CanCurate() {
		p.Classification = b.Classification
		p.Legacy = b.Legacy
		p.WOSType = b.WOSType
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		views.ReplaceModal(publicationviews.EditDetailsDialog(c, p, false, validationErrs.(*okay.Errors))).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditDetailsDialog(c, p, true, nil)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.DetailsBodySelector, publicationviews.DetailsBody(c, p)).Render(r.Context(), w)
}
