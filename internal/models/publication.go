package models

import (
	"fmt"
	"strings"
	"time"
)

type PublicationHits struct {
	Pagination
	Hits   []*Publication     `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty" form:"access_level"`
	CCLicense          string     `json:"cc_license,omitempty" form:"cc_license"`
	ContentType        string     `json:"content_type,omitempty" form:"-"`
	DateCreated        *time.Time `json:"date_created,omitempty" form:"-"`
	DateUpdated        *time.Time `json:"date_updated,omitempty" form:"-"`
	Description        string     `json:"description,omitempty" form:"description"`
	Embargo            string     `json:"embargo,omitempty" form:"embargo"`
	EmbargoTo          string     `json:"embargo_to,omitempty" form:"embargo_to"`
	Filename           string     `json:"file_name,omitempty" form:"-"`
	FileSize           int        `json:"file_size,omitempty" form:"-"`
	ID                 string     `json:"file_id,omitempty" form:"-"`
	SHA256             string     `json:"sha256,omitempty" form:"-"`
	NoLicense          bool       `json:"no_license,omitempty" form:"no_license"`
	OtherLicense       string     `json:"other_license,omitempty" form:"other_license"`
	PublicationVersion string     `json:"publication_version,omitempty" form:"publication_version"`
	Relation           string     `json:"relation,omitempty" form:"relation"`
	ThumbnailURL       string     `json:"thumbnail_url,omitempty" form:"-"`
	Title              string     `json:"title,omitempty" form:"title"`
	URL                string     `json:"url,omitempty" form:"-"`
}

type PublicationLink struct {
	URL         string `json:"url,omitempty" form:"url"`
	Relation    string `json:"relation,omitempty" form:"relation"`
	Description string `json:"description,omitempty" form:"description"`
}

type PublicationConference struct {
	Name      string `json:"name,omitempty" form:"name"`
	Location  string `json:"location,omitempty" form:"location"`
	Organizer string `json:"organizer,omitempty" form:"organizer"`
	StartDate string `json:"start_date,omitempty" form:"start_date"`
	EndDate   string `json:"end_date,omitempty" form:"end_date"`
}

type PublicationDepartment struct {
	ID string `json:"_id,omitempty"`
}

type PublicationProject struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PublicationORCIDWork struct {
	ORCID   string `json:"orcid,omitempty"`
	PutCode int    `json:"put_code,omitempty"`
}

type RelatedDataset struct {
	ID string `json:"id,omitempty"`
}

type Publication struct {
	Abstract          []Text                `json:"abstract,omitempty" form:"abstract"`
	AdditionalInfo    string                `json:"additional_info,omitempty" form:"additional_info"`
	AlternativeTitle  []string              `json:"alternative_title,omitempty" form:"alternative_title"`
	ArticleNumber     string                `json:"article_number,omitempty" form:"article_number"`
	ArxivID           string                `json:"arxiv_id,omitempty" form:"arxiv_id"`
	Author            []*Contributor        `json:"author,omitempty" form:"-"`
	BatchID           string                `json:"batch_id,omitempty" form:"-"`
	Classification    string                `json:"classification,omitempty" form:"classification"`
	CompletenessScore int                   `json:"completeness_score,omitempty" form:"-"`
	Conference        PublicationConference `json:"conference,omitempty" form:"conference"`
	ConferenceType    string                `json:"conference_type,omitempty" form:"conference_type"`
	// CreationContext         string                  `json:"creation_context,omitempty" form:"-"`
	CreatorID               string                  `json:"creator_id,omitempty" form:"-"`
	DateCreated             *time.Time              `json:"date_created,omitempty" form:"-"`
	DateUpdated             *time.Time              `json:"date_updated,omitempty" form:"-"`
	DefenseDate             string                  `json:"defense_date,omitempty" form:"defense_date"`
	DefensePlace            string                  `json:"defense_place,omitempty" form:"defense_place"`
	DefenseTime             string                  `json:"defense_time,omitempty" form:"defense_time"`
	Department              []PublicationDepartment `json:"department,omitempty" form:"-"`
	DOI                     string                  `json:"doi,omitempty" form:"doi"`
	Edition                 string                  `json:"edition,omitempty" form:"edition"`
	Editor                  []*Contributor          `json:"editor,omitempty" form:"-"`
	EISBN                   []string                `json:"eisbn,omitempty" form:"eisbn"`
	EISSN                   []string                `json:"eissn,omitempty" form:"eissn"`
	ESCIID                  string                  `json:"esci_id,omitempty" form:"esci_id"`
	Extern                  bool                    `json:"extern,omitempty" form:"extern"`
	File                    []*PublicationFile      `json:"file,omitempty" form:"-"`
	Handle                  string                  `json:"handle,omitempty" form:"-"`
	HasConfidentialData     string                  `json:"has_confidential_data,omitempty" form:"has_confidential_data"`
	HasPatentApplication    string                  `json:"has_patent_application,omitempty" form:"has_patent_application"`
	HasPublicationsPlanned  string                  `json:"has_publications_planned,omitempty" form:"has_publications_planned"`
	HasPublishedMaterial    string                  `json:"has_published_material,omitempty" form:"has_published_material"`
	ID                      string                  `json:"id,omitempty" form:"-"`
	ISBN                    []string                `json:"isbn,omitempty" form:"isbn"`
	ISSN                    []string                `json:"issn,omitempty" form:"issn"`
	Issue                   string                  `json:"issue,omitempty" form:"issue"`
	IssueTitle              string                  `json:"issue_title,omitempty" form:"issue_title"`
	JournalArticleType      string                  `json:"journal_article_type,omitempty" form:"journal_article_type"`
	Keyword                 []string                `json:"keyword,omitempty" form:"keyword"`
	Language                []string                `json:"language,omitempty" form:"language"`
	LaySummary              []Text                  `json:"lay_summary,omitempty" form:"lay_summary"`
	Link                    []PublicationLink       `json:"link,omitempty" form:"-"`
	Locked                  bool                    `json:"locked,omitempty" form:"-"`
	Message                 string                  `json:"message,omitempty" form:"-"`
	MiscellaneousType       string                  `json:"miscellaneous_type,omitempty" form:"miscellaneous_type"`
	ORCIDWork               []PublicationORCIDWork  `json:"orcid_work,omitempty" form:"-"`
	PageCount               string                  `json:"page_count,omitempty" form:"page_count"`
	PageFirst               string                  `json:"page_first,omitempty" form:"page_first"`
	PageLast                string                  `json:"page_last,omitempty" form:"page_last"`
	PlaceOfPublication      string                  `json:"place_of_publication,omitempty" form:"place_of_publication"`
	Project                 []PublicationProject    `json:"project,omitempty" form:"-"`
	Publication             string                  `json:"publication,omitempty" form:"publication"`
	PublicationAbbreviation string                  `json:"publication_abbreviation,omitempty" form:"publication_abbreviation"`
	PublicationStatus       string                  `json:"publication_status,omitempty" form:"publication_status"`
	Publisher               string                  `json:"publisher,omitempty" form:"publisher"`
	PubMedID                string                  `json:"pubmed_id,omitempty" form:"pubmed_id"`
	RelatedDataset          []RelatedDataset        `json:"related_dataset,omitempty" form:"-"`
	ReportNumber            string                  `json:"report_number,omitempty" form:"report_number"`
	ResearchField           []string                `json:"research_field,omitempty" form:"research_field"`
	ReviewerNote            string                  `json:"reviewer_note,omitempty" form:"-"`
	ReviewerTags            []string                `json:"reviewer_tags,omitempty" form:"-"`
	SeriesTitle             string                  `json:"series_title,omitempty" form:"series_title"`
	SourceDB                string                  `json:"source_db,omitempty" form:"-"`
	SourceID                string                  `json:"source_id,omitempty" form:"-"`
	SourceRecord            string                  `json:"source_record,omitempty" form:"-"`
	Status                  string                  `json:"status,omitempty" form:"-"`
	Supervisor              []*Contributor          `json:"supervisor,omitempty" form:"-"`
	Title                   string                  `json:"title,omitempty" form:"title"`
	Type                    string                  `json:"type,omitempty" form:"-"`
	URL                     string                  `json:"url,omitempty" form:"url"`
	UserID                  string                  `json:"user_id,omitempty" form:"-"`
	// Version                 int                     `json:"_version,omitempty" form:"-"`
	Volume  string `json:"volume,omitempty" form:"volume"`
	WOSID   string `json:"wos_id,omitempty" form:"wos_id"`
	WOSType string `json:"wos_type,omitempty" form:"-"`
	Year    string `json:"year,omitempty" form:"year"`
}

func (p *Publication) AccessLevel() string {
	for _, a := range []string{"open_access", "local", "closed"} {
		for _, file := range p.File {
			if file.AccessLevel == a {
				return a
			}
		}
	}
	return ""
}

func (p *Publication) HasRelatedDataset(id string) bool {
	for _, r := range p.RelatedDataset {
		if r.ID == id {
			return true
		}
	}
	return false
}

func (p *Publication) GetFile(id string) *PublicationFile {
	for _, file := range p.File {
		if file.ID == id {
			return file
		}
	}
	return nil
}

func (p *Publication) ThumbnailURL() string {
	for _, file := range p.File {
		if file.ThumbnailURL != "" {
			return file.ThumbnailURL
		}
	}
	return ""
}

func (p *Publication) ClassificationChoices() []string {
	switch p.Type {
	case "journal_article":
		return []string{
			"U",
			"A1",
			"A2",
			"A3",
			"A4",
			"V",
		}
	case "book":
		return []string{
			"U",
			"B1",
		}
	case "book_chapter":
		return []string{
			"U",
			"B2",
		}
	case "book_editor", "issue_editor":
		return []string{
			"U",
			"B3",
		}
	case "conference":
		return []string{
			"U",
			"P1",
			"C1",
			"C3",
		}
	case "dissertation":
		return []string{
			"U",
			"D1",
		}
	case "miscellaneous", "report", "preprint":
		return []string{
			"U",
			"V",
		}
	default:
		return []string{
			"U",
		}
	}
}

func (p *Publication) Contributors(role string) []*Contributor {
	switch role {
	case "author":
		return p.Author
	case "editor":
		return p.Editor
	case "supervisor":
		return p.Supervisor
	default:
		return nil
	}
}

func (p *Publication) SetContributors(role string, c []*Contributor) {
	switch role {
	case "author":
		p.Author = c
	case "editor":
		p.Editor = c
	case "supervisor":
		p.Supervisor = c
	}
}

func (p *Publication) AddContributor(role string, i int, c *Contributor) {
	cc := p.Contributors(role)

	if len(cc) == i {
		p.SetContributors(role, append(cc, c))
		return
	}

	newCC := append(cc[:i+1], cc[i:]...)
	newCC[i] = c
	p.SetContributors(role, newCC)
}

func (p *Publication) RemoveContributor(role string, i int) {
	cc := p.Contributors(role)

	p.SetContributors(role, append(cc[:i], cc[i+1:]...))
}

func (p *Publication) UsesLaySummary() bool {

	switch p.Type {
		case "dissertation":
			return true
		default:
			return false
	}

}

func (p *Publication) UsesConference() bool {
	switch p.Type {
	case "book_chapter", "book_editor", "conference", "issue_editor", "journal_article":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesAuthor() bool {
	switch p.Type {
	case "book", "book_chapter", "conference", "dissertation", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesEditor() bool {
	switch p.Type {
	case "book", "book_chapter", "book_editor", "conference", "issue_editor", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesSupervisor() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) InORCIDWorks(orcidID string) bool {
	for _, w := range p.ORCIDWork {
		if w.ORCID == orcidID {
			return true
		}
	}
	return false
}

func (d *Publication) Validate() (errs ValidationErrors) {
	if d.ID == "" {
		errs = append(errs, ValidationError{
			Pointer: "/id",
			Code:    "required",
		})
	}
	if d.Type == "" {
		errs = append(errs, ValidationError{
			Pointer: "/type",
			Code:    "required",
		})
	}
	if d.Status == "" {
		errs = append(errs, ValidationError{
			Pointer: "/status",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.Title == "" {
		errs = append(errs, ValidationError{
			Pointer: "/title",
			Code:    "required",
		})
	}

	for i, c := range d.Author {
		for _, err := range c.Validate() {
			errs = append(errs, ValidationError{
				Pointer: fmt.Sprintf("/author/%d/%s", i, err.Pointer),
				Code:    err.Code,
			})
		}
	}
	for i, c := range d.Editor {
		for _, err := range c.Validate() {
			errs = append(errs, ValidationError{
				Pointer: fmt.Sprintf("/editor/%d/%s", i, err.Pointer),
				Code:    err.Code,
			})
		}
	}
	for i, c := range d.Supervisor {
		for _, err := range c.Validate() {
			errs = append(errs, ValidationError{
				Pointer: fmt.Sprintf("/supervisor/%d/%s", i, err.Pointer),
				Code:    err.Code,
			})
		}
	}

	return
}

func (d *Publication) Vacuum() {
	d.AdditionalInfo = strings.TrimSpace(d.AdditionalInfo)
	d.AlternativeTitle = vacuumStringSlice(d.AlternativeTitle)
	d.ArticleNumber = strings.TrimSpace(d.ArticleNumber)
	d.ArxivID = strings.TrimSpace(d.ArxivID)
	d.DefenseDate = strings.TrimSpace(d.DefenseDate)
	d.DefensePlace = strings.TrimSpace(d.DefensePlace)
	d.DefenseTime = strings.TrimSpace(d.DefenseTime)
	d.DOI = strings.TrimSpace(d.DOI)
	d.Edition = strings.TrimSpace(d.Edition)
	d.EISBN = vacuumStringSlice(d.EISBN)
	d.EISSN = vacuumStringSlice(d.EISSN)
	d.ESCIID = strings.TrimSpace(d.ESCIID)
	d.ISBN = vacuumStringSlice(d.ISBN)
	d.ISSN = vacuumStringSlice(d.ISSN)
	d.Issue = strings.TrimSpace(d.Issue)
	d.IssueTitle = strings.TrimSpace(d.IssueTitle)
	d.Keyword = vacuumStringSlice(d.Keyword)
	d.PageCount = strings.TrimSpace(d.PageCount)
	d.PageFirst = strings.TrimSpace(d.PageFirst)
	d.PageLast = strings.TrimSpace(d.PageLast)
	d.PlaceOfPublication = strings.TrimSpace(d.PlaceOfPublication)
	d.Publication = strings.TrimSpace(d.Publication)
	d.PublicationAbbreviation = strings.TrimSpace(d.PublicationAbbreviation)
	d.Publisher = strings.TrimSpace(d.Publisher)
	d.PubMedID = strings.TrimSpace(d.PubMedID)
	d.ReportNumber = strings.TrimSpace(d.ReportNumber)
	d.ResearchField = vacuumStringSlice(d.ResearchField)
	d.SeriesTitle = strings.TrimSpace(d.SeriesTitle)
	d.Title = strings.TrimSpace(d.Title)
	d.URL = strings.TrimSpace(d.URL)
	d.Volume = strings.TrimSpace(d.Volume)
	d.WOSID = strings.TrimSpace(d.WOSID)
	d.Year = strings.TrimSpace(d.Year)
}
