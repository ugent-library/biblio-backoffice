package frontoffice

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caltechlibrary/doitools"
	"github.com/iancoleman/strcase"
	"github.com/jpillora/ipfilter"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

const timestampFmt = "2006-01-02 15:04:05"

var licenses = map[string]string{
	"CC0-1.0":          "Creative Commons Public Domain Dedication (CC0 1.0)",
	"CC-BY-4.0":        "Creative Commons Attribution 4.0 International Public License (CC-BY 4.0)",
	"CC-BY-SA-4.0":     "Creative Commons Attribution-ShareAlike 4.0 International Public License (CC BY-SA 4.0)",
	"CC-BY-NC-4.0":     "Creative Commons Attribution-NonCommercial 4.0 International Public License (CC BY-NC 4.0)",
	"CC-BY-ND-4.0":     "Creative Commons Attribution-NoDerivatives 4.0 International Public License (CC BY-ND 4.0)",
	"CC-BY-NC-SA-4.0":  "Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International Public License (CC BY-NC-SA 4.0)",
	"CC-BY-NC-ND-4.0":  "Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International Public License (CC BY-NC-ND 4.0)",
	"InCopyright":      "No license (in copyright)",
	"LicenseNotListed": "A specific license has been chosen by the rights holder. Get in touch with the rights holder for reuse rights.",
	"CopyrightUnknown": "Information pending",
	"":                 "No license (in copyright)",
}

var openLicenses = map[string]struct{}{
	"CC0-1.0":         {},
	"CC-BY-4.0":       {},
	"CC-BY-SA-4.0":    {},
	"CC-BY-NC-4.0":    {},
	"CC-BY-ND-4.0":    {},
	"CC-BY-NC-SA-4.0": {},
	"CC-BY-NC-ND-4.0": {},
}

var hiddenLicenses = map[string]struct{}{
	"InCopyright":      {},
	"CopyrightUnknown": {},
}

type Handler struct {
	handlers.BaseHandler
	Repo      *repositories.Repo
	FileStore backends.FileStore
	IPFilter  *ipfilter.IPFilter
}

// safe basic auth handling
// see https://www.alexedwards.net/blog/basic-authentication-in-go
func (h *Handler) BasicAuth(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(viper.GetString("frontend-username")))
			expectedPasswordHash := sha256.Sum256([]byte(viper.GetString("frontend-password")))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				fn(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type AffiliationPath struct {
	UGentID string `json:"ugent_id,omitempty"`
}

type Affiliation struct {
	Path    []AffiliationPath `json:"path,omitempty"`
	UGentID string            `json:"ugent_id,omitempty"`
}

type Person struct {
	ID         string   `json:"_id,omitempty"`
	CreditRole []string `json:"credit_role,omitempty"`
	Name       string   `json:"name,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
}

type Conference struct {
	Name      string `json:"name,omitempty"`
	Location  string `json:"location,omitempty"`
	Organizer string `json:"organizer,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type Defense struct {
	Date     string `json:"date,omitempty"`
	Location string `json:"location,omitempty"`
}

type Change struct {
	On string `json:"on,omitempty"`
	To string `json:"to,omitempty"`
}

type File struct {
	ID                 string  `json:"_id"`
	Name               string  `json:"name,omitempty"`
	Access             string  `json:"access,omitempty"`
	Change             *Change `json:"change,omitempty"`
	ContentType        string  `json:"content_type,omitempty"`
	Kind               string  `json:"kind,omitempty"`
	PublicationVersion string  `json:"publication_version,omitempty"`
	Size               string  `json:"size,omitempty"`
	SHA256             string  `json:"sha256,omitempty"`
}

type Page struct {
	Count string `json:"count,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
}

type Parent struct {
	Title      string `json:"title,omitempty"`
	ShortTitle string `json:"short_title,omitempty"`
}

type Project struct {
	ID    string `json:"_id"`
	Title string `json:"title,omitempty"`
}

type Publisher struct {
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

type Source struct {
	DB     string `json:"db,omitempty"`
	ID     string `json:"id,omitempty"`
	Record string `json:"record,omitempty"`
}

type Relation struct {
	ID string `json:"_id,omitempty"`
}

type Link struct {
	Access string `json:"access,omitempty"`
	Kind   string `json:"kind,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Text struct {
	Text string `json:"text,omitempty"`
	Lang string `json:"lang,omitempty"`
}

type Publication struct {
	ID                  string        `json:"_id"`
	Abstract            []string      `json:"abstract,omitempty"`
	AbstractFull        []Text        `json:"abstract_full,omitempty"`
	AccessLevel         string        `json:"access_level,omitempty"`
	AdditionalInfo      string        `json:"additional_info,omitempty"`
	Affiliation         []Affiliation `json:"affiliation,omitempty"`
	AlternativeLocation []Link        `json:"alternative_location,omitempty"`
	AlternativeTitle    []string      `json:"alternative_title,omitempty"`
	ArticleNumber       string        `json:"article_number,omitempty"`
	ArticleType         string        `json:"article_type,omitempty"`
	ArxivID             string        `json:"arxiv_id,omitempty"`
	Author              []Person      `json:"author,omitempty"`
	Classification      string        `json:"classification,omitempty"`
	Conference          *Conference   `json:"conference,omitempty"`
	ConferenceType      string        `json:"conference_type,omitempty"`
	CopyrightStatement  string        `json:"copyright_statement,omitempty"`
	CreatedBy           *Person       `json:"created_by,omitempty"`
	DateFrom            string        `json:"date_from"`
	DateCreated         string        `json:"date_created"`
	DateUpdated         string        `json:"date_updated"`
	Defense             *Defense      `json:"defense,omitempty"`
	DOI                 []string      `json:"doi,omitempty"`
	Edition             string        `json:"edition,omitempty"`
	Editor              []Person      `json:"editor,omitempty"`
	ESCIID              string        `json:"esci_id,omitempty"`
	Embargo             string        `json:"embargo,omitempty"`
	EmbargoTo           string        `json:"embargo_to,omitempty"`
	External            int           `json:"external"`
	File                []File        `json:"file,omitempty"`
	Format              []string      `json:"format,omitempty"`
	Handle              string        `json:"handle,omitempty"`
	ISBN                []string      `json:"isbn,omitempty"`
	ISSN                []string      `json:"issn,omitempty"`
	Issue               string        `json:"issue,omitempty"`
	Issuetitle          string        `json:"issue_title,omitempty"`
	Keyword             []string      `json:"keyword,omitempty"`
	Language            []string      `json:"language,omitempty"`
	License             string        `json:"license,omitempty"`
	MiscType            string        `json:"misc_type,omitempty"`
	OtherLicense        string        `json:"other_license,omitempty"`
	Page                *Page         `json:"page,omitempty"`
	Parent              *Parent       `json:"parent,omitempty"`
	Project             []Project     `json:"project,omitempty"`
	Promoter            []Person      `json:"promoter,omitempty"`
	PublicationStatus   string        `json:"publication_status,omitempty"`
	Publisher           *Publisher    `json:"publisher,omitempty"`
	PubMedID            string        `json:"pubmed_id,omitempty"`
	SeriesTitle         string        `json:"series_title,omitempty"`
	Source              *Source       `json:"source,omitempty"`
	Status              string        `json:"status,omitempty"`
	Subject             []string      `json:"subject,omitempty"`
	Title               string        `json:"title,omitempty"`
	Type                string        `json:"type,omitempty"`
	URL                 string        `json:"url,omitempty"`
	Volume              string        `json:"volume,omitempty"`
	WOSID               string        `json:"wos_id,omitempty"`
	WOSType             string        `json:"wos_type,omitempty"`
	Year                string        `json:"year,omitempty"`
	RelatedPublication  []Relation    `json:"related_publication,omitempty"`
	RelatedDataset      []Relation    `json:"related_dataset,omitempty"`
	VABBID              string        `json:"vabb_id,omitempty"`
	VABBType            string        `json:"vabb_type,omitempty"`
	VABBApproved        *int          `json:"vabb_approved,omitempty"`
	VABBYear            []string      `json:"vabb_year,omitempty"`
}

type Hits struct {
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Total  int            `json:"total"`
	Hits   []*Publication `json:"hits"`
}

func mapContributor(c *models.Contributor) *Person {
	return &Person{
		ID:        c.PersonID,
		FirstName: c.FirstName(),
		LastName:  c.LastName(),
		Name:      c.Name(),
	}
}

func (h *Handler) mapPublication(p *models.Publication) *Publication {
	pp := &Publication{
		ID:             p.ID,
		AdditionalInfo: p.AdditionalInfo,
		ArticleNumber:  p.ArticleNumber,
		ArxivID:        p.ArxivID,
		Classification: p.Classification,
		//biblio used librecat's zulu time and splitted them
		//two types of dates in the loop (old: zulu, new: with timestamp)
		DateCreated: p.DateCreated.UTC().Format(timestampFmt),
		DateUpdated: p.DateUpdated.UTC().Format(timestampFmt),
		//date_from used by biblio indexer only
		DateFrom:    p.DateFrom.Format("2006-01-02T15:04:05.000Z"),
		Edition:     p.Edition,
		ESCIID:      p.ESCIID,
		Handle:      p.Handle,
		Issue:       p.Issue,
		Issuetitle:  p.IssueTitle,
		PubMedID:    p.PubMedID,
		SeriesTitle: p.SeriesTitle,
		Title:       p.Title,
		Volume:      p.Volume,
		WOSID:       p.WOSID,
		WOSType:     p.WOSType,
		VABBID:      p.VABBID,
		VABBType:    p.VABBType,
		VABBYear:    p.VABBYear,
	}

	if p.Type != "" {
		v := strcase.ToLowerCamel(p.Type)
		if v == "miscellaneous" {
			v = "misc"
		}
		pp.Type = v
	}

	if p.JournalArticleType != "" {
		pp.ArticleType = strcase.ToLowerCamel(p.JournalArticleType)
	}

	if p.ConferenceType != "" {
		v := strcase.ToLowerCamel(p.ConferenceType)
		if v == "abstract" {
			v = "meetingAbstract"
		}
		pp.ConferenceType = v
	}

	if p.MiscellaneousType != "" {
		pp.MiscType = strcase.ToLowerCamel(p.MiscellaneousType)
	}

	if pp.Type == "conference" && pp.ConferenceType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Proceeding") {
			pp.ConferenceType = "proceedingsPaper"
		} else if strings.Contains(p.WOSType, "Conference Paper") {
			pp.ConferenceType = "conferencePaper"
		} else if strings.Contains(p.WOSType, "Abstract") {
			pp.ConferenceType = "meetingAbstract"
		} else if strings.Contains(p.WOSType, "Other") {
			pp.ConferenceType = "other"
		}
	}
	if pp.Type == "conference" && pp.ConferenceType == "" && p.Classification == "P1" {
		pp.ConferenceType = "proceedingsPaper"
	}

	if pp.Type == "journalArticle" && pp.ArticleType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Article") || strings.Contains(p.WOSType, "Journal Paper") {
			pp.ArticleType = "original"
		} else if strings.Contains(p.WOSType, "Proceedings Paper") {
			pp.ArticleType = "proceedingsPaper"
		} else if strings.Contains(p.WOSType, "Letter") || strings.Contains(p.WOSType, "Note") {
			pp.ArticleType = "letterNote"
		} else if strings.Contains(p.WOSType, "Review") {
			pp.ArticleType = "review"
		}
	}

	if pp.Type == "misc" && pp.MiscType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Book Review") {
			pp.MiscType = "bookReview"
		} else if strings.Contains(p.WOSType, "Theatre Review") {
			pp.MiscType = "theatreReview"
		} else if strings.Contains(p.WOSType, "Correction") {
			pp.MiscType = "correction"
		} else if strings.Contains(p.WOSType, "Editorial Material") {
			pp.MiscType = "editorialMaterial"
		} else if strings.Contains(p.WOSType, "Biographical-Item") || strings.Contains(p.WOSType, "Item About An Individual") {
			pp.MiscType = "biography"
		} else if strings.Contains(p.WOSType, "News Item") {
			pp.MiscType = "newsArticle"
		} else if strings.Contains(p.WOSType, "Bibliography") {
			pp.MiscType = "bibliography"
		} else if strings.Contains(p.WOSType, "Other") {
			pp.MiscType = "other"
		}
	}

	if pp.Type == "misc" {
		if pp.MiscType == "biographicalItem" {
			pp.MiscType = "biography"
		} else if pp.MiscType == "bibliographicalItem" {
			pp.MiscType = "bibliography"
		}
	}

	for _, v := range p.Abstract {
		pp.Abstract = append(pp.Abstract, v.Text)
		pp.AbstractFull = append(pp.AbstractFull, Text{Text: v.Text, Lang: v.Lang})
	}

	for _, rel := range p.RelatedOrganizations {
		aff := Affiliation{UGentID: rel.OrganizationID, Path: make([]AffiliationPath, len(rel.Organization.Tree))}
		for i, t := range rel.Organization.Tree {
			aff.Path[i] = AffiliationPath{UGentID: t.ID}
		}
		pp.Affiliation = append(pp.Affiliation, aff)
	}

	for _, v := range p.Link {
		pp.AlternativeLocation = append(pp.AlternativeLocation, Link{
			URL:    v.URL,
			Access: "open",
			Kind:   "fullText",
		})
	}

	if p.AlternativeTitle != nil {
		pp.AlternativeTitle = append(pp.AlternativeTitle, p.AlternativeTitle...)
	}

	for _, v := range p.Author {
		c := mapContributor(v)
		c.CreditRole = append(c.CreditRole, v.CreditRole...)
		pp.Author = append(pp.Author, *c)
	}

	for _, v := range p.Editor {
		c := mapContributor(v)
		pp.Editor = append(pp.Editor, *c)
	}

	for _, v := range p.Supervisor {
		c := mapContributor(v)
		pp.Promoter = append(pp.Promoter, *c)
	}

	if p.Keyword != nil {
		pp.Keyword = append(pp.Keyword, p.Keyword...)
	}

	if p.ResearchField != nil {
		pp.Subject = append(pp.Subject, p.ResearchField...)
	}

	if p.Status == "private" {
		pp.Status = "unsubmitted"
	} else if p.Status == "deleted" && p.HasBeenPublic {
		pp.Status = "pdeleted"
	} else {
		pp.Status = p.Status
	}

	if p.CreatorID != "" {
		pp.CreatedBy = &Person{ID: p.CreatorID}
	}

	if p.DOI != "" {
		doi, err := doitools.NormalizeDOI(p.DOI)
		if err != nil {
			pp.DOI = append(pp.DOI, p.DOI)
		} else {
			pp.DOI = append(pp.DOI, doi)
		}
	}

	if p.Language != nil {
		pp.Language = append(pp.Language, p.Language...)
	} else if p.Language == nil || len(p.Language) == 0 {
		pp.Language = []string{"und"}
	}

	if validation.IsYear(p.Year) {
		pp.Year = p.Year
	}

	if p.PublicationStatus == "unpublished" {
		pp.PublicationStatus = "unpublished"
	} else if p.PublicationStatus == "accepted" {
		pp.PublicationStatus = "inpress"
	} else {
		pp.PublicationStatus = "published"
	}

	if p.ISSN != nil {
		pp.ISSN = append(pp.ISSN, p.ISSN...)
	}
	if p.EISSN != nil {
		pp.ISSN = append(pp.ISSN, p.EISSN...)
	}
	if p.ISBN != nil {
		pp.ISBN = append(pp.ISBN, p.ISBN...)
	}
	if p.EISBN != nil {
		pp.ISBN = append(pp.ISBN, p.EISBN...)
	}

	if p.PageFirst != "" {
		if pp.Page == nil {
			pp.Page = &Page{}
		}
		pp.Page.First = p.PageFirst
	}
	if p.PageLast != "" {
		if pp.Page == nil {
			pp.Page = &Page{}
		}
		pp.Page.Last = p.PageLast
	}
	if p.PageCount != "" {
		if pp.Page == nil {
			pp.Page = &Page{}
		}
		pp.Page.Count = p.PageCount
	}

	if p.Publication != "" {
		if pp.Parent == nil {
			pp.Parent = &Parent{}
		}
		pp.Parent.Title = p.Publication
	}
	if p.PublicationAbbreviation != "" {
		if pp.Parent == nil {
			pp.Parent = &Parent{}
		}
		pp.Parent.ShortTitle = p.PublicationAbbreviation
	}

	if p.Publisher != "" {
		if pp.Publisher == nil {
			pp.Publisher = &Publisher{}
		}
		pp.Publisher.Name = p.Publisher
	}
	if p.PlaceOfPublication != "" {
		if pp.Publisher == nil {
			pp.Publisher = &Publisher{}
		}
		pp.Publisher.Location = p.PlaceOfPublication
	}

	if p.ConferenceName != "" {
		if pp.Conference == nil {
			pp.Conference = &Conference{}
		}
		pp.Conference.Name = p.ConferenceName
	}
	if p.ConferenceLocation != "" {
		if pp.Conference == nil {
			pp.Conference = &Conference{}
		}
		pp.Conference.Location = p.ConferenceLocation
	}
	if p.ConferenceStartDate != "" {
		if pp.Conference == nil {
			pp.Conference = &Conference{}
		}
		pp.Conference.StartDate = p.ConferenceStartDate
	}
	if p.ConferenceEndDate != "" {
		if pp.Conference == nil {
			pp.Conference = &Conference{}
		}
		pp.Conference.EndDate = p.ConferenceEndDate
	}
	if p.ConferenceOrganizer != "" {
		if pp.Conference == nil {
			pp.Conference = &Conference{}
		}
		pp.Conference.Organizer = p.ConferenceOrganizer
	}

	if validation.IsDate(p.DefenseDate) {
		if pp.Defense == nil {
			pp.Defense = &Defense{}
		}
		pp.Defense.Date = p.DefenseDate
	}
	if p.DefensePlace != "" {
		if pp.Defense == nil {
			pp.Defense = &Defense{}
		}
		pp.Defense.Location = p.DefensePlace
	}

	if p.RelatedProjects != nil {
		pp.Project = make([]Project, len(p.RelatedProjects))
		for i, v := range p.RelatedProjects {
			pp.Project[i] = Project{
				ID:    v.ProjectID,
				Title: v.Project.Title,
			}
		}
	}

	if p.File != nil {
		pp.File = make([]File, len(p.File))
		for i, v := range p.File {
			f := File{
				ID:          v.ID,
				Name:        v.Name,
				Size:        fmt.Sprintf("%d", v.Size),
				ContentType: v.ContentType,
				SHA256:      v.SHA256,
			}

			switch r := v.Relation; r {
			case "main_file":
				f.Kind = "fullText"
			case "table_of_contents":
				f.Kind = "toc"
			case "colophon":
				f.Kind = "colophon"
			case "data_fact_sheet":
				f.Kind = "dataFactsheet"
			case "peer_review_report":
				f.Kind = "peerReviewReport"
			case "agreement":
				f.Kind = "agreement"
			case "supplementary_material":
				f.Kind = "dataset" //was called "data_set" in old LibreCat
			default:
				f.Kind = "fullText"
			}

			switch a := v.AccessLevel; a {
			case "info:eu-repo/semantics/openAccess":
				f.Access = "open"
			case "info:eu-repo/semantics/restrictedAccess":
				f.Access = "restricted"
			case "info:eu-repo/semantics/closedAccess":
				f.Access = "private"
			case "info:eu-repo/semantics/embargoedAccess":
				switch ae := v.AccessLevelDuringEmbargo; ae {
				case "info:eu-repo/semantics/openAccess":
					f.Access = "open"
				case "info:eu-repo/semantics/restrictedAccess":
					f.Access = "restricted"
				case "info:eu-repo/semantics/closedAccess":
					f.Access = "private"
				}

				c := &Change{On: v.EmbargoDate}

				switch ae := v.AccessLevelAfterEmbargo; ae {
				case "info:eu-repo/semantics/openAccess":
					c.To = "open"
				case "info:eu-repo/semantics/restrictedAccess":
					c.To = "restricted"
				case "info:eu-repo/semantics/closedAccess":
					c.To = "private"
				}

				f.Change = c
			}

			if f.Kind == "fullText" {
				f.PublicationVersion = v.PublicationVersion
			}

			pp.File[i] = f
		}

		bestLicense := ""
		for _, f := range p.File {
			if bestLicense == "" {
				if _, isLicense := licenses[f.License]; isLicense {
					bestLicense = f.License
				}
			}
			if _, isOpenLicense := openLicenses[f.License]; isOpenLicense {
				bestLicense = f.License
				break
			}
		}
		pp.CopyrightStatement = licenses[bestLicense]
	}

	if p.SourceDB != "" {
		if pp.Source == nil {
			pp.Source = &Source{}
		}
		pp.Source.DB = p.SourceDB
	}
	if p.SourceID != "" {
		if pp.Source == nil {
			pp.Source = &Source{}
		}
		pp.Source.ID = p.SourceID
	}
	if p.SourceRecord != "" {
		if pp.Source == nil {
			pp.Source = &Source{}
		}
		pp.Source.Record = p.SourceRecord
	}

	if p.RelatedDataset != nil {
		rel_ids := make([]string, 0, len(p.RelatedDataset))
		for _, rd := range p.RelatedDataset {
			rel_ids = append(rel_ids, rd.ID)
		}
		related_datasets, _ := h.Repo.GetDatasets(rel_ids)
		pp.RelatedDataset = make([]Relation, 0, len(related_datasets))
		for _, rd := range related_datasets {
			if rd.Status != "public" {
				continue
			}
			pp.RelatedDataset = append(pp.RelatedDataset, Relation{ID: rd.ID})
		}
	}

	if p.Extern {
		pp.External = 1
	} else {
		pp.External = 0
	}

	if p.VABBID != "" {
		if p.VABBApproved {
			v := 1
			pp.VABBApproved = &v
		} else {
			v := 0
			pp.VABBApproved = &v
		}
	}

	return pp
}

func (h *Handler) mapDataset(d *models.Dataset) *Publication {
	pp := &Publication{
		ID: d.ID,
		//biblio used librecat's zulu time and splitted them
		//two types of dates in the loop (old: zulu, new: with timestamp)
		DateCreated: d.DateCreated.UTC().Format(timestampFmt),
		DateUpdated: d.DateUpdated.UTC().Format(timestampFmt),
		//date_from used by biblio indexer only
		DateFrom:    d.DateFrom.Format("2006-01-02T15:04:05.000Z"),
		AccessLevel: d.AccessLevel,
		Embargo:     d.EmbargoDate,
		EmbargoTo:   d.AccessLevelAfterEmbargo,
		Title:       d.Title,
		Type:        "researchData",
		Year:        d.Year,
		Language:    d.Language,
	}

	if len(d.Link) > 0 {
		pp.URL = d.Link[0].URL
	}

	for _, v := range d.Abstract {
		pp.Abstract = append(pp.Abstract, v.Text)
		pp.AbstractFull = append(pp.AbstractFull, Text{Text: v.Text, Lang: v.Lang})
	}

	for _, rel := range d.RelatedOrganizations {
		aff := Affiliation{UGentID: rel.OrganizationID}
		for i := len(rel.Organization.Tree) - 1; i >= 0; i-- {
			aff.Path = append(aff.Path, AffiliationPath{UGentID: rel.Organization.Tree[i].ID})
		}
		aff.Path = append(aff.Path, AffiliationPath{UGentID: rel.OrganizationID})
		pp.Affiliation = append(pp.Affiliation, aff)
	}

	for _, v := range d.Author {
		c := mapContributor(v)
		pp.Author = append(pp.Author, *c)
	}

	if d.CreatorID != "" {
		pp.CreatedBy = &Person{ID: d.CreatorID}
	}

	if val := d.Identifiers.Get("DOI"); val != "" {
		doi, err := doitools.NormalizeDOI(val)
		if err != nil {
			pp.DOI = append(pp.DOI, val)
		} else {
			pp.DOI = append(pp.DOI, doi)
		}
	}

	if d.Format != nil {
		pp.Format = append(pp.Format, d.Format...)
	}

	if d.Keyword != nil {
		pp.Keyword = append(pp.Keyword, d.Keyword...)
	}

	// hide keywords like LicenseNotListed or UnknownCopyright
	if _, isHidden := hiddenLicenses[d.License]; !isHidden {
		pp.License = d.License
	}
	pp.OtherLicense = d.OtherLicense
	if v, ok := licenses[d.License]; ok {
		pp.CopyrightStatement = v
	}

	if d.RelatedProjects != nil {
		pp.Project = make([]Project, len(d.RelatedProjects))
		for i, v := range d.RelatedProjects {
			pp.Project[i] = Project{
				ID:    v.ProjectID,
				Title: v.Project.Title,
			}
		}
	}

	if d.Publisher != "" {
		if pp.Publisher == nil {
			pp.Publisher = &Publisher{}
		}
		pp.Publisher.Name = d.Publisher
	}

	if d.RelatedPublication != nil {
		rel_ids := make([]string, 0, len(d.RelatedPublication))
		for _, rp := range d.RelatedPublication {
			rel_ids = append(rel_ids, rp.ID)
		}
		related_publications, _ := h.Repo.GetPublications(rel_ids)
		pp.RelatedPublication = make([]Relation, 0, len(related_publications))
		for _, rp := range related_publications {
			if rp.Status != "public" {
				continue
			}
			pp.RelatedPublication = append(pp.RelatedPublication, Relation{ID: rp.ID})
		}
	}

	if d.Status == "private" {
		pp.Status = "unsubmitted"
	} else if d.Status == "deleted" && d.HasBeenPublic {
		pp.Status = "pdeleted"
	} else {
		pp.Status = d.Status
	}

	return pp
}

func (h *Handler) GetPublication(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.GetPublication(bind.PathValues(r).Get("id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}
	j, err := json.Marshal(h.mapPublication(p))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) GetDataset(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.GetDataset(bind.PathValues(r).Get("id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}
	j, err := json.Marshal(h.mapDataset(p))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

type BindGetAll struct {
	Limit        int    `query:"limit"`
	Offset       int    `query:"offset"`
	UpdatedSince string `query:"updated_since"`
}

func (h *Handler) GetAllPublications(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.RequestQuery(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Logger.Errorw("updatedSince error", err)
			render.InternalServerError(w, r, err)
			return
		}
		updatedSince = t.Local()
	}

	mappedHits := &Hits{
		Limit:  b.Limit,
		Offset: b.Offset,
	}

	n, publications, err := h.Repo.PublicationsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Logger.Errorw("select error", err)
		render.InternalServerError(w, r, err)
		return
	}

	mappedHits.Total = n
	mappedHits.Hits = make([]*Publication, 0, len(publications))
	for _, p := range publications {
		mappedHits.Hits = append(mappedHits.Hits, h.mapPublication(p))
	}

	j, err := json.Marshal(mappedHits)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(j)
}

func (h *Handler) GetAllDatasets(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.RequestQuery(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Logger.Errorw("updatedSince error", err)
			render.InternalServerError(w, r, err)
			return
		}
		updatedSince = t.Local()
	}

	mappedHits := &Hits{
		Limit:  b.Limit,
		Offset: b.Offset,
	}

	n, datasets, err := h.Repo.DatasetsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Logger.Errorw("select error", err)
		render.InternalServerError(w, r, err)
		return
	}

	mappedHits.Total = n
	mappedHits.Hits = make([]*Publication, 0, len(datasets))
	for _, d := range datasets {
		mappedHits.Hits = append(mappedHits.Hits, h.mapDataset(d))
	}

	j, err := json.Marshal(mappedHits)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(j)
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	vals := bind.PathValues(r)

	p, err := h.Repo.GetPublication(vals.Get("id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}

	if p.Status != "public" {
		render.Forbidden(w, r)
		return
	}

	f := p.GetFile(vals.Get("file_id"))
	if f == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	accessLevel := f.AccessLevel
	if accessLevel == "info:eu-repo/semantics/embargoedAccess" {
		accessLevel = f.AccessLevelDuringEmbargo
	}

	switch accessLevel {
	case "info:eu-repo/semantics/openAccess":
		// ok
	case "info:eu-repo/semantics/restrictedAccess":
		// check ip
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
			ip = remoteIP
		}
		if !h.IPFilter.Allowed(ip) {
			h.Logger.Warnw("ip not allowed, allowed", "ip", ip, "allowed", viper.GetString("ip-ranges"))
			render.Forbidden(w, r)
			return
		}
	default:
		render.Forbidden(w, r)
		return
	}

	var reader io.ReadCloser
	var readerErr error

	if r.Method == "GET" {
		reader, readerErr = h.FileStore.Get(r.Context(), f.SHA256)
		if readerErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
	}

	responseHeaders := [][]string{}
	responseHeaders = append(responseHeaders, []string{"Content-Type", f.ContentType})
	responseHeaders = append(responseHeaders, []string{"Content-Length", fmt.Sprintf("%d", f.Size)})
	responseHeaders = append(responseHeaders, []string{"Last-Modified", f.DateUpdated.UTC().Format(http.TimeFormat)})
	responseHeaders = append(responseHeaders, []string{"ETag", f.SHA256})
	responseHeaders = append(responseHeaders, []string{"Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name))})

	/*
		Important: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304 dictates that all
		headers in a 304 response should be sent, that would have been sent in 200 response,
		but https://github.com/golang/go/issues/6157 dictates that Content-Length
		and Content-Type should not. Weird.
	*/

	// Status 304: If-Modified-Since (Last-Modified)
	if r.Header.Get("If-Modified-Since") != "" {
		sinceTime, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// http time format does not register milliseconds
		if !f.DateUpdated.Truncate(time.Second).After(sinceTime) {
			for _, pairs := range responseHeaders {
				w.Header().Set(pairs[0], pairs[1])
			}
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Status 304: If-None-Match (ETag)
	if r.Header.Get("If-None-Match") == f.SHA256 {
		for _, pairs := range responseHeaders {
			w.Header().Set(pairs[0], pairs[1])
		}
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Status 200
	for _, pairs := range responseHeaders {
		w.Header().Set(pairs[0], pairs[1])
	}
	w.WriteHeader(http.StatusOK)

	if r.Method == "GET" {
		io.Copy(w, reader)
	}

}
