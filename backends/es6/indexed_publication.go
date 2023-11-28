package es6

import (
	"regexp"

	"slices"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	internal_time "github.com/ugent-library/biblio-backoffice/time"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type indexedPublication struct {
	AuthorID                []string `json:"author_id,omitempty"`
	AlternativeTitle        []string `json:"alternative_title,omitempty"`
	BatchID                 string   `json:"batch_id,omitempty"`
	Classification          string   `json:"classification,omitempty"`
	Contributor             []string `json:"contributor,omitempty"`
	CreatorID               string   `json:"creator_id,omitempty"`
	ConferenceName          string   `json:"conference_name,omitempty"`
	DateCreated             string   `json:"date_created"`
	DateUpdated             string   `json:"date_updated"`
	DOI                     string   `json:"doi,omitempty"`
	Extern                  bool     `json:"extern"`
	FacultyID               []string `json:"faculty_id,omitempty"`
	FileRelation            []string `json:"file_relation,omitempty"`
	HasMessage              bool     `json:"has_message"`
	HasFiles                bool     `json:"has_files"`
	ID                      string   `json:"id,omitempty"`
	Identifier              []string `json:"identifier,omitempty"`
	IssueTitle              string   `json:"issue_title,omitempty"`
	ISXN                    []string `json:"isxn,omitempty"`
	Keyword                 []string `json:"keyword,omitempty"`
	LastUserID              string   `json:"last_user_id,omitempty"`
	Legacy                  bool     `json:"legacy"`
	Locked                  bool     `json:"locked"`
	OrganizationID          []string `json:"organization_id,omitempty"`
	Publication             string   `json:"publication,omitempty"`
	PublicationAbbreviation string   `json:"publication_abbreviation,omitempty"`
	PublicationStatus       string   `json:"publication_status,omitempty"`
	Publisher               string   `json:"publisher,omitempty"`
	ReviewerTags            []string `json:"reviewer_tags,omitempty"`
	SeriesTitle             string   `json:"series_title,omitempty"`
	Status                  string   `json:"status,omitempty"`
	Title                   string   `json:"title,omitempty"`
	Type                    string   `json:"type,omitempty"`
	UserID                  string   `json:"user_id,omitempty"`
	VABBType                string   `json:"vabb_type,omitempty"`
	WOSType                 []string `json:"wos_type,omitempty"`
	Year                    string   `json:"year,omitempty"`
}

var reSplitWOS *regexp.Regexp = regexp.MustCompile(`\s*[,;]\s*`)

func NewIndexedPublication(p *models.Publication) *indexedPublication {
	ip := &indexedPublication{
		BatchID:                 p.BatchID,
		Classification:          p.Classification,
		CreatorID:               p.CreatorID,
		ConferenceName:          p.ConferenceName,
		DateCreated:             internal_time.FormatTimeUTC(p.DateCreated),
		DateUpdated:             internal_time.FormatTimeUTC(p.DateUpdated),
		DOI:                     p.DOI,
		Extern:                  p.Extern,
		ID:                      p.ID,
		IssueTitle:              p.IssueTitle,
		LastUserID:              p.LastUserID,
		Legacy:                  p.Legacy,
		Locked:                  p.Locked,
		HasMessage:              len(p.Message) > 0,
		Publication:             p.Publication,
		PublicationAbbreviation: p.PublicationAbbreviation,
		PublicationStatus:       p.PublicationStatus,
		Publisher:               p.Publisher,
		ReviewerTags:            p.ReviewerTags,
		SeriesTitle:             p.SeriesTitle,
		Status:                  p.Status,
		Title:                   p.Title,
		Type:                    p.Type,
		UserID:                  p.UserID,
		VABBType:                p.VABBType,
		Year:                    p.Year,
		HasFiles:                len(p.File) > 0,
		Keyword:                 p.Keyword,
		AlternativeTitle:        p.AlternativeTitle,
	}

	faculties := vocabularies.Map["faculties"]

	// extract faculty_id and all organization id's from organization trees
	for _, rel := range p.RelatedOrganizations {
		for _, org := range rel.Organization.Tree {
			if !slices.Contains(ip.OrganizationID, org.ID) {
				ip.OrganizationID = append(ip.OrganizationID, org.ID)
			}

			if slices.Contains(faculties, org.ID) && !slices.Contains(ip.FacultyID, org.ID) {
				ip.FacultyID = append(ip.FacultyID, org.ID)
			}
		}
	}

	if len(ip.FacultyID) == 0 {
		ip.FacultyID = append(ip.FacultyID, backends.MissingValue)
	}

	if ip.PublicationStatus == "" {
		ip.PublicationStatus = backends.MissingValue
	}

	if p.WOSType != "" {
		ip.WOSType = reSplitWOS.Split(p.WOSType, -1)
	}

	for _, author := range p.Author {
		ip.Contributor = append(ip.Contributor, author.Name())
		if author.PersonID != "" {
			ip.AuthorID = append(ip.AuthorID, author.PersonID)
		}
	}
	ip.AuthorID = lo.Uniq(ip.AuthorID)

	for _, supervisor := range p.Supervisor {
		ip.Contributor = append(ip.Contributor, supervisor.Name())
	}

	for _, editor := range p.Editor {
		ip.Contributor = append(ip.Contributor, editor.Name())
	}

	ip.Contributor = lo.Uniq(ip.Contributor)

	for _, file := range p.File {
		ip.FileRelation = append(ip.FileRelation, file.Relation)
	}
	ip.FileRelation = lo.Uniq(ip.FileRelation)

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
