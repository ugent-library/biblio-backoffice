package models

import (
	"strings"
	"time"
)

type DatasetHits struct {
	Pagination
	// Total        int                `json:"total"`
	// Page         int                `json:"page"`
	// LastPage     int                `json:"last_page"`
	// PreviousPage bool               `json:"previous_page"`
	// NextPage     bool               `json:"next_page"`
	// FirstOnPage  int                `json:"first_on_page"`
	// LastOnPage   int                `json:"last_on_page"`
	Hits   []*Dataset         `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type DatasetDepartment struct {
	ID string `json:"id,omitempty"`
}

type DatasetProject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Dataset struct {
	Abstract          []Text         `json:"abstract,omitempty" form:"abstract"`
	AccessLevel       string         `json:"access_level,omitempty" form:"access_level"`
	Author            []*Contributor `json:"author,omitempty" form:"-"` // TODO rename to Creator
	CompletenessScore int            `json:"completeness_score,omitempty" form:"-"`
	Contributor       []*Contributor `json:"contributor,omitempty" form:"-"`
	// CreationContext         string              `json:"creation_context,omitempty" form:"-"`
	CreatorID               string              `json:"creator_id,omitempty" form:"-"`
	DateCreated             *time.Time          `json:"date_created,omitempty" form:"-"`
	DateUpdated             *time.Time          `json:"date_updated,omitempty" form:"-"`
	Department              []DatasetDepartment `json:"department,omitempty" form:"-"`
	DOI                     string              `json:"doi,omitempty" form:"-"`
	Embargo                 string              `json:"embargo,omitempty" form:"embargo"`
	EmbargoTo               string              `json:"embargo_to,omitempty" form:"embargo_to"`
	Format                  []string            `json:"format,omitempty" form:"format"`
	ID                      string              `json:"id,omitempty" form:"-"`
	Keyword                 []string            `json:"keyword,omitempty" form:"keyword"`
	License                 string              `json:"license,omitempty" form:"license"`
	Locked                  bool                `json:"locked,omitempty" form:"-"`
	Message                 string              `json:"message,omitempty" form:"-"`
	OtherLicense            string              `json:"other_license,omitempty" form:"other_license"`
	Project                 []DatasetProject    `json:"project,omitempty" form:"-"`
	Publisher               string              `json:"publisher,omitempty" form:"publisher"`
	RelatedPublicationCount int                 `json:"related_publication_count" form:"-"`
	ReviewerNote            string              `json:"reviewer_note,omitempty" form:"-"`
	ReviewerTags            []string            `json:"reviewer_tags,omitempty" form:"-"`
	Status                  string              `json:"status,omitempty" form:"-"`
	Title                   string              `json:"title,omitempty" form:"title"`
	URL                     string              `json:"url,omitempty" form:"url"`
	UserID                  string              `json:"user_id,omitempty" form:"-"`
	// Version                 int                 `json:"_version,omitempty" form:"-"`
	Year string `json:"year,omitempty" form:"year"`
}

func (d *Dataset) Contributors(role string) []*Contributor {
	switch role {
	case "author":
		return d.Author
	case "contributor":
		return d.Contributor
	default:
		return nil
	}
}

func (p *Dataset) SetContributors(role string, c []*Contributor) {
	switch role {
	case "author":
		p.Author = c
	case "contributor":
		p.Contributor = c
	}
}

func (p *Dataset) AddContributor(role string, i int, c *Contributor) {
	cc := p.Contributors(role)

	if len(cc) == i {
		p.SetContributors(role, append(cc, c))
		return
	}

	newCC := append(cc[:i+1], cc[i:]...)
	newCC[i] = c
	p.SetContributors(role, newCC)
}

func (p *Dataset) RemoveContributor(role string, i int) {
	cc := p.Contributors(role)

	p.SetContributors(role, append(cc[:i], cc[i+1:]...))
}

func (d *Dataset) ResolveDOI() string {
	if d.DOI != "" {
		return "https://doi.org/" + d.DOI

	}
	return ""
}

func (d *Dataset) Vacuum() {
	d.AccessLevel = strings.TrimSpace(d.AccessLevel)
	d.Embargo = strings.TrimSpace(d.Embargo)
	d.EmbargoTo = strings.TrimSpace(d.EmbargoTo)
	d.Format = vacuumStringSlice(d.Format)
	d.Keyword = vacuumStringSlice(d.Keyword)
	d.License = strings.TrimSpace(d.License)
	d.OtherLicense = strings.TrimSpace(d.OtherLicense)
	d.Publisher = strings.TrimSpace(d.Publisher)
	d.Title = strings.TrimSpace(d.Title)
	d.URL = strings.TrimSpace(d.URL)
	d.Year = strings.TrimSpace(d.Year)
}

func vacuumStringSlice(vals []string) []string {
	newVals := []string{}
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			newVals = append(newVals, v)
		}
	}
	return newVals
}
