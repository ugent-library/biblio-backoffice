package models

import "time"

type PublicationHits struct {
	Total        int                `json:"total"`
	Page         int                `json:"page"`
	LastPage     int                `json:"last_page"`
	PreviousPage bool               `json:"previous_page"`
	NextPage     bool               `json:"next_page"`
	Hits         []*Publication     `json:"hits"`
	Facets       map[string][]Facet `json:"facets"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty" form:"access_level"`
	ContentType        string     `json:"content_type,omitempty" form:"-"`
	DateCreated        *time.Time `json:"date_created,omitempty" form:"-"`
	DateUpdated        *time.Time `json:"date_updated,omitempty" form:"-"`
	Description        string     `json:"description,omitempty" form:"description"`
	Embargo            string     `json:"embargo,omitempty" form:"embargo"`
	EmbargoTo          string     `json:"embargo_to,omitempty" form:"embargo_to"`
	Filename           string     `json:"file_name,omitempty" form:"-"`
	FileSize           int        `json:"file_size,omitempty" form:"-"`
	ID                 string     `json:"file_id,omitempty" form:"-"`
	PublicationVersion string     `json:"publication_version,omitempty" form:"publication_version"`
	Relation           string     `json:"relation,omitempty" form:"relation"`
	Title              string     `json:"title,omitempty" form:"title"`
	ThumbnailURL       string     `json:"thumbnail_url,omitempty" form:"-"`
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

type PublicationContributor struct {
	ID         string   `json:"id,omitempty"`
	ORCID      string   `json:"orcid,omitempty"`
	UGentID    []string `json:"ugent_id,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	FullName   string   `json:"full_name,omitempty"`
	CreditRole []string `json:"credit_role,omitempty"`
}

type Publication struct {
	Abstract                []Text                   `json:"abstract,omitempty" form:"abstract"`
	AdditionalInfo          string                   `json:"additional_info,omitempty" form:"additional_info"`
	AlternativeTitle        []string                 `json:"alternative_title,omitempty" form:"alternative_title"`
	ArticleNumber           string                   `json:"article_number,omitempty" form:"article_number"`
	ArxivID                 string                   `json:"arxiv_id,omitempty" form:"arxiv_id"`
	Author                  []PublicationContributor `json:"author,omitempty" form:"-"`
	BatchId                 string                   `json:"batch_id,omitempty" form:""`
	Classification          string                   `json:"classification,omitempty" form:"classification"`
	CompletenessScore       int                      `json:"completeness_score,omitempty" form:"-"`
	Conference              PublicationConference    `json:"conference,omitempty" form:"conference"`
	ConferenceType          string                   `json:"conference_type,omitempty" form:"conference_type"`
	CreatorID               string                   `json:"creator_id,omitempty" form:"-"`
	DateCreated             *time.Time               `json:"date_created,omitempty" form:"-"`
	DateUpdated             *time.Time               `json:"date_updated,omitempty" form:"-"`
	DefenseDate             string                   `json:"defense_date,omitempty" form:"defense_date"`
	DefensePlace            string                   `json:"defense_place,omitempty" form:"defense_place"`
	DefenseTime             string                   `json:"defense_time,omitempty" form:"defense_time"`
	Department              []PublicationDepartment  `json:"department,omitempty" form:"-"`
	DOI                     string                   `json:"doi,omitempty" form:"doi"`
	Edition                 string                   `json:"edition,omitempty" form:"edition"`
	Editor                  []PublicationContributor `json:"editor,omitempty" form:"-"`
	EISBN                   []string                 `json:"eisbn,omitempty" form:"eisbn"`
	EISSN                   []string                 `json:"eissn,omitempty" form:"eissn"`
	ESCIID                  string                   `json:"esci_id,omitempty" form:"esci_id"`
	Extern                  bool                     `json:"extern,omitempty" form:"extern"`
	File                    []*PublicationFile       `json:"file,omitempty" form:"-"`
	Handle                  string                   `json:"handle,omitempty" form:"-"`
	HasConfidentialData     string                   `json:"has_confidential_data,omitempty" form:"-"`
	HasPatentApplication    string                   `json:"has_patent_application,omitempty" form:"-"`
	HasPublicationsPlanned  string                   `json:"has_publications_planned,omitempty" form:"-"`
	HasPublishedMaterial    string                   `json:"has_published_material,omitempty" form:"-"`
	ID                      string                   `json:"_id,omitempty" form:"-"`
	ISBN                    []string                 `json:"isbn,omitempty" form:"isbn"`
	ISSN                    []string                 `json:"issn,omitempty" form:"issn"`
	Issue                   string                   `json:"issue,omitempty" form:"issue"`
	IssueTitle              string                   `json:"issue_title,omitempty" form:"issue_title"`
	JournalArticleType      string                   `json:"journal_article_type,omitempty" form:"journal_article_type"`
	Keyword                 []string                 `json:"keyword,omitempty" form:"keyword"`
	Language                []string                 `json:"language,omitempty" form:"language"`
	LaySummary              []Text                   `json:"lay_summary,omitempty" form:"lay_summary"`
	Link                    []PublicationLink        `json:"link,omitempty" form:"-"`
	Locked                  bool                     `json:"locked,omitempty" form:"-"`
	Message                 string                   `json:"message,omitempty" form:"-"`
	MiscellaneousType       string                   `json:"miscellaneous_type,omitempty" form:"miscellaneous_type"`
	PageCount               string                   `json:"page_count,omitempty" form:"page_count"`
	PageFirst               string                   `json:"page_first,omitempty" form:"page_first"`
	PageLast                string                   `json:"page_last,omitempty" form:"page_last"`
	PlaceOfPublication      string                   `json:"place_of_publication,omitempty" form:"place_of_publication"`
	Project                 []PublicationProject     `json:"project,omitempty" form:"-"`
	Publication             string                   `json:"publication,omitempty" form:"publication"`
	PublicationAbbreviation string                   `json:"publication_abbreviation,omitempty" form:"publication_abbreviation"`
	PublicationStatus       string                   `json:"publication_status,omitempty" form:"publication_status"`
	Publisher               string                   `json:"publisher,omitempty" form:"publisher"`
	PubMedID                string                   `json:"pubmed_id,omitempty" form:"pubmed_id"`
	ReportNumber            string                   `json:"report_number,omitempty" form:"report_number"`
	ResearchField           []string                 `json:"research_field,omitempty" form:"research_field"`
	ReviewerNote            string                   `json:"reviewer_note,omitempty" form:"-"`
	ReviewerTags            []string                 `json:"reviewer_tags,omitempty" form:"-"`
	SeriesTitle             string                   `json:"series_title,omitempty" form:"series_title"`
	SourceDB                string                   `json:"source_db,omitempty" form:"-"`
	SourceID                string                   `json:"source_id,omitempty" form:"-"`
	SourceRecord            string                   `json:"source_record,omitempty" form:"-"`
	Status                  string                   `json:"status,omitempty" form:"-"`
	Supervisor              []PublicationContributor `json:"supervisor,omitempty" form:"-"`
	Title                   string                   `json:"title,omitempty" form:"title"`
	Type                    string                   `json:"type,omitempty" form:"-"`
	URL                     string                   `json:"url,omitempty" form:"url"`
	UserID                  string                   `json:"user_id,omitempty" form:"-"`
	Version                 int                      `json:"_version,omitempty" form:"-"`
	Volume                  string                   `json:"volume,omitempty" form:"volume"`
	WOSID                   string                   `json:"wos_id,omitempty" form:"wos_id"`
	WOSType                 string                   `json:"wos_type,omitempty" form:"-"`
	Year                    string                   `json:"year,omitempty" form:"year"`
}

func (p *Publication) OpenAccess() bool {
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
