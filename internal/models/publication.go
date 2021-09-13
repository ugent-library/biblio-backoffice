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
	AccessLevel string `json:"access_level"`
}

type PublicationContributor struct {
	ID        string `json:"id"`
	ORCID     string `json:"orcid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
}

type Publication struct {
	ID                string                   `json:"_id"`
	Version           int                      `json:"_version"`
	Type              string                   `json:"type"`
	DateCreated       *time.Time               `json:"date_created"`
	DateUpdated       *time.Time               `json:"date_updated"`
	UserID            string                   `json:"user_id"`
	Locked            bool                     `json:"locked"`
	Extern            bool                     `json:"extern"`
	Status            string                   `json:"status"`
	PublicationStatus string                   `json:"publication_status"`
	File              []PublicationFile        `json:"file"`
	Classification    string                   `json:"classification"`
	DOI               string                   `json:"doi"`
	ISBN              []string                 `json:"isbn"`
	ISSN              []string                 `json:"issn"`
	Title             string                   `json:"title"`
	AlternativeTitle  []string                 `json:"alternative_title"`
	Publication       string                   `json:"publication"`
	Year              string                   `json:"year"`
	Author            []PublicationContributor `json:"author"`
}

func (p *Publication) IsOpenAccess() bool {
	for _, file := range p.File {
		if file.AccessLevel == "open_access" {
			return true
		}
	}
	return false
}
