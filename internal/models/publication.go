package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backend/internal/pagination"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type PublicationHits struct {
	pagination.Pagination
	Hits   []*Publication     `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty"`
	License            string     `json:"license,omitempty"`
	ContentType        string     `json:"content_type,omitempty"`
	DateCreated        *time.Time `json:"date_created,omitempty"`
	DateUpdated        *time.Time `json:"date_updated,omitempty"`
	Description        string     `json:"description,omitempty"`
	Embargo            string     `json:"embargo,omitempty"`
	EmbargoTo          string     `json:"embargo_to,omitempty"`
	Name               string     `json:"name,omitempty"`
	Size               int        `json:"size,omitempty"`
	ID                 string     `json:"id,omitempty"`
	SHA256             string     `json:"sha256,omitempty"`
	OtherLicense       string     `json:"other_license,omitempty"`
	PublicationVersion string     `json:"publication_version,omitempty"`
	Relation           string     `json:"relation,omitempty"`
	ThumbnailURL       string     `json:"thumbnail_url,omitempty"`
	Title              string     `json:"title,omitempty"`
	URL                string     `json:"url,omitempty"`
}

type PublicationLink struct {
	ID          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Description string `json:"description,omitempty"`
}

type PublicationDepartmentRef struct {
	ID string `json:"id,omitempty"`
}

type PublicationDepartment struct {
	ID   string                     `json:"id,omitempty"`
	Tree []PublicationDepartmentRef `json:"tree,omitempty"`
}

type PublicationProject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PublicationORCIDWork struct {
	ORCID   string `json:"orcid,omitempty"`
	PutCode int    `json:"put_code,omitempty"`
}

type RelatedDataset struct {
	ID string `json:"id,omitempty"`
}

type Publication struct {
	Abstract                []Text                  `json:"abstract,omitempty"`
	AdditionalInfo          string                  `json:"additional_info,omitempty"`
	AlternativeTitle        []string                `json:"alternative_title,omitempty"`
	ArticleNumber           string                  `json:"article_number,omitempty"`
	ArxivID                 string                  `json:"arxiv_id,omitempty"`
	Author                  []*Contributor          `json:"author,omitempty"`
	BatchID                 string                  `json:"batch_id,omitempty"`
	Classification          string                  `json:"classification,omitempty"`
	CompletenessScore       int                     `json:"completeness_score,omitempty"`
	ConferenceName          string                  `json:"conference_name,omitempty"`
	ConferenceLocation      string                  `json:"conference_location,omitempty"`
	ConferenceOrganizer     string                  `json:"conference_organizer,omitempty"`
	ConferenceStartDate     string                  `json:"conference_start_date,omitempty"`
	ConferenceEndDate       string                  `json:"conference_end_date,omitempty"`
	ConferenceType          string                  `json:"conference_type,omitempty"`
	CreatorID               string                  `json:"creator_id,omitempty"`
	DateCreated             *time.Time              `json:"date_created,omitempty"`
	DateUpdated             *time.Time              `json:"date_updated,omitempty"`
	DateFrom                *time.Time              `json:"date_from,omitempty"`
	DateUntil               *time.Time              `json:"date_until,omitempty"`
	DefenseDate             string                  `json:"defense_date,omitempty"`
	DefensePlace            string                  `json:"defense_place,omitempty"`
	DefenseTime             string                  `json:"defense_time,omitempty"`
	Department              []PublicationDepartment `json:"department,omitempty"`
	DOI                     string                  `json:"doi,omitempty"`
	Edition                 string                  `json:"edition,omitempty"`
	Editor                  []*Contributor          `json:"editor,omitempty"`
	EISBN                   []string                `json:"eisbn,omitempty"`
	EISSN                   []string                `json:"eissn,omitempty"`
	ESCIID                  string                  `json:"esci_id,omitempty"`
	Extern                  bool                    `json:"extern,omitempty"`
	File                    []*PublicationFile      `json:"file,omitempty"`
	Handle                  string                  `json:"handle,omitempty"`
	HasConfidentialData     string                  `json:"has_confidential_data,omitempty"`
	HasPatentApplication    string                  `json:"has_patent_application,omitempty"`
	HasPublicationsPlanned  string                  `json:"has_publications_planned,omitempty"`
	HasPublishedMaterial    string                  `json:"has_published_material,omitempty"`
	ID                      string                  `json:"id,omitempty"`
	ISBN                    []string                `json:"isbn,omitempty"`
	ISSN                    []string                `json:"issn,omitempty"`
	Issue                   string                  `json:"issue,omitempty"`
	IssueTitle              string                  `json:"issue_title,omitempty"`
	JournalArticleType      string                  `json:"journal_article_type,omitempty"`
	Keyword                 []string                `json:"keyword,omitempty"`
	Language                []string                `json:"language,omitempty"`
	LaySummary              []Text                  `json:"lay_summary,omitempty"`
	Link                    []PublicationLink       `json:"link,omitempty"`
	Locked                  bool                    `json:"locked,omitempty"`
	Message                 string                  `json:"message,omitempty"`
	MiscellaneousType       string                  `json:"miscellaneous_type,omitempty"`
	ORCIDWork               []PublicationORCIDWork  `json:"orcid_work,omitempty"`
	PageCount               string                  `json:"page_count,omitempty"`
	PageFirst               string                  `json:"page_first,omitempty"`
	PageLast                string                  `json:"page_last,omitempty"`
	PlaceOfPublication      string                  `json:"place_of_publication,omitempty"`
	Project                 []PublicationProject    `json:"project,omitempty"`
	Publication             string                  `json:"publication,omitempty"`
	PublicationAbbreviation string                  `json:"publication_abbreviation,omitempty"`
	PublicationStatus       string                  `json:"publication_status,omitempty"`
	Publisher               string                  `json:"publisher,omitempty"`
	PubMedID                string                  `json:"pubmed_id,omitempty"`
	RelatedDataset          []RelatedDataset        `json:"related_dataset,omitempty"`
	ReportNumber            string                  `json:"report_number,omitempty"`
	ResearchField           []string                `json:"research_field,omitempty"`
	ReviewerNote            string                  `json:"reviewer_note,omitempty"`
	ReviewerTags            []string                `json:"reviewer_tags,omitempty"`
	SeriesTitle             string                  `json:"series_title,omitempty"`
	SnapshotID              string                  `json:"-"`
	SourceDB                string                  `json:"source_db,omitempty"`
	SourceID                string                  `json:"source_id,omitempty"`
	SourceRecord            string                  `json:"source_record,omitempty"`
	Status                  string                  `json:"status,omitempty"`
	Supervisor              []*Contributor          `json:"supervisor,omitempty"`
	Title                   string                  `json:"title,omitempty"`
	Type                    string                  `json:"type,omitempty"`
	URL                     string                  `json:"url,omitempty"`
	UserID                  string                  `json:"user_id,omitempty"`
	Volume                  string                  `json:"volume,omitempty"`
	WOSID                   string                  `json:"wos_id,omitempty"`
	WOSType                 string                  `json:"wos_type,omitempty"`
	Year                    string                  `json:"year,omitempty"`
}

func (p *Publication) AccessLevel() string {
	for _, a := range []string{"open_access", "local", "closed"} {
		for _, file := range p.File {
			if file.AccessLevel == a {
				return a
			}
		}
	}
	return ""
}

func (p *Publication) HasRelatedDataset(id string) bool {
	for _, r := range p.RelatedDataset {
		if r.ID == id {
			return true
		}
	}
	return false
}

func (p *Publication) RemoveRelatedDataset(id string) {
	var datasets []RelatedDataset
	for _, r := range p.RelatedDataset {
		if r.ID != id {
			datasets = append(datasets, r)
		}
	}
	p.RelatedDataset = datasets
}

func (p *Publication) GetFile(id string) *PublicationFile {
	for _, file := range p.File {
		if file.ID == id {
			return file
		}
	}
	return nil
}

func (p *Publication) ThumbnailURL() string {
	for _, file := range p.File {
		if file.ThumbnailURL != "" {
			return file.ThumbnailURL
		}
	}
	return ""
}

func (p *Publication) ClassificationChoices() []string {
	switch p.Type {
	case "journal_article":
		return []string{
			"U",
			"A1",
			"A2",
			"A3",
			"A4",
			"V",
		}
	case "book":
		return []string{
			"U",
			"B1",
		}
	case "book_chapter":
		return []string{
			"U",
			"B2",
		}
	case "book_editor", "issue_editor":
		return []string{
			"U",
			"B3",
		}
	case "conference":
		return []string{
			"U",
			"P1",
			"C1",
			"C3",
		}
	case "dissertation":
		return []string{
			"U",
			"D1",
		}
	case "miscellaneous", "report", "preprint":
		return []string{
			"U",
			"V",
		}
	default:
		return []string{
			"U",
		}
	}
}

func (p *Publication) Contributors(role string) []*Contributor {
	switch role {
	case "author":
		return p.Author
	case "editor":
		return p.Editor
	case "supervisor":
		return p.Supervisor
	default:
		return nil
	}
}

func (p *Publication) SetContributors(role string, c []*Contributor) {
	switch role {
	case "author":
		p.Author = c
	case "editor":
		p.Editor = c
	case "supervisor":
		p.Supervisor = c
	}
}

func (p *Publication) GetContributor(role string, i int) (*Contributor, error) {
	cc := p.Contributors(role)
	if i >= len(cc) {
		return nil, errors.New("index out of bounds")
	}

	return cc[i], nil
}

func (p *Publication) AddContributor(role string, c *Contributor) {
	p.SetContributors(role, append(p.Contributors(role), c))
}

func (p *Publication) SetContributor(role string, i int, c *Contributor) error {
	cc := p.Contributors(role)
	if i >= len(cc) {
		return errors.New("index out of bounds")
	}

	cc[i] = c

	return nil
}

func (p *Publication) RemoveContributor(role string, i int) error {
	cc := p.Contributors(role)
	if i >= len(cc) {
		return errors.New("index out of bounds")
	}

	p.SetContributors(role, append(cc[:i], cc[i+1:]...))

	return nil
}

func (p *Publication) GetLink(id string) *PublicationLink {
	for _, pl := range p.Link {
		if pl.ID == id {
			return &pl
		}
	}
	return nil
}

func (p *Publication) AddLink(l *PublicationLink) {
	l.ID = ulid.MustGenerate()
	p.Link = append(p.Link, *l)
}

func (p *Publication) RemoveLink(id string) {
	links := make([]PublicationLink, 0)
	for _, pl := range p.Link {
		if pl.ID != id {
			links = append(links, pl)
		}
	}
	p.Link = links
}

func (p *Publication) GetAbstract(id string) *Text {
	for _, abstract := range p.Abstract {
		if abstract.ID == id {
			return &abstract
		}
	}
	return nil
}

func (p *Publication) AddAbstract(t *Text) {
	t.ID = ulid.MustGenerate()
	p.Abstract = append(p.Abstract, *t)
}

func (p *Publication) RemoveAbstract(id string) {
	abstracts := make([]Text, 0)
	for _, abstract := range p.Abstract {
		if abstract.ID != id {
			abstracts = append(abstracts, abstract)
		}
	}
	p.Abstract = abstracts
}

func (p *Publication) GetLaySummary(id string) *Text {
	for _, ls := range p.LaySummary {
		if ls.ID == id {
			return &ls
		}
	}
	return nil
}

func (p *Publication) AddLaySummary(t *Text) {
	t.ID = ulid.MustGenerate()
	p.LaySummary = append(p.LaySummary, *t)
}

func (p *Publication) RemoveLaySummary(id string) {
	lay_summaries := make([]Text, 0)
	for _, ls := range p.LaySummary {
		if ls.ID != id {
			lay_summaries = append(lay_summaries, ls)
		}
	}
	p.LaySummary = lay_summaries
}

func (p *Publication) GetProject(id string) *PublicationProject {
	for _, p := range p.Project {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func (p *Publication) RemoveProject(id string) {
	projects := make([]PublicationProject, 0)
	for _, pl := range p.Project {
		if pl.ID != id {
			projects = append(projects, pl)
		}
	}
	p.Project = projects
}

func (p *Publication) AddProject(pr *PublicationProject) {
	p.RemoveProject(pr.ID)
	p.Project = append(p.Project, *pr)
}

func (p *Publication) RemoveDepartment(id string) {
	deps := make([]PublicationDepartment, 0)
	for _, dep := range p.Department {
		if dep.ID != id {
			deps = append(deps, dep)
		}
	}
	p.Department = deps
}

func (p *Publication) AddDepartmentByOrg(org *Organization) {
	// remove if added before
	p.RemoveDepartment(org.ID)

	publicationDepartment := PublicationDepartment{ID: org.ID}
	for _, d := range org.Tree {
		publicationDepartment.Tree = append(publicationDepartment.Tree, PublicationDepartmentRef(d))
	}
	p.Department = append(p.Department, publicationDepartment)
}

func (p *Publication) AddFile(file *PublicationFile) {
	file.ID = ulid.MustGenerate()
	now := time.Now()
	file.DateCreated = &now
	file.DateUpdated = &now
	if file.AccessLevel == "" {
		file.AccessLevel = "local"
	}
	p.File = append(p.File, file)
}

func (p *Publication) RemoveFile(id string) {
	newFile := []*PublicationFile{}
	for _, f := range p.File {
		if f.ID != id {
			newFile = append(newFile, f)
		}
	}
	p.File = newFile
}

func (p *Publication) UsesLaySummary() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesConference() bool {
	switch p.Type {
	case "book_chapter", "book_editor", "conference", "issue_editor", "journal_article":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesContributors(role string) bool {
	switch role {
	case "author":
		return p.UsesAuthor()
	case "editor":
		return p.UsesEditor()
	case "supervisor":
		return p.UsesSupervisor()
	default:
		return false
	}
}

func (p *Publication) UsesAuthor() bool {
	switch p.Type {
	case "book", "book_chapter", "conference", "dissertation", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesEditor() bool {
	switch p.Type {
	case "book", "book_chapter", "book_editor", "conference", "issue_editor", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesSupervisor() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) InORCIDWorks(orcidID string) bool {
	for _, w := range p.ORCIDWork {
		if w.ORCID == orcidID {
			return true
		}
	}
	return false
}

func (p *Publication) Validate() error {
	var errs validation.Errors

	if p.ID == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/id",
			Code:    "publication.id.required",
		})
	}
	if p.Type == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/type",
			Code:    "publication.type.required",
		})
	} else if !validation.IsPublicationType(p.Type) {
		errs = append(errs, &validation.Error{
			Pointer: "/type",
			Code:    "publication.type.invalid",
		})
	}
	// TODO check classification validity
	if p.Classification == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/classification",
			Code:    "publication.classification.required",
		})
	}
	if p.Status == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "publication.status.required",
		})
	} else if !validation.IsStatus(p.Status) {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "publication.status.invalid",
		})
	}

	if p.Status == "public" && p.Title == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/title",
			Code:    "publication.title.required",
		})
	}

	if p.Status == "public" {
		if p.Year == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "publication.year.required",
			})
		} else if !validation.IsYear(p.Year) {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "publication.year.invalid",
			})
		}
	}

	if p.Status == "public" {
		for i, a := range p.Abstract {
			for _, err := range a.Validate() {
				errs = append(errs, &validation.Error{
					Pointer: fmt.Sprintf("/abstract/%d%s", i, err.Pointer),
					Code:    "publication.abstract." + err.Code,
				})
			}
		}
	}

	if p.Status == "public" && p.UsesLaySummary() {
		for i, l := range p.LaySummary {
			for _, err := range l.Validate() {
				errs = append(errs, &validation.Error{
					Pointer: fmt.Sprintf("/lay_summary/%d%s", i, err.Pointer),
					Code:    "publication.lay_summary." + err.Code,
				})
			}
		}
	}

	if p.Status == "public" && p.UsesAuthor() && len(p.Author) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/author",
			Code:    "publication.author.required",
		})

	}

	// at least one ugent author if not external
	if p.Status == "public" && p.UsesAuthor() && !p.Extern {
		var hasUgentAuthors bool = false
		for _, a := range p.Author {
			if a.ID != "" {
				hasUgentAuthors = true
				break
			}
		}
		if !hasUgentAuthors {
			errs = append(errs, &validation.Error{
				Pointer: "/author",
				Code:    "publication.author.min_ugent_authors",
			})
		}
	}

	if p.Status == "public" && p.UsesEditor() && !p.UsesAuthor() && len(p.Editor) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/editor",
			Code:    "publication.editor.required",
		})

		// at least one ugent editor if not external
		if !p.Extern {
			var hasUgentEditors bool = false
			for _, a := range p.Editor {
				if a.ID != "" {
					hasUgentEditors = true
					break
				}
			}
			if !hasUgentEditors {
				errs = append(errs, &validation.Error{
					Pointer: "/editor",
					Code:    "publication.editor.min_ugent_editors",
				})
			}
		}
	}

	// at least one ugent editor for editor types if not external
	if p.Status == "public" && p.UsesEditor() && !p.UsesAuthor() && !p.Extern {
		var hasUgentEditors bool = false
		for _, a := range p.Editor {
			if a.ID != "" {
				hasUgentEditors = true
				break
			}
		}
		if !hasUgentEditors {
			errs = append(errs, &validation.Error{
				Pointer: "/editor",
				Code:    "publication.editor.min_ugent_editors",
			})
		}
	}

	if p.Status == "public" && p.UsesSupervisor() && len(p.Supervisor) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/supervisor",
			Code:    "publication.supervisor.required",
		})
	}

	for i, c := range p.Author {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/author/%d%s", i, err.Pointer),
				Code:    "publication.author." + err.Code,
			})
		}
	}
	for i, c := range p.Editor {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/editor/%d%s", i, err.Pointer),
				Code:    "publication.editor." + err.Code,
			})
		}
	}
	for i, c := range p.Supervisor {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/supervisor/%d%s", i, err.Pointer),
				Code:    "publication.supervisor." + err.Code,
			})
		}
	}

	for i, pr := range p.Project {
		if pr.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/project/%d/id", i),
				Code:    "publication.project.id.required",
			})
		}
	}

	for i, dep := range p.Department {
		if dep.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/department/%d/id", i),
				Code:    "publication.department.id.required",
			})
		}
	}

	for i, rd := range p.RelatedDataset {
		if rd.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/related_dataset/%d/id", i),
				Code:    "publication.related_dataset.id.required",
			})
		}
	}

	for i, f := range p.File {
		for _, err := range f.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/file/%d%s", i, err.Pointer),
				Code:    "publication.file" + err.Code,
			})
		}
	}

	if p.Status == "public" {
		for i, pl := range p.Link {
			for _, err := range pl.Validate() {
				errs = append(errs, &validation.Error{
					Pointer: fmt.Sprintf("/link/%d%s", i, err.Pointer),
					Code:    "publication.link" + err.Code,
				})
			}
		}
	}

	// type specific validation
	switch p.Type {
	case "dissertation":
		errs = append(errs, p.validateDissertation()...)
	case "journal_article":
		errs = append(errs, p.validateJournalArticle()...)
	case "miscellaneous":
		errs = append(errs, p.validateMiscellaneous()...)
	case "book":
		errs = append(errs, p.validateBook()...)
	case "book_chapter":
		errs = append(errs, p.validateBookChapter()...)
	case "conference":
		errs = append(errs, p.validateConference()...)
	case "book_editor":
		errs = append(errs, p.validateBookEditor()...)
	case "issue_editor":
		errs = append(errs, p.validateIssueEditor()...)
	}

	// TODO: why is the nil slice validationErrors(nil) != nil in mutantdb validation?
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (p *Publication) validateBookEditor() (errs validation.Errors) {
	return
}

func (p *Publication) validateIssueEditor() (errs validation.Errors) {
	return
}

func (p *Publication) validateJournalArticle() (errs validation.Errors) {
	// TODO: confusing: gui shows select without empty element
	// but first creation sets this value to empty
	if p.JournalArticleType != "" && !validation.InArray(vocabularies.Map["journal_article_types"], p.JournalArticleType) {
		errs = append(errs, &validation.Error{
			Pointer: "/journal_article_type",
			Code:    "publication.journal_article_type.invalid",
		})
	}
	if p.Status != "public" {
		return
	}
	if p.Publication == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/publication",
			Code:    "publication.journal_article.publication.required",
		})
	}
	return
}

func (p *Publication) validateBook() (errs validation.Errors) {
	return
}

func (p *Publication) validateConference() (errs validation.Errors) {
	return
}

func (p *Publication) validateBookChapter() (errs validation.Errors) {
	return
}

func (p *Publication) validateDissertation() (errs validation.Errors) {
	if p.Status != "public" {
		return
	}
	if p.DefensePlace == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_place",
			Code:    "publication.defense_place.required",
		})
	}
	if p.DefenseDate == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_date",
			Code:    "publication.defense_date.required",
		})
	} else if !validation.IsDate(p.DefenseDate) {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_date",
			Code:    "publication.defense_date.invalid",
		})
	}
	if p.DefenseTime == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_time",
			Code:    "publication.defense_time.required",
		})
	} else if !validation.IsTime(p.DefenseTime) {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_time",
			Code:    "publication.defense_time.invalid",
		})
	}
	return
}

func (p *Publication) validateMiscellaneous() (errs validation.Errors) {
	// TODO confusing: gui shows select without empty element
	// but first creation sets this value to empty
	if p.MiscellaneousType != "" && !validation.InArray(vocabularies.Map["miscellaneous_types"], p.MiscellaneousType) {
		errs = append(errs, &validation.Error{
			Pointer: "/miscellaneous_type",
			Code:    "publication.miscellaneous_type.invalid",
		})
	}
	return
}

func (pf *PublicationFile) Validate() (errs validation.Errors) {
	if !validation.InArray(vocabularies.Map["publication_file_access_levels"], pf.AccessLevel) {
		errs = append(errs, &validation.Error{
			Pointer: "/access_level",
			Code:    "access_level.invalid",
		})
	}

	if pf.ContentType == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/content_type",
			Code:    "content_type.required",
		})
	}

	if pf.ID == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/id",
			Code:    "id.required",
		})
	}

	if pf.Relation != "" && !validation.InArray(vocabularies.Map["publication_file_relations"], pf.Relation) {
		errs = append(errs, &validation.Error{
			Pointer: "/relation",
			Code:    "relation.invalid",
		})
	}

	if pf.Embargo != "" {
		if !validation.IsDate(pf.Embargo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo",
				Code:    "embargo.invalid",
			})
		}
		if pf.EmbargoTo == pf.AccessLevel {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "embargo_to.invalid", // TODO better code
			})
		} else if !validation.InArray(vocabularies.Map["publication_file_access_levels"], pf.EmbargoTo) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_to",
				Code:    "embargo_to.invalid", // TODO better code
			})
		}
	}

	if pf.PublicationVersion != "" && !validation.InArray(vocabularies.Map["publication_versions"], pf.PublicationVersion) {
		errs = append(errs, &validation.Error{
			Pointer: "/publication_version",
			Code:    "publication_version.invalid",
		})
	}

	return
}

func (pl *PublicationLink) Validate() (errs validation.Errors) {
	if !validation.InArray(vocabularies.Map["publication_link_relations"], pl.Relation) {
		errs = append(errs, &validation.Error{
			Pointer: "/relation",
			Code:    "relation.invalid",
		})
	}
	return
}
