package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/pagination"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type DatasetHits struct {
	pagination.Pagination
	Hits   []*Dataset         `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type DatasetDepartmentRef struct {
	ID string `json:"id,omitempty"`
}

type DatasetDepartment struct {
	ID   string                 `json:"id,omitempty"`
	Tree []DatasetDepartmentRef `json:"tree,omitempty"`
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
	DateFrom           *time.Time           `json:"date_from,omitempty" form:"-"`
	DateUntil          *time.Time           `json:"date_until,omitempty" form:"-"`
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

func (d *Dataset) GetAbstract(i int) (Text, error) {
	if i >= len(d.Abstract) {
		return Text{}, errors.New("index out of bounds")
	}

	return d.Abstract[i], nil
}

func (d *Dataset) SetAbstract(i int, t Text) error {
	if i >= len(d.Abstract) {
		return errors.New("index out of bounds")
	}

	d.Abstract[i] = t

	return nil
}

func (d *Dataset) RemoveAbstract(i int) error {
	if i >= len(d.Abstract) {
		return errors.New("index out of bounds")
	}

	d.Abstract = append(d.Abstract[:i], d.Abstract[i+1:]...)

	return nil
}

func (d *Dataset) GetProject(i int) (DatasetProject, error) {
	if i >= len(d.Project) {
		return DatasetProject{}, errors.New("index out of bounds")
	}

	return d.Project[i], nil
}

func (d *Dataset) RemoveProject(i int) error {
	if i >= len(d.Project) {
		return errors.New("index out of bounds")
	}

	d.Project = append(d.Project[:i], d.Project[i+1:]...)

	return nil
}

func (d *Dataset) GetDepartment(i int) (DatasetDepartment, error) {
	if i >= len(d.Department) {
		return DatasetDepartment{}, errors.New("index out of bounds")
	}

	return d.Department[i], nil
}

func (d *Dataset) RemoveDepartment(i int) error {
	if i >= len(d.Department) {
		return errors.New("index out of bounds")
	}

	d.Department = append(d.Department[:i], d.Department[i+1:]...)

	return nil
}

func (d *Dataset) ResolveDOI() string {
	if d.DOI != "" {
		return "https://doi.org/" + d.DOI

	}
	return ""
}

func (d *Dataset) Validate() error {
	var errs validation.Errors

	// if d.ID == "" {
	// 	errs = append(errs, &validation.Error{
	// 		Pointer: "/id",
	// 		Code:    "required",
	// 		Field:   "id",
	// 	})
	// }
	if d.Status == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "required",
			Field:   "status",
		})
	} else if !validation.IsStatus(d.Status) {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "invalid",
			Field:   "status",
		})
	}
	if d.Status == "public" {
		if d.AccessLevel == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level",
				Code:    "required",
				Field:   "access_level",
			})
		} else if !validation.IsDatasetAccessLevel(d.AccessLevel) {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level",
				Code:    "invalid",
				Field:   "access_level",
			})
		}
	}
	if d.Status == "public" && d.DOI == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/doi",
			Code:    "required",
			Field:   "doi",
		})
	}
	if d.Status == "public" {
		if len(d.Format) == 0 {
			errs = append(errs, &validation.Error{
				Pointer: "/format",
				Code:    "required",
				Field:   "format",
			})
		}
		for i, f := range d.Format {
			if f == "" {
				errs = append(errs, &validation.Error{
					Pointer: fmt.Sprintf("/format/%d", i),
					Code:    "required",
					Field:   "format",
				})
			}
		}
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

	if d.Status == "public" {
		// year =~ /^\d{4}$/
		if d.Year == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "required",
				Field:   "year",
			})
		} else if !validation.IsYear(d.Year) {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "invalid",
				Field:   "year",
			})
		}
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

	// at least one ugent author
	if d.Status == "public" {
		var hasUgentAuthors bool = false
		for _, a := range d.Author {
			if a.ID != "" {
				hasUgentAuthors = true
				break
			}
		}
		if !hasUgentAuthors {
			errs = append(errs, &validation.Error{
				Pointer: "/author",
				Code:    "min_ugent_authors",
				Field:   "author",
			})
		}
	}

	// license or other_license -> TODO: base error?
	// now "fixed" by (incorrectly) pointing at license
	if d.Status == "public" && d.License == "" && d.OtherLicense == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/license",
			Code:    "required",
			Field:   "license",
		})
	}

	for i, rp := range d.RelatedPublication {
		if rp.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/related_publication/%d/id", i),
				Code:    "required",
				Field:   "related_publication",
			})
		}
	}

	for i, pr := range d.Project {
		if pr.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/project/%d/id", i),
				Code:    "required",
				Field:   "project",
			})
		}
	}

	for i, dep := range d.Department {
		if dep.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/department/%d/id", i),
				Code:    "required",
				Field:   "department",
			})
		}
	}

	if d.Status == "public" && d.AccessLevel == vocabularies.EmbargoedAccess {
		if d.Embargo == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo",
				Code:    "required",
				Field:   "embargo",
			})
		} else if !validation.IsDate(d.Embargo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo",
				Code:    "invalid",
				Field:   "embargo",
			})
		}
		if d.EmbargoTo == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "required",
				Field:   "embargo_to",
			})
		} else if d.AccessLevel == d.EmbargoTo {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "invalid", //TODO: better code
				Field:   "embargo_to",
			})
		} else if !validation.IsDatasetAccessLevel(d.EmbargoTo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "invalid", //TODO: better code
				Field:   "embargo_to",
			})
		}
	}

	for i, abstract := range d.Abstract {
		for _, err := range abstract.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/abstract/%d%s", i, err.Pointer),
				Code:    "dataset.abstract." + err.Code,
			})
		}
	}

	// TODO: why is the nil slice validationErrors(nil) != nil in mutantdb validation?
	if len(errs) > 0 {
		return errs
	}
	return nil
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
