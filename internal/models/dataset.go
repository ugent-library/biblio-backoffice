package models

import "time"

type DatasetHits struct {
	Total        int                `json:"total"`
	Page         int                `json:"page"`
	LastPage     int                `json:"last_page"`
	PreviousPage bool               `json:"previous_page"`
	NextPage     bool               `json:"next_page"`
	Hits         []Dataset          `json:"hits"`
	Facets       map[string][]Facet `json:"facets"`
}

type Dataset struct {
	Abstract          []Text               `json:"abstract,omitempty"`
	CompletenessScore int                  `json:"completeness_score,omitempty"`
	Contributor       []Contributor        `json:"contributors,omitempty"`
	Creator           []Contributor        `json:"author,omitempty"`
	CreatorID         string               `json:"creator_id,omitempty"`
	DateCreated       *time.Time           `json:"date_created,omitempty"`
	DateUpdated       *time.Time           `json:"date_updated,omitempty"`
	DOI               string               `json:"doi,omitempty"`
	ID                string               `json:"_id,omitempty"`
	Keyword           []string             `json:"keyword,omitempty"`
	Locked            bool                 `json:"locked,omitempty"`
	Message           string               `json:"message,omitempty"`
	Project           []PublicationProject `json:"project,omitempty"`
	Publisher         string               `json:"publisher,omitempty"`
	ReviewerNote      string               `json:"reviewer_note,omitempty"`
	ReviewerTags      []string             `json:"reviewer_tags,omitempty"`
	Status            string               `json:"status,omitempty"`
	Title             string               `json:"title,omitempty"`
	Type              string               `json:"type,omitempty"`
	URL               string               `json:"url,omitempty"`
	UserID            string               `json:"user_id,omitempty"`
	Version           int                  `json:"_version,omitempty"`
	Year              string               `json:"year,omitempty"`
}
