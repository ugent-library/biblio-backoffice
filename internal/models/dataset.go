package models

import (
	"fmt"
	"strings"
	"time"
)

type DatasetHits struct {
	Pagination
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

type RelatedPublication struct {
	ID string `json:"id,omitempty"`
}

type Dataset struct {
	Abstract          []Text         `json:"abstract,omitempty" form:"abstract"`
	AccessLevel       string         `json:"access_level,omitempty" form:"access_level"`
	Author            []*Contributor `json:"author,omitempty" form:"-"` // TODO rename to Creator
	CompletenessScore int            `json:"completeness_score,omitempty" form:"-"`
	Contributor       []*Contributor `json:"contributor,omitempty" form:"-"`
	// CreationContext         string              `json:"creation_context,omitempty" form:"-"`
	CreatorID          string               `json:"creator_id,omitempty" form:"-"`
	DateCreated        *time.Time           `json:"date_created,omitempty" form:"-"`
	DateUpdated        *time.Time           `json:"date_updated,omitempty" form:"-"`
	Department         []DatasetDepartment  `json:"department,omitempty" form:"-"`
	DOI                string               `json:"doi,omitempty" form:"-"`
	Embargo            string               `json:"embargo,omitempty" form:"embargo"`
	EmbargoTo          string               `json:"embargo_to,omitempty" form:"embargo_to"`
	Format             []string             `json:"format,omitempty" form:"format"`
	ID                 string               `json:"id,omitempty" form:"-"`
	Keyword            []string             `json:"keyword,omitempty" form:"keyword"`
	License            string               `json:"license,omitempty" form:"license"`
	Locked             bool                 `json:"locked,omitempty" form:"-"`
	Message            string               `json:"message,omitempty" form:"-"`
	OtherLicense       string               `json:"other_license,omitempty" form:"other_license"`
	Project            []DatasetProject     `json:"project,omitempty" form:"-"`
	Publisher          string               `json:"publisher,omitempty" form:"publisher"`
	RelatedPublication []RelatedPublication `json:"related_publication,omitempty" form:"-"`
	ReviewerNote       string               `json:"reviewer_note,omitempty" form:"-"`
	ReviewerTags       []string             `json:"reviewer_tags,omitempty" form:"-"`
	Status             string               `json:"status,omitempty" form:"-"`
	Title              string               `json:"title,omitempty" form:"title"`
	URL                string               `json:"url,omitempty" form:"url"`
	UserID             string               `json:"user_id,omitempty" form:"-"`
	// Version                 int                 `json:"_version,omitempty" form:"-"`
	Year string `json:"year,omitempty" form:"year"`
}

func (d *Dataset) HasRelatedPublication(id string) bool {
	for _, r := range d.RelatedPublication {
		if r.ID == id {
			return true
		}
	}
	return false
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

func (d *Dataset) Validate() (errs ValidationErrors) {
	if d.ID == "" {
		errs = append(errs, ValidationError{
			Pointer: "/id",
			Code:    "required",
		})
	}
	if d.Status == "" {
		errs = append(errs, ValidationError{
			Pointer: "/status",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.AccessLevel == "" {
		errs = append(errs, ValidationError{
			Pointer: "/access_level",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.DOI == "" {
		errs = append(errs, ValidationError{
			Pointer: "/doi",
			Code:    "required",
		})
	}
	if d.Status == "public" && len(d.Format) == 0 {
		errs = append(errs, ValidationError{
			Pointer: "/format",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.Publisher == "" {
		errs = append(errs, ValidationError{
			Pointer: "/publisher",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.Title == "" {
		errs = append(errs, ValidationError{
			Pointer: "/title",
			Code:    "required",
		})
	}
	if d.Status == "public" && d.Year == "" {
		errs = append(errs, ValidationError{
			Pointer: "/year",
			Code:    "required",
		})
	}
	if d.Status == "public" && len(d.Author) == 0 {
		errs = append(errs, ValidationError{
			Pointer: "/author",
			Code:    "required",
		})
	}

	for i, c := range d.Author {
		for _, err := range c.Validate() {
			errs = append(errs, ValidationError{
				Pointer: fmt.Sprintf("/author/%d/%s", i, err.Pointer),
				Code:    err.Code,
			})
		}
	}

	return
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
