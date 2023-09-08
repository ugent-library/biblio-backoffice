package es6

import (
	"regexp"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/models"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
	"github.com/ugent-library/biblio-backoffice/internal/util"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

type indexedPublication struct {
	AuthorID          []string `json:"author_id,omitempty"`
	BatchID           string   `json:"batch_id,omitempty"`
	Classification    string   `json:"classification,omitempty"`
	Contributor       []string `json:"contributor,omitempty"`
	CreatorID         string   `json:"creator_id,omitempty"`
	DateCreated       string   `json:"date_created"`
	DateUpdated       string   `json:"date_updated"`
	DOI               string   `json:"doi,omitempty"`
	Extern            bool     `json:"extern"`
	FacultyID         []string `json:"faculty_id,omitempty"`
	FileRelation      []string `json:"file_relation,omitempty"`
	HasMessage        bool     `json:"has_message"`
	HasFiles          bool     `json:"has_files"`
	ID                string   `json:"id,omitempty"`
	Identifier        []string `json:"identifier,omitempty"`
	ISXN              []string `json:"isxn,omitempty"`
	LastUserID        string   `json:"last_user_id,omitempty"`
	Locked            bool     `json:"locked"`
	OrganizationID    []string `json:"organization_id,omitempty"`
	PublicationStatus string   `json:"publication_status,omitempty"`
	ReviewerTags      []string `json:"reviewer_tags,omitempty"`
	Status            string   `json:"status,omitempty"`
	Title             string   `json:"title,omitempty"`
	Type              string   `json:"type,omitempty"`
	UserID            string   `json:"user_id,omitempty"`
	VABBType          string   `json:"vabb_type,omitempty"`
	WOSType           []string `json:"wos_type,omitempty"`
	Year              string   `json:"year,omitempty"`
}

var reSplitWOS *regexp.Regexp = regexp.MustCompile("[,;]")

func NewIndexedPublication(p *models.Publication) *indexedPublication {
	ip := &indexedPublication{
		BatchID:           p.BatchID,
		Classification:    p.Classification,
		CreatorID:         p.CreatorID,
		DateCreated:       internal_time.FormatTimeUTC(p.DateCreated),
		DateUpdated:       internal_time.FormatTimeUTC(p.DateUpdated),
		DOI:               p.DOI,
		Extern:            p.Extern,
		ID:                p.ID,
		LastUserID:        p.LastUserID,
		Locked:            p.Locked,
		HasMessage:        len(p.Message) > 0,
		PublicationStatus: p.PublicationStatus,
		ReviewerTags:      p.ReviewerTags,
		Status:            p.Status,
		Title:             p.Title,
		Type:              p.Type,
		UserID:            p.UserID,
		VABBType:          p.VABBType,
		Year:              p.Year,
		HasFiles:          len(p.File) > 0,
	}

	faculties := vocabularies.Map["faculties"]

	// extract faculty_id and all organization id's from organization trees
	for _, rel := range p.RelatedOrganizations {
		for _, org := range rel.Organization.Tree {
			if !validation.InArray(ip.OrganizationID, org.ID) {
				ip.OrganizationID = append(ip.OrganizationID, org.ID)
			}

			if validation.InArray(faculties, org.ID) && !validation.InArray(ip.FacultyID, org.ID) {
				ip.FacultyID = append(ip.FacultyID, org.ID)
			}
		}
	}

	if len(ip.FacultyID) == 0 {
		ip.FacultyID = append(ip.FacultyID, models.MissingValue)
	}

	if ip.PublicationStatus == "" {
		ip.PublicationStatus = models.MissingValue
	}

	if p.WOSType != "" {
		wos_types := reSplitWOS.Split(p.WOSType, -1)
		for _, wos_type := range wos_types {
			wt := strings.TrimSpace(wos_type)
			ip.WOSType = append(ip.WOSType, wt)
		}
	}

	for _, author := range p.Author {
		ip.Contributor = append(ip.Contributor, author.Name())
		if author.PersonID != "" {
			ip.AuthorID = append(ip.AuthorID, author.PersonID)
		}
	}
	ip.AuthorID = util.UniqStrings(ip.AuthorID)

	for _, supervisor := range p.Supervisor {
		ip.Contributor = append(ip.Contributor, supervisor.Name())
	}

	for _, editor := range p.Editor {
		ip.Contributor = append(ip.Contributor, editor.Name())
	}

	ip.Contributor = util.UniqStrings(ip.Contributor)

	for _, file := range p.File {
		ip.FileRelation = append(ip.FileRelation, file.Relation)
	}
	ip.FileRelation = util.UniqStrings(ip.FileRelation)

	if p.DOI != "" {
		ip.Identifier = append(ip.Identifier, p.DOI)
	}
	ip.Identifier = append(ip.Identifier, p.ISBN...)
	ip.Identifier = append(ip.Identifier, p.EISBN...)
	ip.Identifier = append(ip.Identifier, p.ISSN...)
	ip.Identifier = append(ip.Identifier, p.EISSN...)
	if p.WOSID != "" {
		ip.Identifier = append(ip.Identifier, p.WOSID)
	}
	if p.ArxivID != "" {
		ip.Identifier = append(ip.Identifier, p.ArxivID)
	}
	if p.PubMedID != "" {
		ip.Identifier = append(ip.Identifier, p.PubMedID)
	}
	if p.VABBID != "" {
		ip.Identifier = append(ip.Identifier, p.VABBID)
	}
	if p.SourceID != "" {
		ip.Identifier = append(ip.Identifier, p.SourceID)
	}

	// issn/isbn may have dashes or not, so separate analyzing is necessary
	ip.ISXN = append(ip.ISXN, p.ISBN...)
	ip.ISXN = append(ip.ISXN, p.EISBN...)
	ip.ISXN = append(ip.ISXN, p.ISSN...)
	ip.ISXN = append(ip.ISXN, p.EISSN...)

	return ip
}
