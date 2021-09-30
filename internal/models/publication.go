package models

import "time"

type PublicationHits struct {
	Total        int                `json:"total"`
	Page         int                `json:"page"`
	LastPage     int                `json:"last_page"`
	PreviousPage bool               `json:"previous_page"`
	NextPage     bool               `json:"next_page"`
	Hits         []Publication      `json:"hits"`
	Facets       map[string][]Facet `json:"facets"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty"`
	ContentType        string     `json:"content_type,omitempty"`
	DateCreated        *time.Time `json:"date_created,omitempty"`
	DateUpdated        *time.Time `json:"date_updated,omitempty"`
	Description        string     `json:"description,omitempty"`
	Embargo            string     `json:"embargo,omitempty"`
	EmbargoTo          string     `json:"embargo_to,omitempty"`
	Filename           string     `json:"file_name,omitempty"`
	FileSize           int        `json:"file_size,omitempty"`
	ID                 string     `json:"file_id,omitempty"`
	PublicationVersion string     `json:"publication_version,omitempty"`
	Relation           string     `json:"relation,omitempty"`
	ThumbnailURL       string     `json:"thumbnail_url,omitempty"`
}

type PublicationLink struct {
	URL         string `json:"url,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Description string `json:"description,omitempty"`
}

type PublicationContributor struct {
	ID         string   `json:"id,omitempty"`
	ORCID      string   `json:"orcid,omitempty"`
	UGentID    []string `json:"ugent_id,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	FullName   string   `json:"full_name,omitempty"`
	CreditRole []string `json:"credit_role,omitempty"`
}

type PublicationConference struct {
	Name      string `json:"name,omitempty"`
	Location  string `json:"location,omitempty"`
	Organizer string `json:"organizer,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type PublicationText struct {
	Text string `json:"text,omitempty"`
	Lang string `json:"lang,omitempty"`
}

type PublicationDepartment struct {
	ID string `json:"_id,omitempty"`
}

type PublicationProject struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Publication struct {
	Abstract                []PublicationText        `json:"abstract,omitempty"`
	AdditionalInfo          string                   `json:"additional_info,omitempty"`
	AlternativeTitle        []string                 `json:"alternative_title,omitempty"`
	ArticleNumber           string                   `json:"article_number,omitempty"`
	ArxivID                 string                   `json:"arxiv_id,omitempty"`
	Author                  []PublicationContributor `json:"author,omitempty"`
	Classification          string                   `json:"classification,omitempty"`
	Conference              PublicationConference    `json:"conference,omitempty"`
	ConferenceType          string                   `json:"conference_type,omitempty"`
	CreatorID               string                   `json:"creator_id,omitempty"`
	DateCreated             *time.Time               `json:"date_created,omitempty"`
	DateUpdated             *time.Time               `json:"date_updated,omitempty"`
	DefenseDate             string                   `json:"defense_date,omitempty"`
	DefensePlace            string                   `json:"defense_place,omitempty"`
	DefenseTime             string                   `json:"defense_time,omitempty"`
	Department              []PublicationDepartment  `json:"department,omitempty"`
	DOI                     string                   `json:"doi,omitempty"`
	Edition                 string                   `json:"edition,omitempty"`
	Editor                  []PublicationContributor `json:"editor,omitempty"`
	EISBN                   []string                 `json:"eisbn,omitempty"`
	EISSN                   []string                 `json:"eissn,omitempty"`
	ESCIID                  string                   `json:"esci_id,omitempty"`
	Extern                  bool                     `json:"extern,omitempty"`
	File                    []PublicationFile        `json:"file,omitempty"`
	Handle                  string                   `json:"handle,omitempty"`
	HasConfidentialData     string                   `json:"has_confidential_data,omitempty"`
	HasPatentApplication    string                   `json:"has_patent_application,omitempty"`
	HasPublicationsPlanned  string                   `json:"has_publications_planned,omitempty"`
	HasPublishedMaterial    string                   `json:"has_published_material,omitempty"`
	ID                      string                   `json:"_id,omitempty"`
	ISBN                    []string                 `json:"isbn,omitempty"`
	ISSN                    []string                 `json:"issn,omitempty"`
	Issue                   string                   `json:"issue,omitempty"`
	IssueTitle              string                   `json:"issue_title,omitempty"`
	JournalArticleType      string                   `json:"journal_article_type,omitempty"`
	Keyword                 []string                 `json:"keyword,omitempty"`
	Language                []string                 `json:"language,omitempty"`
	LaySummary              []PublicationText        `json:"lay_summary,omitempty"`
	Link                    []PublicationLink        `json:"link,omitempty"`
	Locked                  bool                     `json:"locked,omitempty"`
	Message                 string                   `json:"message,omitempty"`
	MiscellaneousType       string                   `json:"miscellaneous_type,omitempty"`
	PageCount               string                   `json:"page_count,omitempty"`
	PageFirst               string                   `json:"page_first,omitempty"`
	PageLast                string                   `json:"page_last,omitempty"`
	Project                 []PublicationProject     `json:"project,omitempty"`
	Publication             string                   `json:"publication,omitempty"`
	PublicationAbbreviation string                   `json:"publication_abbreviation,omitempty"`
	PublicationStatus       string                   `json:"publication_status,omitempty"`
	PubMedID                string                   `json:"pubmed_id,omitempty"`
	ReviewerNote            string                   `json:"reviewer_note,omitempty"`
	ReviewerTags            []string                 `json:"reviewer_tags,omitempty"`
	SourceDB                string                   `json:"source_db,omitempty"`
	SourceID                string                   `json:"source_id,omitempty"`
	SourceRecord            string                   `json:"source_record,omitempty"`
	SeriesTitle             string                   `json:"series_title,omitempty"`
	Status                  string                   `json:"status,omitempty"`
	ReportNumber            string                   `json:"report_number,omitempty"`
	ResearchField           []string                 `json:"research_field,omitempty"`
	Supervisor              []PublicationContributor `json:"supervisor,omitempty"`
	Title                   string                   `json:"title,omitempty"`
	Type                    string                   `json:"type,omitempty"`
	UserID                  string                   `json:"user_id,omitempty"`
	Version                 int                      `json:"_version,omitempty"`
	Volume                  string                   `json:"volume,omitempty"`
	WOSID                   string                   `json:"wos_id,omitempty"`
	WOSType                 string                   `json:"wos_type,omitempty"`
	Year                    string                   `json:"year,omitempty"`
	Publisher               string                   `json:"publisher,omitempty"`
	PlaceOfPublication      string                   `json:"place_of_publication,omitempty"`
}

func (p *Publication) IsOpenAccess() bool {
	for _, file := range p.File {
		if file.AccessLevel == "open_access" {
			return true
		}
	}
	return false
}

func (p *Publication) ThumbnailURL() string {
	for _, file := range p.File {
		if file.ThumbnailURL != "" {
			return file.ThumbnailURL
		}
	}
	return ""
}
