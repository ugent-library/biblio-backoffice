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
	AccessLevel string `json:"access_level,omitempty"`
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
	ID                string                   `json:"_id,omitempty"`
	Version           int                      `json:"_version,omitempty"`
	Type              string                   `json:"type,omitempty"`
	DateCreated       *time.Time               `json:"date_created,omitempty"`
	DateUpdated       *time.Time               `json:"date_updated,omitempty"`
	UserID            string                   `json:"user_id,omitempty"`
	Locked            bool                     `json:"locked,omitempty"`
	Extern            bool                     `json:"extern,omitempty"`
	Status            string                   `json:"status,omitempty"`
	PublicationStatus string                   `json:"publication_status,omitempty"`
	File              []PublicationFile        `json:"file,omitempty"`
	Classification    string                   `json:"classification,omitempty"`
	DOI               string                   `json:"doi,omitempty"`
	ISBN              []string                 `json:"isbn,omitempty"`
	ISSN              []string                 `json:"issn,omitempty"`
	Title             string                   `json:"title,omitempty"`
	AlternativeTitle  []string                 `json:"alternative_title,omitempty"`
	Publication       string                   `json:"publication,omitempty"`
	Year              string                   `json:"year,omitempty"`
	Author            []PublicationContributor `json:"author,omitempty"`
	Abstract          []PublicationAbstract    `json:"abstract,omitempty"`
	Conference        PublicationConference    `json:"conference,omitempty"`
}

func (p *Publication) IsOpenAccess() bool {
	for _, file := range p.File {
		if file.AccessLevel == "open_access" {
			return true
		}
	}
	return false
}
