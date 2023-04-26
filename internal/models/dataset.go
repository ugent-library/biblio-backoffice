package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/internal/pagination"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

type DatasetHits struct {
	pagination.Pagination
	Hits   []*Dataset             `json:"hits"`
	Facets map[string]FacetValues `json:"facets"`
}

type DatasetUser struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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
	Abstract    []Text         `json:"abstract,omitempty"`
	AccessLevel string         `json:"access_level,omitempty"`
	Author      []*Contributor `json:"author,omitempty"` // TODO rename to Creator
	// CompletenessScore  int                  `json:"completeness_score,omitempty"`
	BatchID                 string               `json:"batch_id,omitempty"`
	Contributor             []*Contributor       `json:"contributor,omitempty"`
	Creator                 *DatasetUser         `json:"creator,omitempty"`
	DateCreated             *time.Time           `json:"date_created,omitempty"`
	DateUpdated             *time.Time           `json:"date_updated,omitempty"`
	DateFrom                *time.Time           `json:"date_from,omitempty"`
	DateUntil               *time.Time           `json:"date_until,omitempty"`
	Department              []DatasetDepartment  `json:"department,omitempty"`
	DOI                     string               `json:"doi,omitempty"`
	EmbargoDate             string               `json:"embargo_date,omitempty"`
	AccessLevelAfterEmbargo string               `json:"access_level_after_embargo,omitempty"`
	Format                  []string             `json:"format,omitempty"`
	Handle                  string               `json:"handle,omitempty"`
	ID                      string               `json:"id,omitempty"`
	Keyword                 []string             `json:"keyword,omitempty"`
	HasBeenPublic           bool                 `json:"has_been_public"`
	Language                []string             `json:"language,omitempty"`
	LastUser                *DatasetUser         `json:"last_user,omitempty"`
	License                 string               `json:"license,omitempty"`
	Locked                  bool                 `json:"locked"`
	Message                 string               `json:"message,omitempty"`
	OtherLicense            string               `json:"other_license,omitempty"`
	Project                 []DatasetProject     `json:"project,omitempty"`
	Publisher               string               `json:"publisher,omitempty"`
	RelatedPublication      []RelatedPublication `json:"related_publication,omitempty"`
	ReviewerNote            string               `json:"reviewer_note,omitempty"`
	ReviewerTags            []string             `json:"reviewer_tags,omitempty"`
	SnapshotID              string               `json:"snapshot_id,omitempty"`
	Status                  string               `json:"status,omitempty"`
	Title                   string               `json:"title,omitempty"`
	URL                     string               `json:"url,omitempty"`
	User                    *DatasetUser         `json:"user,omitempty"`
	Year                    string               `json:"year,omitempty"`
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

func (d *Dataset) GetAbstract(id string) *Text {
	for _, abstract := range d.Abstract {
		if abstract.ID == id {
			return &abstract
		}
	}
	return nil
}

func (d *Dataset) SetAbstract(t *Text) {
	for i, abstract := range d.Abstract {
		if abstract.ID == t.ID {
			d.Abstract[i] = *t
		}
	}
}

func (d *Dataset) AddAbstract(t *Text) {
	t.ID = ulid.Make().String()
	d.Abstract = append(d.Abstract, *t)
}

func (d *Dataset) RemoveAbstract(id string) {
	abstracts := make([]Text, 0)
	for _, abstract := range d.Abstract {
		if abstract.ID != id {
			abstracts = append(abstracts, abstract)
		}
	}
	d.Abstract = abstracts
}

func (d *Dataset) GetProject(id string) *DatasetProject {
	for _, project := range d.Project {
		if project.ID == id {
			return &project
		}
	}
	return nil
}

func (d *Dataset) RemoveProject(id string) {
	projects := make([]DatasetProject, 0)
	for _, project := range d.Project {
		if project.ID != id {
			projects = append(projects, project)
		}
	}
	d.Project = projects
}

func (d *Dataset) AddProject(pr *DatasetProject) {
	d.RemoveProject(pr.ID)
	d.Project = append(d.Project, *pr)
}

func (d *Dataset) GetDepartment(id string) *DatasetDepartment {
	for _, department := range d.Department {
		if department.ID == id {
			return &department
		}
	}
	return nil
}

func (d *Dataset) RemoveDepartment(id string) {
	departments := make([]DatasetDepartment, 0)
	for _, department := range d.Department {
		if department.ID != id {
			departments = append(departments, department)
		}
	}
	d.Department = departments
}

func (d *Dataset) AddDepartmentByOrg(org *Organization) {
	// remove if added before
	d.RemoveDepartment(org.ID)

	datasetDepartment := DatasetDepartment{ID: org.ID}
	for _, d := range org.Tree {
		datasetDepartment.Tree = append(datasetDepartment.Tree, DatasetDepartmentRef(d))
	}
	d.Department = append(d.Department, datasetDepartment)
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

	if d.Status == "public" && d.AccessLevel == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/access_level",
			Code:    "dataset.access_level.required",
		})
	}
	if d.AccessLevel != "" && !validation.IsDatasetAccessLevel(d.AccessLevel) {
		errs = append(errs, &validation.Error{
			Pointer: "/access_level",
			Code:    "dataset.access_level.invalid",
		})
	}

	if d.Status == "public" && d.DOI == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/doi",
			Code:    "dataset.doi.required",
		})
	}

	if d.Status == "public" && len(d.Format) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/format",
			Code:    "dataset.format.required",
		})
	}
	for i, f := range d.Format {
		if f == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/format/%d", i),
				Code:    "dataset.format.invalid",
			})
		}
	}

	// for i, k := range d.Keyword {
	// 	if k == "" {
	// 		errs = append(errs, &validation.Error{
	// 			Pointer: fmt.Sprintf("/keyword/%d", i),
	// 			Code:    "dataset.keyword.invalid",
	// 		})
	// 	}
	// }

	for i, l := range d.Language {
		if !validation.InArray(vocabularies.Map["language_codes"], l) {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/language/%d", i),
				Code:    "dataset.language.invalid",
			})
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

	if d.Status == "public" && d.Year == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/year",
			Code:    "dataset.year.required",
		})
	}
	if d.Year != "" && !validation.IsYear(d.Year) {
		errs = append(errs, &validation.Error{
			Pointer: "/year",
			Code:    "dataset.year.invalid",
		})
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
				Code:    "dataset.author." + err.Code,
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

	// TODO IsDate and co. are only checked when dataset is public
	if d.Status == "public" && d.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
		if d.EmbargoDate == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_date",
				Code:    "dataset.embargo_date.required",
			})
		} else if !validation.IsDate(d.EmbargoDate) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_date",
				Code:    "dataset.embargo_date.invalid",
			})
		}

		invalid := false
		if d.AccessLevelAfterEmbargo == "" {
			invalid = true
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_after_embargo",
				Code:    "dataset.access_level_after_embargo.required",
			})
		}

		if d.AccessLevel == d.AccessLevelAfterEmbargo && !invalid {
			invalid = true
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_after_embargo",
				Code:    "dataset.access_level_after_embargo.similar", // TODO better code
			})
		}

		if !validation.IsDatasetAccessLevel(d.AccessLevelAfterEmbargo) && !invalid {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_after_embargo",
				Code:    "dataset.access_level_after_embargo.invalid", // TODO better code
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

func (d *Dataset) ClearEmbargo() {
	d.AccessLevel = d.AccessLevelAfterEmbargo
	d.AccessLevelAfterEmbargo = ""
	d.EmbargoDate = ""
}
