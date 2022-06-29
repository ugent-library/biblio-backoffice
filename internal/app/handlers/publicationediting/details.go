package publicationediting

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
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
	Classification          string   `form:"classification"`
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
	Publisher               string   `form:"page_count"`
	PubMedID                string   `form:"pubmed_id"`
	ReportNumber            string   `form:"report_number"`
	SeriesTitle             string   `form:"series_title"`
	Title                   string   `form:"title"`
	Type                    string   `form:"-"`
	Volume                  string   `form:"volume"`
	WOSID                   string   `form:"wos_id"`
	WOSType                 string   `form:"-"`
	Year                    string   `form:"year"`
}

func (b *BindDetails) CleanValues() {

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
	b := BindDetails{}

	render.Render(w, "publication/edit_details", YieldEditDetails{
		Context: ctx,
		Form:    FormPublicationDetails(ctx.Locale, b, nil),
	})
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindDetails{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	b.CleanValues()

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := FormPublicationDetails(ctx.Locale, b, validationErrs.(validation.Errors))

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
