package models

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/pagination"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
)

var reYear = regexp.MustCompile("^[0-9]{4}$")

type DatasetHits struct {
	pagination.Pagination
	Hits   []*Dataset             `json:"hits"`
	Facets map[string]FacetValues `json:"facets"`
}

type DatasetLink struct {
	ID          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Description string `json:"description,omitempty"`
}

type RelatedPublication struct {
	ID string `json:"id,omitempty"`
}

type Dataset struct {
	Abstract    []*Text        `json:"abstract,omitempty"`
	AccessLevel string         `json:"access_level,omitempty"`
	Author      []*Contributor `json:"author,omitempty"` // TODO rename to Creator
	// CompletenessScore  int                  `json:"completeness_score,omitempty"`
	BatchID                 string                 `json:"batch_id,omitempty"`
	Contributor             []*Contributor         `json:"contributor,omitempty"`
	CreatorID               string                 `json:"creator_id,omitempty"`
	Creator                 *Person                `json:"-"`
	DateCreated             *time.Time             `json:"date_created,omitempty"`
	DateUpdated             *time.Time             `json:"date_updated,omitempty"`
	DateFrom                *time.Time             `json:"date_from,omitempty"`
	DateUntil               *time.Time             `json:"date_until,omitempty"`
	EmbargoDate             string                 `json:"embargo_date,omitempty"`
	AccessLevelAfterEmbargo string                 `json:"access_level_after_embargo,omitempty"`
	Format                  []string               `json:"format,omitempty"`
	Handle                  string                 `json:"handle,omitempty"`
	ID                      string                 `json:"id,omitempty"`
	Identifiers             Values                 `json:"identifiers,omitempty"`
	Keyword                 []string               `json:"keyword,omitempty"`
	HasBeenPublic           bool                   `json:"has_been_public"`
	Language                []string               `json:"language,omitempty"`
	LastUserID              string                 `json:"last_user_id,omitempty"`
	LastUser                *Person                `json:"-"`
	License                 string                 `json:"license,omitempty"`
	Link                    []*DatasetLink         `json:"link,omitempty"`
	Locked                  bool                   `json:"locked"`
	Message                 string                 `json:"message,omitempty"`
	OtherLicense            string                 `json:"other_license,omitempty"`
	Publisher               string                 `json:"publisher,omitempty"`
	RelatedOrganizations    []*RelatedOrganization `json:"related_organizations,omitempty"`
	RelatedProjects         []*RelatedProject      `json:"related_projects,omitempty"`
	RelatedPublication      []RelatedPublication   `json:"related_publication,omitempty"`
	ReviewerNote            string                 `json:"reviewer_note,omitempty"`
	ReviewerTags            []string               `json:"reviewer_tags,omitempty"`
	SnapshotID              string                 `json:"snapshot_id,omitempty"`
	Status                  string                 `json:"status,omitempty"`
	Title                   string                 `json:"title,omitempty"`
	UserID                  string                 `json:"user_id,omitempty"`
	User                    *Person                `json:"-"`
	Year                    string                 `json:"year,omitempty"`
}

func (d *Dataset) FirstIdentifier() (string, []string) { // TODO eliminate need for this
	for typ, vals := range d.Identifiers {
		return typ, vals
	}
	return "", nil
}

// TODO: remove when we really support multiple identifiers
func (d *Dataset) IdentifierType() string {
	for key := range d.Identifiers {
		return key
	}
	return ""
}

// TODO: remove when we really support multiple identifiers
func (d *Dataset) IdentifierValue() string {
	for _, vals := range d.Identifiers {
		return vals[0]
	}
	return ""
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
		return nil, fmt.Errorf("dataset.GetContributor: index %d out of bounds", i)
	}

	return cc[i], nil
}

func (d *Dataset) AddContributor(role string, c *Contributor) {
	d.SetContributors(role, append(d.Contributors(role), c))
}

func (d *Dataset) SetContributor(role string, i int, c *Contributor) error {
	cc := d.Contributors(role)
	if i >= len(cc) {
		return fmt.Errorf("dataset.SetContributor: index %d out of bounds", i)
	}

	cc[i] = c

	return nil
}

func (d *Dataset) RemoveContributor(role string, i int) error {
	cc := d.Contributors(role)
	if i >= len(cc) {
		return fmt.Errorf("dataset.RemoveContributor: index %d out of bounds", i)
	}

	d.SetContributors(role, append(cc[:i], cc[i+1:]...))

	return nil
}

func (d *Dataset) HasContributor(role string, u *Person) bool {
	cc := d.Contributors(role)

	return slices.ContainsFunc(cc, func(c *Contributor) bool {
		return c.PersonID == u.ID
	})
}

func (d *Dataset) GetLink(id string) *DatasetLink {
	for _, pl := range d.Link {
		if pl.ID == id {
			return pl
		}
	}
	return nil
}

func (d *Dataset) SetLink(l *DatasetLink) {
	for i, link := range d.Link {
		if link.ID == l.ID {
			d.Link[i] = l
		}
	}
}

func (d *Dataset) AddLink(l *DatasetLink) {
	l.ID = ulid.Make().String()
	d.Link = append(d.Link, l)
}

func (d *Dataset) RemoveLink(id string) {
	links := make([]*DatasetLink, 0)
	for _, pl := range d.Link {
		if pl.ID != id {
			links = append(links, pl)
		}
	}
	d.Link = links
}

func (d *Dataset) GetAbstract(id string) *Text {
	for _, abstract := range d.Abstract {
		if abstract.ID == id {
			return abstract
		}
	}
	return nil
}

func (d *Dataset) SetAbstract(t *Text) {
	for i, abstract := range d.Abstract {
		if abstract.ID == t.ID {
			d.Abstract[i] = t
		}
	}
}

func (d *Dataset) AddAbstract(t *Text) {
	t.ID = ulid.Make().String()
	d.Abstract = append(d.Abstract, t)
}

func (d *Dataset) RemoveAbstract(id string) {
	abstracts := make([]*Text, 0)
	for _, abstract := range d.Abstract {
		if abstract.ID != id {
			abstracts = append(abstracts, abstract)
		}
	}
	d.Abstract = abstracts
}

func (d *Dataset) AddProject(project *Project) {
	d.RemoveProject(project.ID)
	d.RelatedProjects = append(d.RelatedProjects, &RelatedProject{
		ProjectID: project.ID,
		Project:   project,
	})
}

func (d *Dataset) RemoveProject(id string) {
	rels := make([]*RelatedProject, 0)
	for _, rel := range d.RelatedProjects {
		if rel.ProjectID != id {
			rels = append(rels, rel)
		}
	}
	d.RelatedProjects = rels
}

func (d *Dataset) AddOrganization(org *Organization) {
	d.RemoveOrganization(org.ID)
	d.RelatedOrganizations = append(d.RelatedOrganizations, &RelatedOrganization{
		OrganizationID: org.ID,
		Organization:   org,
	})
}

func (d *Dataset) RemoveOrganization(id string) {
	rels := make([]*RelatedOrganization, 0)
	for _, rel := range d.RelatedOrganizations {
		if rel.OrganizationID != id {
			rels = append(rels, rel)
		}
	}
	d.RelatedOrganizations = rels
}

func (dl *DatasetLink) Validate() error {
	errs := okay.NewErrors()

	if dl.ID == "" {
		errs.Add(okay.NewError("/id", "id.required"))
	}
	if dl.URL == "" {
		errs.Add(okay.NewError("/url", "url.required"))
	}
	if !slices.Contains(vocabularies.Map["dataset_link_relations"], dl.Relation) {
		errs.Add(okay.NewError("/relation", "relation.invalid"))
	}

	return errs.ErrorOrNil()
}

func (d *Dataset) Validate() error {
	errs := okay.NewErrors()

	if d.ID == "" {
		errs.Add(okay.NewError("/id", "dataset.id.required"))
	}

	if d.Status == "" {
		errs.Add(okay.NewError("/status", "dataset.status.required"))
	} else if !slices.Contains(vocabularies.Map["publication_statuses"], d.Status) {
		errs.Add(okay.NewError("/status", "dataset.status.invalid"))
	}

	if d.Status == "public" && d.AccessLevel == "" {
		errs.Add(okay.NewError("/access_level", "dataset.access_level.required"))
	}
	if d.AccessLevel != "" && !slices.Contains(vocabularies.Map["dataset_access_levels"], d.AccessLevel) {
		errs.Add(okay.NewError("/access_level", "dataset.access_level.invalid"))
	}

	if d.Status == "public" && len(d.Identifiers) == 0 {
		errs.Add(okay.NewError("/identifier_type", "dataset.identifier_type.required"))
		errs.Add(okay.NewError("/identifier", "dataset.identifier.required"))
	}
	for key, vals := range d.Identifiers {
		if key == "" {
			errs.Add(okay.NewError("/identifier_type", "dataset.identifier_type.required"))
		} else if !slices.Contains(vocabularies.Map["dataset_identifier_types"], key) {
			errs.Add(okay.NewError("/identifier_type", "dataset.identifier_type.invalid"))
		}
		if len(vals) == 0 {
			errs.Add(okay.NewError("/identifier", "dataset.identifier.required"))
		}
		for _, val := range vals {
			if val == "" {
				errs.Add(okay.NewError("/identifier", "dataset.identifier.required"))
				break
			}
		}
	}

	if d.Status == "public" && len(d.Format) == 0 {
		errs.Add(okay.NewError("/format", "dataset.format.required"))
	}
	for i, f := range d.Format {
		if f == "" {
			errs.Add(okay.NewError(fmt.Sprintf("/format/%d", i), "dataset.format.invalid"))
		}
	}

	for i, l := range d.Language {
		if !slices.Contains(vocabularies.Map["language_codes"], l) {
			errs.Add(okay.NewError(fmt.Sprintf("/language/%d", i), "dataset.language.invalid"))
		}
	}

	if d.Status == "public" && d.Publisher == "" {
		errs.Add(okay.NewError("/publisher", "dataset.publisher.required"))
	}

	if d.Status == "public" && d.Title == "" {
		errs.Add(okay.NewError("/title", "dataset.title.required"))
	}

	if d.Status == "public" && d.Year == "" {
		errs.Add(okay.NewError("/year", "dataset.year.required"))
	}
	if d.Year != "" && !reYear.MatchString(d.Year) {
		errs.Add(okay.NewError("/year", "dataset.year.invalid"))
	}

	if d.Status == "public" && len(d.Author) == 0 {
		errs.Add(okay.NewError("/author", "dataset.author.required"))
	}

	for i, c := range d.Author {
		var e *okay.Errors
		if errors.As(c.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Add(okay.NewError(fmt.Sprintf("/author/%d%s", i, err.Key), "dataset.author."+err.Rule))
			}
		}
	}

	// at least one ugent author
	if d.Status == "public" {
		var hasUgentAuthors bool = false
		for _, a := range d.Author {
			if a.PersonID != "" {
				hasUgentAuthors = true
				break
			}
		}
		if !hasUgentAuthors {
			errs.Add(okay.NewError("/author", "dataset.author.min_ugent_authors"))
		}
	}

	// license or other_license -> TODO: base error?
	// now "fixed" by (incorrectly) pointing at license
	if d.Status == "public" && d.License == "" && d.OtherLicense == "" {
		errs.Add(okay.NewError("/license", "dataset.license.required"))
	}

	for i, rel := range d.RelatedPublication {
		if rel.ID == "" {
			errs.Add(okay.NewError(fmt.Sprintf("/related_publication/%d/id", i), "dataset.related_publication.required"))
		}
	}

	for i, rel := range d.RelatedProjects {
		if rel.ProjectID == "" {
			errs.Add(okay.NewError(fmt.Sprintf("/project/%d/id", i), "dataset.project.id.required"))
		}
	}

	for i, rel := range d.RelatedOrganizations {
		if rel.OrganizationID == "" {
			errs.Add(okay.NewError(fmt.Sprintf("/department/%d/id", i), "dataset.department.id.required"))
		}
	}

	for i, l := range d.Link {
		var e *okay.Errors
		if errors.As(l.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Add(okay.NewError(fmt.Sprintf("/link/%d%s", i, err.Key), "dataset.link."+err.Rule))
			}
		}
	}

	// TODO dates and co. are only checked when dataset is public
	if d.Status == "public" && d.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
		if d.EmbargoDate == "" {
			errs.Add(okay.NewError("/embargo_date", "dataset.embargo_date.required"))
		} else if _, err := time.Parse("2006-01-02", d.EmbargoDate); err != nil {
			errs.Add(okay.NewError("/embargo_date", "dataset.embargo_date.invalid"))
		}

		invalid := false
		if d.AccessLevelAfterEmbargo == "" {
			invalid = true
			errs.Add(okay.NewError("/access_level_after_embargo", "dataset.access_level_after_embargo.required"))
		}

		if d.AccessLevel == d.AccessLevelAfterEmbargo && !invalid {
			invalid = true
			errs.Add(okay.NewError("/access_level_after_embargo", "dataset.access_level_after_embargo.similar"))
		}

		if !slices.Contains(vocabularies.Map["dataset_access_levels"], d.AccessLevelAfterEmbargo) && !invalid {
			errs.Add(okay.NewError("/access_level_after_embargo", "dataset.access_level_after_embargo.invalid"))
		}
	}

	for i, a := range d.Abstract {
		var e *okay.Errors
		if errors.As(a.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Add(okay.NewError(fmt.Sprintf("/abstract/%d%s", i, err.Key), "dataset.abstract."+err.Rule))
			}
		}
	}

	return errs.ErrorOrNil()
}

func (d *Dataset) ClearEmbargo() {
	d.AccessLevel = d.AccessLevelAfterEmbargo
	d.AccessLevelAfterEmbargo = ""
	d.EmbargoDate = ""
}
