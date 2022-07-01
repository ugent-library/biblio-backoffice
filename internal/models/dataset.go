package models

import (
	"errors"
	"fmt"
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
	Abstract           []Text               `json:"abstract,omitempty"`
	AccessLevel        string               `json:"access_level,omitempty"`
	Author             []*Contributor       `json:"author,omitempty"` // TODO rename to Creator
	CompletenessScore  int                  `json:"completeness_score,omitempty"`
	Contributor        []*Contributor       `json:"contributor,omitempty"`
	CreatorID          string               `json:"creator_id,omitempty"`
	DateCreated        *time.Time           `json:"date_created,omitempty"`
	DateUpdated        *time.Time           `json:"date_updated,omitempty"`
	DateFrom           *time.Time           `json:"date_from,omitempty"`
	DateUntil          *time.Time           `json:"date_until,omitempty"`
	Department         []DatasetDepartment  `json:"department,omitempty"`
	DOI                string               `json:"doi,omitempty"`
	Embargo            string               `json:"embargo,omitempty"`
	EmbargoTo          string               `json:"embargo_to,omitempty"`
	Format             []string             `json:"format,omitempty"`
	ID                 string               `json:"id,omitempty"`
	Keyword            []string             `json:"keyword,omitempty"`
	License            string               `json:"license,omitempty"`
	Locked             bool                 `json:"locked,omitempty"`
	Message            string               `json:"message,omitempty"`
	OtherLicense       string               `json:"other_license,omitempty"`
	Project            []DatasetProject     `json:"project,omitempty"`
	Publisher          string               `json:"publisher,omitempty"`
	RelatedPublication []RelatedPublication `json:"related_publication,omitempty"`
	ReviewerNote       string               `json:"reviewer_note,omitempty"`
	ReviewerTags       []string             `json:"reviewer_tags,omitempty"`
	SnapshotID         string               `json:"-"`
	Status             string               `json:"status,omitempty"`
	Title              string               `json:"title,omitempty"`
	URL                string               `json:"url,omitempty"`
	UserID             string               `json:"user_id,omitempty"`
	Year               string               `json:"year,omitempty"`
}

// TODO still needed?
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

func (d *Dataset) SetContributors(role string, c []*Contributor) {
	switch role {
	case "author":
		d.Author = c
	case "contributor":
		d.Contributor = c
	}
}

func (d *Dataset) GetContributor(role string, i int) (*Contributor, error) {
	cc := d.Contributors(role)
	if i >= len(cc) {
		return nil, errors.New("index out of bounds")
	}

	return cc[i], nil
}

func (d *Dataset) AddContributor(role string, c *Contributor) {
	d.SetContributors(role, append(d.Contributors(role), c))
}

func (d *Dataset) SetContributor(role string, i int, c *Contributor) error {
	cc := d.Contributors(role)
	if i >= len(cc) {
		return errors.New("index out of bounds")
	}

	cc[i] = c

	return nil
}

func (d *Dataset) RemoveContributor(role string, i int) error {
	cc := d.Contributors(role)
	if i >= len(cc) {
		return errors.New("index out of bounds")
	}

	d.SetContributors(role, append(cc[:i], cc[i+1:]...))

	return nil
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

	if d.ID == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/id",
			Code:    "dataset.id.required",
		})
	}
	if d.Status == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "dataset.status.required",
		})
	} else if !validation.IsStatus(d.Status) {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "dataset.status.invalid",
		})
	}
	if d.Status == "public" {
		if d.AccessLevel == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level",
				Code:    "dataset.access_level.required",
			})
		} else if !validation.IsDatasetAccessLevel(d.AccessLevel) {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level",
				Code:    "dataset.access_level.invalid",
			})
		}
	}
	if d.Status == "public" && d.DOI == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/doi",
			Code:    "dataset.doi.required",
		})
	}
	if d.Status == "public" {
		if len(d.Format) == 0 {
			errs = append(errs, &validation.Error{
				Pointer: "/format",
				Code:    "dataset.format.required",
			})
		}
		for i, f := range d.Format {
			if f == "" {
				errs = append(errs, &validation.Error{
					Pointer: fmt.Sprintf("/format/%d", i),
					Code:    "dataset.format.required",
				})
			}
		}
	}
	if d.Status == "public" && d.Publisher == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/publisher",
			Code:    "dataset.publisher.required",
		})
	}
	if d.Status == "public" && d.Title == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/title",
			Code:    "dataset.title.required",
		})
	}

	if d.Status == "public" {
		// year =~ /^\d{4}$/
		if d.Year == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "dataset.year.required",
			})
		} else if !validation.IsYear(d.Year) {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "dataset.year.invalid",
			})
		}
	}
	if d.Status == "public" && len(d.Author) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/author",
			Code:    "dataset.author.required",
		})
	}

	for i, c := range d.Author {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/author/%d%s", i, err.Pointer),
				Code:    "dataset.author" + err.Code,
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
				Code:    "dataset.author.min_ugent_authors",
			})
		}
	}

	// license or other_license -> TODO: base error?
	// now "fixed" by (incorrectly) pointing at license
	if d.Status == "public" && d.License == "" && d.OtherLicense == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/license",
			Code:    "dataset.license.required",
		})
	}

	for i, rp := range d.RelatedPublication {
		if rp.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/related_publication/%d/id", i),
				Code:    "dataset.related_publication.required",
			})
		}
	}

	for i, pr := range d.Project {
		if pr.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/project/%d/id", i),
				Code:    "dataset.project.id.required",
			})
		}
	}

	for i, dep := range d.Department {
		if dep.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/department/%d/id", i),
				Code:    "dataset.department.id.required",
			})
		}
	}

	if d.Status == "public" && d.AccessLevel == vocabularies.EmbargoedAccess {
		if d.Embargo == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo",
				Code:    "dataset.embargo.required",
			})
		} else if !validation.IsDate(d.Embargo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo",
				Code:    "dataset.embargo.invalid",
			})
		}
		if d.EmbargoTo == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "dataset.embargo_to.required",
			})
		} else if d.AccessLevel == d.EmbargoTo {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "dataset.embarg_to.invalid", //TODO: better code
			})
		} else if !validation.IsDatasetAccessLevel(d.EmbargoTo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "dataset.embargo_to.invalid", //TODO: better code
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
