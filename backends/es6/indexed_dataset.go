package es6

import (
	"slices"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	internal_time "github.com/ugent-library/biblio-backoffice/time"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type indexedDataset struct {
	AuthorID       []string `json:"author_id,omitempty"`
	BatchID        string   `json:"batch_id,omitempty"`
	CreatorID      string   `json:"creator_id,omitempty"`
	DateCreated    string   `json:"date_created,omitempty"`
	DateUpdated    string   `json:"date_updated,omitempty"`
	Contributor    []string `json:"contributor,omitempty"`
	FacultyID      []string `json:"faculty_id,omitempty"`
	HasMessage     bool     `json:"has_message"`
	ID             string   `json:"id,omitempty"`
	IdentifierType []string `json:"identifier_type,omitempty"`
	Identifier     []string `json:"identifier,omitempty"`
	Keyword        []string `json:"keyword,omitempty"`
	LastUserID     string   `json:"last_user_id,omitempty"`
	Locked         bool     `json:"locked"`
	OrganizationID []string `json:"organization_id,omitempty"`
	Publisher      string   `json:"publisher,omitempty"`
	ReviewerTags   []string `json:"reviewer_tags,omitempty"`
	Status         string   `json:"status,omitempty"`
	Title          string   `json:"title,omitempty"`
	UserID         string   `json:"user_id,omitempty"`
	Year           string   `json:"year,omitempty"`
}

func NewIndexedDataset(d *models.Dataset) *indexedDataset {
	id := &indexedDataset{
		BatchID:      d.BatchID,
		CreatorID:    d.CreatorID,
		DateCreated:  internal_time.FormatTimeUTC(d.DateCreated),
		DateUpdated:  internal_time.FormatTimeUTC(d.DateUpdated),
		HasMessage:   len(d.Message) > 0,
		ID:           d.ID,
		LastUserID:   d.LastUserID,
		Locked:       d.Locked,
		Keyword:      d.Keyword,
		ReviewerTags: d.ReviewerTags,
		Publisher:    d.Publisher,
		Status:       d.Status,
		Title:        d.Title,
		Year:         d.Year,
	}

	faculties := vocabularies.Map["faculties"]

	// extract faculty_id and organization_id from department trees
	for _, rel := range d.RelatedOrganizations {
		for _, org := range rel.Organization.Tree {
			if !slices.Contains(id.OrganizationID, org.ID) {
				id.OrganizationID = append(id.OrganizationID, org.ID)
			}

			if slices.Contains(faculties, org.ID) && !slices.Contains(id.FacultyID, org.ID) {
				id.FacultyID = append(id.FacultyID, org.ID)
			}
		}
	}

	if len(id.FacultyID) == 0 {
		id.FacultyID = append(id.FacultyID, backends.MissingValue)
	}

	for k, vals := range d.Identifiers {
		id.IdentifierType = append(id.IdentifierType, k)
		id.Identifier = append(id.Identifier, vals...)
	}

	for _, author := range d.Author {
		id.Contributor = append(id.Contributor, author.Name())
		if author.PersonID != "" {
			id.AuthorID = append(id.AuthorID, author.PersonID)
		}
	}
	for _, contributor := range d.Contributor {
		id.Contributor = append(id.Contributor, contributor.Name())
	}
	id.AuthorID = lo.Uniq(id.AuthorID)
	id.Contributor = lo.Uniq(id.Contributor)

	return id
}
