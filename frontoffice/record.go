package frontoffice

import (
	"fmt"
	"strconv"
	"strings"

	"slices"

	"github.com/caltechlibrary/doitools"
	"github.com/iancoleman/strcase"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/biblio-backoffice/validation"
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

type AffiliationPath struct {
	UGentID string `json:"ugent_id,omitempty"`
}

type Affiliation struct {
	Path    []AffiliationPath `json:"path,omitempty"`
	UGentID string            `json:"ugent_id,omitempty"`
}

type Person struct {
	ID            string        `json:"_id,omitempty"`
	BiblioID      string        `json:"biblio_id,omitempty"`
	CreditRole    []string      `json:"credit_role,omitempty"`
	Name          string        `json:"name,omitempty"`
	FirstName     string        `json:"first_name,omitempty"`
	LastName      string        `json:"last_name,omitempty"`
	NameLastFirst string        `json:"name_last_first,omitempty"`
	UGentID       []string      `json:"ugent_id,omitempty"`
	ORCID         string        `json:"orcid_id,omitempty"`
	Affiliation   []Affiliation `json:"affiliation,omitempty"`
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
	ID                   string `json:"_id"`
	Title                string `json:"title,omitempty"`
	StartDate            string `json:"start_date,omitempty"`
	EndDate              string `json:"end_date,omitempty"`
	EUID                 string `json:"eu_id,omitempty"`
	EUCallID             string `json:"eu_call_id,omitempty"`
	EUFrameworkProgramme string `json:"eu_framework_programme,omitempty"`
	EUAcronym            string `json:"eu_acronym,omitempty"`
	GISMOID              string `json:"gismo_id,omitempty"`
	IWETOID              string `json:"iweto_id,omitempty"`
	Abstract             string `json:"abstract,omitempty"`
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

type ECOOMFund struct {
	CSS                        string   `json:"css,omitempty"`
	InternationalCollaboration string   `json:"international_collaboration,omitempty"`
	Sector                     []string `json:"sector,omitempty"`
	Validation                 string   `json:"validation,omitempty"`
	Weight                     string   `json:"weight,omitempty"`
}

type JCR struct {
	Eigenfactor          *float64 `json:"eigenfactor,omitempty"`
	ImmediacyIndex       *float64 `json:"immediacy_index,omitempty"`
	ImpactFactor         *float64 `json:"impact_factor,omitempty"`
	ImpactFactor5Year    *float64 `json:"impact_factor_5year,omitempty"`
	TotalCites           *int     `json:"total_cites,omitempty"`
	Category             *string  `json:"category,omitempty"`
	CategoryRank         *string  `json:"category_rank,omitempty"`
	CategoryQuartile     *int     `json:"category_quartile,omitempty"`
	CategoryDecile       *int     `json:"category_decile,omitempty"`
	CategoryVigintile    *int     `json:"category_vigintile,omitempty"`
	PrevImpactFactor     *float64 `json:"prev_impact_factor,omitempty"`
	PrevCategoryQuartile *int     `json:"prev_category_quartile,omitempty"`
}

type Record struct {
	ID                  string               `json:"_id"`
	Abstract            []string             `json:"abstract,omitempty"`
	AbstractFull        []Text               `json:"abstract_full,omitempty"`
	AccessLevel         string               `json:"access_level,omitempty"`
	AdditionalInfo      string               `json:"additional_info,omitempty"`
	Affiliation         []Affiliation        `json:"affiliation,omitempty"`
	AlternativeLocation []Link               `json:"alternative_location,omitempty"`
	AlternativeTitle    []string             `json:"alternative_title,omitempty"`
	ArticleNumber       string               `json:"article_number,omitempty"`
	ArticleType         string               `json:"article_type,omitempty"`
	ArxivID             string               `json:"arxiv_id,omitempty"`
	Author              []Person             `json:"author,omitempty"`
	AuthorSort          string               `json:"author_sort,omitempty"`
	Classification      string               `json:"classification,omitempty"`
	Conference          *Conference          `json:"conference,omitempty"`
	ConferenceType      string               `json:"conference_type,omitempty"`
	CopyrightStatement  string               `json:"copyright_statement,omitempty"`
	CreatedBy           *Person              `json:"created_by,omitempty"`
	DateFrom            string               `json:"date_from"`
	DateCreated         string               `json:"date_created"`
	DateUpdated         string               `json:"date_updated"`
	Defense             *Defense             `json:"defense,omitempty"`
	DOI                 []string             `json:"doi,omitempty"`
	ECOOM               map[string]ECOOMFund `json:"ecoom,omitempty"`
	Edition             string               `json:"edition,omitempty"`
	Editor              []Person             `json:"editor,omitempty"`
	ESCIID              string               `json:"esci_id,omitempty"`
	Embargo             string               `json:"embargo,omitempty"`
	EmbargoTo           string               `json:"embargo_to,omitempty"`
	External            int                  `json:"external"`
	File                []File               `json:"file,omitempty"`
	FirstAuthor         []Person             `json:"first_author,omitempty"`
	Format              []string             `json:"format,omitempty"`
	Handle              string               `json:"handle,omitempty"`
	ISBN                []string             `json:"isbn,omitempty"`
	ISSN                []string             `json:"issn,omitempty"`
	Issue               string               `json:"issue,omitempty"`
	IssueTitle          string               `json:"issue_title,omitempty"`
	JCR                 *JCR                 `json:"jcr,omitempty"`
	Keyword             []string             `json:"keyword,omitempty"`
	Language            []string             `json:"language,omitempty"`
	LastAuthor          []Person             `json:"last_author,omitempty"`
	License             string               `json:"license,omitempty"`
	MiscType            string               `json:"misc_type,omitempty"`
	OtherLicense        string               `json:"other_license,omitempty"`
	Page                *Page                `json:"page,omitempty"`
	Parent              *Parent              `json:"parent,omitempty"`
	Project             []Project            `json:"project,omitempty"`
	Promoter            []Person             `json:"promoter,omitempty"`
	PublicationStatus   string               `json:"publication_status,omitempty"`
	Publisher           *Publisher           `json:"publisher,omitempty"`
	PubMedID            string               `json:"pubmed_id,omitempty"`
	SeriesTitle         string               `json:"series_title,omitempty"`
	SoleAuthor          *Person              `json:"sole_author,omitempty"`
	Source              *Source              `json:"source,omitempty"`
	Status              string               `json:"status,omitempty"`
	Subject             []string             `json:"subject,omitempty"`
	Title               string               `json:"title,omitempty"`
	Type                string               `json:"type,omitempty"`
	URL                 string               `json:"url,omitempty"`
	Volume              string               `json:"volume,omitempty"`
	WOSID               string               `json:"wos_id,omitempty"`
	WOSType             string               `json:"wos_type,omitempty"`
	Year                string               `json:"year,omitempty"`
	RelatedPublication  []Relation           `json:"related_publication,omitempty"`
	RelatedDataset      []Relation           `json:"related_dataset,omitempty"`
	VABBID              string               `json:"vabb_id,omitempty"`
	VABBType            string               `json:"vabb_type,omitempty"`
	VABBApproved        *int                 `json:"vabb_approved,omitempty"`
	VABBYear            []string             `json:"vabb_year,omitempty"`
}

func (r *Record) IsExternal() bool {
	return r.External == 1
}

func (r *Record) IsVABBApproved() bool {
	return r.VABBApproved != nil && *r.VABBApproved == 1
}

type Hits struct {
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	Total  int       `json:"total"`
	Hits   []*Record `json:"hits"`
}

func mapContributor(c *models.Contributor) *Person {
	p := &Person{
		ID:        c.PersonID,
		BiblioID:  c.PersonID,
		FirstName: c.FirstName(),
		LastName:  c.LastName(),
		Name:      c.Name(),
		ORCID:     c.ORCID(),
	}
	if p.LastName != "" && p.FirstName != "" {
		p.NameLastFirst = fmt.Sprintf("%s, %s", p.LastName, p.FirstName)
	}
	if c.Person != nil {
		p.UGentID = c.Person.UGentID
		for _, a := range c.Person.Affiliations {
			aff := Affiliation{UGentID: a.OrganizationID, Path: make([]AffiliationPath, len(a.Organization.Tree))}
			for i, t := range a.Organization.Tree {
				aff.Path[i].UGentID = t.ID
			}
			p.Affiliation = append(p.Affiliation, aff)
		}
	}
	return p
}

func MapPublication(p *models.Publication, repo *repositories.Repo) *Record {
	rec := &Record{
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
		IssueTitle:  p.IssueTitle,
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
		rec.Type = v
	}

	if p.JournalArticleType != "" {
		rec.ArticleType = strcase.ToLowerCamel(p.JournalArticleType)
	}

	if p.ConferenceType != "" {
		v := strcase.ToLowerCamel(p.ConferenceType)
		if v == "abstract" {
			v = "meetingAbstract"
		}
		rec.ConferenceType = v
	}

	if p.MiscellaneousType != "" {
		rec.MiscType = strcase.ToLowerCamel(p.MiscellaneousType)
	}

	if rec.Type == "conference" && rec.ConferenceType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Proceeding") {
			rec.ConferenceType = "proceedingsPaper"
		} else if strings.Contains(p.WOSType, "Conference Paper") {
			rec.ConferenceType = "conferencePaper"
		} else if strings.Contains(p.WOSType, "Abstract") {
			rec.ConferenceType = "meetingAbstract"
		} else if strings.Contains(p.WOSType, "Other") {
			rec.ConferenceType = "other"
		}
	}
	if rec.Type == "conference" && rec.ConferenceType == "" && p.Classification == "P1" {
		rec.ConferenceType = "proceedingsPaper"
	}

	if rec.Type == "journalArticle" && rec.ArticleType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Article") || strings.Contains(p.WOSType, "Journal Paper") {
			rec.ArticleType = "original"
		} else if strings.Contains(p.WOSType, "Proceedings Paper") {
			rec.ArticleType = "proceedingsPaper"
		} else if strings.Contains(p.WOSType, "Letter") || strings.Contains(p.WOSType, "Note") {
			rec.ArticleType = "letterNote"
		} else if strings.Contains(p.WOSType, "Review") {
			rec.ArticleType = "review"
		}
	}

	if rec.Type == "misc" && rec.MiscType == "" && p.WOSType != "" {
		if strings.Contains(p.WOSType, "Book Review") {
			rec.MiscType = "bookReview"
		} else if strings.Contains(p.WOSType, "Theatre Review") {
			rec.MiscType = "theatreReview"
		} else if strings.Contains(p.WOSType, "Correction") {
			rec.MiscType = "correction"
		} else if strings.Contains(p.WOSType, "Editorial Material") {
			rec.MiscType = "editorialMaterial"
		} else if strings.Contains(p.WOSType, "Biographical-Item") || strings.Contains(p.WOSType, "Item About An Individual") {
			rec.MiscType = "biography"
		} else if strings.Contains(p.WOSType, "News Item") {
			rec.MiscType = "newsArticle"
		} else if strings.Contains(p.WOSType, "Bibliography") {
			rec.MiscType = "bibliography"
		} else if strings.Contains(p.WOSType, "Other") {
			rec.MiscType = "other"
		}
	}

	if rec.Type == "misc" {
		if rec.MiscType == "biographicalItem" {
			rec.MiscType = "biography"
		} else if rec.MiscType == "bibliographicalItem" {
			rec.MiscType = "bibliography"
		}
	}

	for _, v := range p.Abstract {
		rec.Abstract = append(rec.Abstract, v.Text)
		rec.AbstractFull = append(rec.AbstractFull, Text{Text: v.Text, Lang: v.Lang})
	}

	for _, rel := range p.RelatedOrganizations {
		aff := Affiliation{UGentID: rel.OrganizationID, Path: make([]AffiliationPath, len(rel.Organization.Tree))}
		for i, t := range rel.Organization.Tree {
			aff.Path[i] = AffiliationPath{UGentID: t.ID}
		}
		rec.Affiliation = append(rec.Affiliation, aff)
	}

	for _, v := range p.Link {
		rec.AlternativeLocation = append(rec.AlternativeLocation, Link{
			URL:    v.URL,
			Access: "open",
			Kind:   "fullText",
		})
	}

	if p.AlternativeTitle != nil {
		rec.AlternativeTitle = append(rec.AlternativeTitle, p.AlternativeTitle...)
	}

	for _, v := range p.Author {
		c := mapContributor(v)
		c.CreditRole = append(c.CreditRole, v.CreditRole...)
		rec.Author = append(rec.Author, *c)
	}

	for _, v := range p.Editor {
		c := mapContributor(v)
		rec.Editor = append(rec.Editor, *c)
	}

	for _, v := range p.Supervisor {
		c := mapContributor(v)
		rec.Promoter = append(rec.Promoter, *c)
	}

	if len(rec.Author) > 0 && rec.Author[0].NameLastFirst != "" {
		rec.AuthorSort = rec.Author[0].NameLastFirst
	}

	if len(rec.Author) == 1 {
		rec.SoleAuthor = &rec.Author[0]
	} else if len(rec.Author) > 1 {
		firstAuthor := make([]Person, 0)
		lastAuthor := make([]Person, 0)
		for _, person := range rec.Author {
			if slices.Contains(person.CreditRole, "first_author") {
				firstAuthor = append(firstAuthor, person)
			}
			if slices.Contains(person.CreditRole, "last_author") {
				lastAuthor = append(lastAuthor, person)
			}
		}
		if len(firstAuthor) == 0 {
			firstAuthor = append(firstAuthor, rec.Author[0])
		}
		if len(lastAuthor) == 0 {
			lastAuthor = append(lastAuthor, rec.Author[len(rec.Author)-1])
		}
		rec.FirstAuthor = firstAuthor
		rec.LastAuthor = lastAuthor
	}

	if p.Keyword != nil {
		rec.Keyword = append(rec.Keyword, p.Keyword...)
	}

	if p.ResearchField != nil {
		rec.Subject = append(rec.Subject, p.ResearchField...)
	}

	if p.Status == "private" {
		rec.Status = "unsubmitted"
	} else if p.Status == "deleted" && p.HasBeenPublic {
		rec.Status = "pdeleted"
	} else {
		rec.Status = p.Status
	}

	if p.CreatorID != "" {
		rec.CreatedBy = mapContributor(&models.Contributor{
			PersonID: p.CreatorID,
			Person:   p.Creator,
		})
	}

	if p.DOI != "" {
		doi, err := doitools.NormalizeDOI(p.DOI)
		if err != nil {
			rec.DOI = append(rec.DOI, p.DOI)
		} else {
			rec.DOI = append(rec.DOI, doi)
		}
	}

	if p.Language != nil {
		rec.Language = append(rec.Language, p.Language...)
	} else if p.Language == nil || len(p.Language) == 0 {
		rec.Language = []string{"und"}
	}

	if validation.IsYear(p.Year) {
		rec.Year = p.Year
	}

	if p.PublicationStatus == "unpublished" {
		rec.PublicationStatus = "unpublished"
	} else if p.PublicationStatus == "accepted" {
		rec.PublicationStatus = "inpress"
	} else {
		rec.PublicationStatus = "published"
	}

	if p.ISSN != nil {
		rec.ISSN = append(rec.ISSN, p.ISSN...)
	}
	if p.EISSN != nil {
		rec.ISSN = append(rec.ISSN, p.EISSN...)
	}
	if p.ISBN != nil {
		rec.ISBN = append(rec.ISBN, p.ISBN...)
	}
	if p.EISBN != nil {
		rec.ISBN = append(rec.ISBN, p.EISBN...)
	}

	if p.PageFirst != "" {
		if rec.Page == nil {
			rec.Page = &Page{}
		}
		rec.Page.First = p.PageFirst
	}
	if p.PageLast != "" {
		if rec.Page == nil {
			rec.Page = &Page{}
		}
		rec.Page.Last = p.PageLast
	}
	if p.PageCount != "" {
		if rec.Page == nil {
			rec.Page = &Page{}
		}
		rec.Page.Count = p.PageCount
	}

	if p.Publication != "" {
		if rec.Parent == nil {
			rec.Parent = &Parent{}
		}
		rec.Parent.Title = p.Publication
	}
	if p.PublicationAbbreviation != "" {
		if rec.Parent == nil {
			rec.Parent = &Parent{}
		}
		rec.Parent.ShortTitle = p.PublicationAbbreviation
	}

	if p.Publisher != "" {
		if rec.Publisher == nil {
			rec.Publisher = &Publisher{}
		}
		rec.Publisher.Name = p.Publisher
	}
	if p.PlaceOfPublication != "" {
		if rec.Publisher == nil {
			rec.Publisher = &Publisher{}
		}
		rec.Publisher.Location = p.PlaceOfPublication
	}

	if p.ConferenceName != "" {
		if rec.Conference == nil {
			rec.Conference = &Conference{}
		}
		rec.Conference.Name = p.ConferenceName
	}
	if p.ConferenceLocation != "" {
		if rec.Conference == nil {
			rec.Conference = &Conference{}
		}
		rec.Conference.Location = p.ConferenceLocation
	}
	if p.ConferenceStartDate != "" {
		if rec.Conference == nil {
			rec.Conference = &Conference{}
		}
		rec.Conference.StartDate = p.ConferenceStartDate
	}
	if p.ConferenceEndDate != "" {
		if rec.Conference == nil {
			rec.Conference = &Conference{}
		}
		rec.Conference.EndDate = p.ConferenceEndDate
	}
	if p.ConferenceOrganizer != "" {
		if rec.Conference == nil {
			rec.Conference = &Conference{}
		}
		rec.Conference.Organizer = p.ConferenceOrganizer
	}

	if validation.IsDate(p.DefenseDate) {
		if rec.Defense == nil {
			rec.Defense = &Defense{}
		}
		rec.Defense.Date = p.DefenseDate
	}
	if p.DefensePlace != "" {
		if rec.Defense == nil {
			rec.Defense = &Defense{}
		}
		rec.Defense.Location = p.DefensePlace
	}

	if p.RelatedProjects != nil {
		rec.Project = make([]Project, len(p.RelatedProjects))
		for i, v := range p.RelatedProjects {
			p := Project{
				ID:        v.ProjectID,
				Title:     v.Project.Title,
				StartDate: v.Project.StartDate,
				EndDate:   v.Project.EndDate,
				GISMOID:   v.Project.GISMOID,
				IWETOID:   v.Project.IWETOID,
			}
			if eu := v.Project.EUProject; eu != nil {
				p.EUID = eu.ID
				p.EUCallID = eu.CallID
				p.EUAcronym = eu.Acronym
				p.EUFrameworkProgramme = eu.FrameworkProgramme
			}
			rec.Project[i] = p
		}
	}

	if p.File != nil {
		rec.File = make([]File, len(p.File))
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

			rec.File[i] = f
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
		rec.CopyrightStatement = licenses[bestLicense]
	}

	if p.SourceDB != "" {
		if rec.Source == nil {
			rec.Source = &Source{}
		}
		rec.Source.DB = p.SourceDB
	}
	if p.SourceID != "" {
		if rec.Source == nil {
			rec.Source = &Source{}
		}
		rec.Source.ID = p.SourceID
	}
	if p.SourceRecord != "" {
		if rec.Source == nil {
			rec.Source = &Source{}
		}
		rec.Source.Record = p.SourceRecord
	}

	if p.RelatedDataset != nil {
		rel_ids := make([]string, 0, len(p.RelatedDataset))
		for _, rd := range p.RelatedDataset {
			rel_ids = append(rel_ids, rd.ID)
		}
		related_datasets, _ := repo.GetDatasets(rel_ids)
		rec.RelatedDataset = make([]Relation, 0, len(related_datasets))
		for _, rd := range related_datasets {
			if rd.Status != "public" {
				continue
			}
			rec.RelatedDataset = append(rec.RelatedDataset, Relation{ID: rd.ID})
		}
	}

	if p.Extern {
		rec.External = 1
	} else {
		rec.External = 0
	}

	if p.VABBID != "" {
		if p.VABBApproved {
			v := 1
			rec.VABBApproved = &v
		} else {
			v := 0
			rec.VABBApproved = &v
		}
	}

	if p.ExternalFields != nil {
		rec.ECOOM = make(map[string]ECOOMFund)
		for _, fund := range []string{"bof", "iof"} {
			fundPrefix := fmt.Sprintf("ecoom-%s", fund)
			fundFields := models.Values{}
			for _, f := range []string{"css", "weight", "sector", "validation", "international_collaboration"} {
				key := fmt.Sprintf("%s-%s", fundPrefix, f)
				vals := p.ExternalFields.GetAll(key)
				if len(vals) > 0 {
					fundFields.SetAll(key, vals...)
				}
			}
			if len(fundFields) == 0 {
				continue
			}
			rec.ECOOM[fund] = ECOOMFund{
				CSS:                        fundFields.Get(fmt.Sprintf("%s-%s", fundPrefix, "css")),
				Weight:                     fundFields.Get(fmt.Sprintf("%s-%s", fundPrefix, "weight")),
				Sector:                     fundFields.GetAll(fmt.Sprintf("%s-%s", fundPrefix, "sector")),
				Validation:                 fundFields.Get(fmt.Sprintf("%s-%s", fundPrefix, "validation")),
				InternationalCollaboration: fundFields.Get(fmt.Sprintf("%s-%s", fundPrefix, "international_collaboration")),
			}
		}
	}

	if p.ExternalFields != nil {
		jcrKeys := []string{
			"eigenfactor", "immediacy_index", "impact_factor", "impact_factor_5year",
			"total_cites", "category", "category_rank", "category_quartile",
			"category_decile", "category_vigintile", "prev_impact_factor",
			"prev_category_quartile"}
		jcrFields := &models.Values{}
		for _, key := range jcrKeys {
			field := "jcr-" + key
			if v := p.ExternalFields.Get(field); v != "" {
				jcrFields.Set(field, v)
			}
		}
		if len(*jcrFields) > 0 {
			rec.JCR = &JCR{}
			if v := jcrFields.Get("jcr-eigenfactor"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					rec.JCR.Eigenfactor = &f
				}
			}
			if v := jcrFields.Get("jcr-immediacy_index"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					rec.JCR.ImmediacyIndex = &f
				}
			}
			if v := jcrFields.Get("jcr-impact_factor"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					rec.JCR.ImpactFactor = &f
				}
			}
			if v := jcrFields.Get("jcr-impact_factor_5year"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					rec.JCR.ImpactFactor5Year = &f
				}
			}
			if v := jcrFields.Get("jcr-total_cites"); v != "" {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					i32 := int(i)
					rec.JCR.TotalCites = &i32
				}
			}
			if v := jcrFields.Get("jcr-category"); v != "" {
				rec.JCR.Category = &v
			}
			if v := jcrFields.Get("jcr-category_rank"); v != "" {
				rec.JCR.CategoryRank = &v
			}
			if v := jcrFields.Get("jcr-category_decile"); v != "" {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					i32 := int(i)
					rec.JCR.CategoryDecile = &i32
				}
			}
			if v := jcrFields.Get("jcr-category_quartile"); v != "" {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					i32 := int(i)
					rec.JCR.CategoryQuartile = &i32
				}
			}
			if v := jcrFields.Get("jcr-category_vigintile"); v != "" {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					i32 := int(i)
					rec.JCR.CategoryVigintile = &i32
				}
			}
			if v := jcrFields.Get("jcr-prev_impact_factor"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					rec.JCR.PrevImpactFactor = &f
				}
			}
			if v := jcrFields.Get("jcr-prev_category_quartile"); v != "" {
				if i, err := strconv.ParseInt(v, 10, 32); err == nil {
					i32 := int(i)
					rec.JCR.PrevCategoryQuartile = &i32
				}
			}
		}

	}

	return rec
}

func MapDataset(d *models.Dataset, repo *repositories.Repo) *Record {
	rec := &Record{
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
		rec.URL = d.Link[0].URL
	}

	for _, v := range d.Abstract {
		rec.Abstract = append(rec.Abstract, v.Text)
		rec.AbstractFull = append(rec.AbstractFull, Text{Text: v.Text, Lang: v.Lang})
	}

	for _, rel := range d.RelatedOrganizations {
		aff := Affiliation{UGentID: rel.OrganizationID}
		for i := len(rel.Organization.Tree) - 1; i >= 0; i-- {
			aff.Path = append(aff.Path, AffiliationPath{UGentID: rel.Organization.Tree[i].ID})
		}
		aff.Path = append(aff.Path, AffiliationPath{UGentID: rel.OrganizationID})
		rec.Affiliation = append(rec.Affiliation, aff)
	}

	for _, v := range d.Author {
		c := mapContributor(v)
		rec.Author = append(rec.Author, *c)
	}

	if d.CreatorID != "" {
		rec.CreatedBy = &Person{ID: d.CreatorID}
	}

	if val := d.Identifiers.Get("DOI"); val != "" {
		doi, err := doitools.NormalizeDOI(val)
		if err != nil {
			rec.DOI = append(rec.DOI, val)
		} else {
			rec.DOI = append(rec.DOI, doi)
		}
	}

	if d.Format != nil {
		rec.Format = append(rec.Format, d.Format...)
	}

	if d.Keyword != nil {
		rec.Keyword = append(rec.Keyword, d.Keyword...)
	}

	// hide keywords like LicenseNotListed or UnknownCopyright
	if _, isHidden := hiddenLicenses[d.License]; !isHidden {
		rec.License = d.License
	}
	rec.OtherLicense = d.OtherLicense
	if v, ok := licenses[d.License]; ok {
		rec.CopyrightStatement = v
	}

	if d.RelatedProjects != nil {
		rec.Project = make([]Project, len(d.RelatedProjects))
		for i, v := range d.RelatedProjects {
			rec.Project[i] = Project{
				ID:    v.ProjectID,
				Title: v.Project.Title,
			}
		}
	}

	if d.Publisher != "" {
		if rec.Publisher == nil {
			rec.Publisher = &Publisher{}
		}
		rec.Publisher.Name = d.Publisher
	}

	if d.RelatedPublication != nil {
		rel_ids := make([]string, 0, len(d.RelatedPublication))
		for _, rp := range d.RelatedPublication {
			rel_ids = append(rel_ids, rp.ID)
		}
		related_publications, _ := repo.GetPublications(rel_ids)
		rec.RelatedPublication = make([]Relation, 0, len(related_publications))
		for _, rp := range related_publications {
			if rp.Status != "public" {
				continue
			}
			rec.RelatedPublication = append(rec.RelatedPublication, Relation{ID: rp.ID})
		}
	}

	if d.Status == "private" {
		rec.Status = "unsubmitted"
	} else if d.Status == "deleted" && d.HasBeenPublic {
		rec.Status = "pdeleted"
	} else {
		rec.Status = d.Status
	}

	return rec
}
