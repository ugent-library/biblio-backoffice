package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/pagination"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type DatasetHits struct {
	pagination.Pagination
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
	Abstract           []Text               `json:"abstract,omitempty" form:"abstract"`
	AccessLevel        string               `json:"access_level,omitempty" form:"access_level"`
	Author             []*Contributor       `json:"author,omitempty" form:"-"` // TODO rename to Creator
	CompletenessScore  int                  `json:"completeness_score,omitempty" form:"-"`
	Contributor        []*Contributor       `json:"contributor,omitempty" form:"-"`
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
	SnapshotID         string               `json:"-" form:"-"`
	Status             string               `json:"status,omitempty" form:"-"`
	Title              string               `json:"title,omitempty" form:"title"`
	URL                string               `json:"url,omitempty" form:"url"`
	UserID             string               `json:"user_id,omitempty" form:"-"`
	Year               string               `json:"year,omitempty" form:"year"`
}

func (d *Dataset) Clone() *Dataset {
	clone := *d
	clone.Abstract = nil
	clone.Abstract = append(clone.Abstract, d.Abstract...)
	clone.Author = nil
	for _, c := range d.Author {
		clone.Author = append(clone.Author, c.Clone())
	}
	clone.Contributor = nil
	for _, c := range d.Contributor {
		clone.Contributor = append(clone.Contributor, c.Clone())
	}
	clone.Department = nil
	clone.Department = append(clone.Department, d.Department...)
	clone.Format = nil
	clone.Format = append(clone.Format, d.Format...)
	clone.Keyword = nil
	clone.Keyword = append(clone.Keyword, d.Keyword...)
	clone.Project = nil
	clone.Project = append(clone.Project, d.Project...)
	clone.RelatedPublication = nil
	clone.RelatedPublication = append(clone.RelatedPublication, d.RelatedPublication...)
	clone.ReviewerTags = nil
	clone.ReviewerTags = append(clone.ReviewerTags, d.ReviewerTags...)
	return &clone
}

func (d *Dataset) HasRelatedPublication(id string) bool {
	for _, r := range d.RelatedPublication {
		if r.ID == id {
			return true
		}
	}
	return false
}

func (d *Dataset) RemoveRelatedPublication(id string) {
	var publications []RelatedPublication
	for _, r := range d.RelatedPublication {
		if r.ID != id {
			publications = append(publications, r)
		}
	}
	d.RelatedPublication = publications
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

func (d *Dataset) Validate() (errs validation.Errors) {
	if d.ID == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/id",
			Code:    "required",
			Field:   "id",
		})
	}
	if d.Status == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "required",
			Field:   "status",
		})
	}
	if d.Status == "public" && d.AccessLevel == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/access_level",
			Code:    "required",
			Field:   "access_level",
		})
	}
	if d.Status == "public" && d.DOI == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/doi",
			Code:    "required",
			Field:   "doi",
		})
	}
	if d.Status == "public" && len(d.Format) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/format",
			Code:    "required",
			Field:   "format",
		})
	}
	if d.Status == "public" && d.Publisher == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/publisher",
			Code:    "required",
			Field:   "publisher",
		})
	}
	if d.Status == "public" && d.Title == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/title",
			Code:    "required",
			Field:   "title",
		})
	}
	if d.Status == "public" && d.Year == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/year",
			Code:    "required",
			Field:   "year",
		})
	}
	if d.Status == "public" && len(d.Author) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/author",
			Code:    "required",
			Field:   "author",
		})
	}

	for i, c := range d.Author {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/author/%d%s", i, err.Pointer),
				Code:    err.Code,
				Field:   "author." + err.Field,
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
