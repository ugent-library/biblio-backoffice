package models

import "time"

type PublicationHits struct {
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	LastPage     int           `json:"last_page"`
	PreviousPage bool          `json:"previous_page"`
	NextPage     bool          `json:"next_page"`
	Hits         []Publication `json:"hits"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty"`
	ContentType        string     `json:"content_type,omitempty"`
	DateCreated        *time.Time `json:"date_created,omitempty"`
	DateUpdated        *time.Time `json:"date_updated,omitempty"`
	Embargo            string     `json:"embargo,omitempty"`
	EmbargoTo          string     `json:"embargo_to,omitempty"`
	FileID             string     `json:"file_id,omitempty"`
	Filename           string     `json:"file_name,omitempty"`
	FileSize           int        `json:"file_size,omitempty"`
	PublicationVersion string     `json:"publication_version,omitempty"`
	Relation           string     `json:"relation,omitempty"`
}

type PublicationContributor struct {
	ID        string `json:"id,omitempty"`
	ORCID     string `json:"orcid,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	FullName  string `json:"full_name,omitempty"`
}

type PublicationConference struct {
	Name      string `json:"name,omitempty"`
	Location  string `json:"location,omitempty"`
	Organizer string `json:"organizer,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type PublicationAbstract struct {
	Text string `json:"text,omitempty"`
	Lang string `json:"lang,omitempty"`
}

type Publication struct {
	Abstract          []PublicationAbstract    `json:"abstract,omitempty"`
	AlternativeTitle  []string                 `json:"alternative_title,omitempty"`
	ArticleNumber     string                   `json:"article_number,omitempty"`
	ArticleType       string                   `json:"article_type,omitempty"`
	Author            []PublicationContributor `json:"author,omitempty"`
	Classification    string                   `json:"classification,omitempty"`
	Conference        PublicationConference    `json:"conference,omitempty"`
	DateCreated       *time.Time               `json:"date_created,omitempty"`
	DateUpdated       *time.Time               `json:"date_updated,omitempty"`
	DOI               string                   `json:"doi,omitempty"`
	Extern            bool                     `json:"extern,omitempty"`
	File              []PublicationFile        `json:"file,omitempty"`
	ID                string                   `json:"_id,omitempty"`
	ISBN              []string                 `json:"isbn,omitempty"`
	ISSN              []string                 `json:"issn,omitempty"`
	Locked            bool                     `json:"locked,omitempty"`
	Publication       string                   `json:"publication,omitempty"`
	PublicationStatus string                   `json:"publication_status,omitempty"`
	Status            string                   `json:"status,omitempty"`
	Title             string                   `json:"title,omitempty"`
	Type              string                   `json:"type,omitempty"`
	UserID            string                   `json:"user_id,omitempty"`
	Version           int                      `json:"_version,omitempty"`
	WOSID             string                   `json:"wos_id,omitempty"`
	WOSType           string                   `json:"wos_type,omitempty"`
	Year              string                   `json:"year,omitempty"`
}

func (p *Publication) IsOpenAccess() bool {
	for _, file := range p.File {
		if file.AccessLevel == "open_access" {
			return true
		}
	}
	return false
}
