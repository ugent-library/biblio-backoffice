package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backend/internal/pagination"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

type PublicationHits struct {
	pagination.Pagination
	Hits   []*Publication     `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type PublicationFile struct {
	AccessLevel        string     `json:"access_level,omitempty" form:"access_level"`
	CCLicense          string     `json:"cc_license,omitempty" form:"cc_license"`
	ContentType        string     `json:"content_type,omitempty" form:"-"`
	DateCreated        *time.Time `json:"date_created,omitempty" form:"-"`
	DateUpdated        *time.Time `json:"date_updated,omitempty" form:"-"`
	Description        string     `json:"description,omitempty" form:"description"`
	Embargo            string     `json:"embargo,omitempty" form:"embargo"`
	EmbargoTo          string     `json:"embargo_to,omitempty" form:"embargo_to"`
	Filename           string     `json:"file_name,omitempty" form:"-"`
	FileSize           int        `json:"file_size,omitempty" form:"-"`
	ID                 string     `json:"file_id,omitempty" form:"-"`
	SHA256             string     `json:"sha256,omitempty" form:"-"`
	NoLicense          bool       `json:"no_license,omitempty" form:"no_license"`
	OtherLicense       string     `json:"other_license,omitempty" form:"other_license"`
	PublicationVersion string     `json:"publication_version,omitempty" form:"publication_version"`
	Relation           string     `json:"relation,omitempty" form:"relation"`
	ThumbnailURL       string     `json:"thumbnail_url,omitempty" form:"-"`
	Title              string     `json:"title,omitempty" form:"title"`
	URL                string     `json:"url,omitempty" form:"-"`
}

func (f *PublicationFile) Clone() *PublicationFile {
	clone := *f
	return &clone
}

type PublicationLink struct {
	URL         string `json:"url,omitempty" form:"url"`
	Relation    string `json:"relation,omitempty" form:"relation"`
	Description string `json:"description,omitempty" form:"description"`
}

type PublicationConference struct {
	Name      string `json:"name,omitempty" form:"name"`
	Location  string `json:"location,omitempty" form:"location"`
	Organizer string `json:"organizer,omitempty" form:"organizer"`
	StartDate string `json:"start_date,omitempty" form:"start_date"`
	EndDate   string `json:"end_date,omitempty" form:"end_date"`
}

type PublicationDepartmentRef struct {
	ID string `json:"_id,omitempty"`
}

type PublicationDepartment struct {
	ID   string                     `json:"_id,omitempty"`
	Tree []PublicationDepartmentRef `json:"tree,omitempty"`
}

type PublicationProject struct {
	ID   string `json:"_id,omitempty"`
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
	Abstract                []Text                  `json:"abstract,omitempty" form:"abstract"`
	AdditionalInfo          string                  `json:"additional_info,omitempty" form:"additional_info"`
	AlternativeTitle        []string                `json:"alternative_title,omitempty" form:"alternative_title"`
	ArticleNumber           string                  `json:"article_number,omitempty" form:"article_number"`
	ArxivID                 string                  `json:"arxiv_id,omitempty" form:"arxiv_id"`
	Author                  []*Contributor          `json:"author,omitempty" form:"-"`
	BatchID                 string                  `json:"batch_id,omitempty" form:"-"`
	Classification          string                  `json:"classification,omitempty" form:"classification"`
	CompletenessScore       int                     `json:"completeness_score,omitempty" form:"-"`
	Conference              PublicationConference   `json:"conference,omitempty" form:"conference"` // TODO should be pointer or just fields
	ConferenceType          string                  `json:"conference_type,omitempty" form:"conference_type"`
	CreatorID               string                  `json:"creator_id,omitempty" form:"-"`
	DateCreated             *time.Time              `json:"date_created,omitempty" form:"-"`
	DateUpdated             *time.Time              `json:"date_updated,omitempty" form:"-"`
	DateFrom                *time.Time              `json:"date_from,omitempty" form:"-"`
	DateUntil               *time.Time              `json:"date_until,omitempty" form:"-"`
	DefenseDate             string                  `json:"defense_date,omitempty" form:"defense_date"`
	DefensePlace            string                  `json:"defense_place,omitempty" form:"defense_place"`
	DefenseTime             string                  `json:"defense_time,omitempty" form:"defense_time"`
	Department              []PublicationDepartment `json:"department,omitempty" form:"-"`
	DOI                     string                  `json:"doi,omitempty" form:"doi"`
	Edition                 string                  `json:"edition,omitempty" form:"edition"`
	Editor                  []*Contributor          `json:"editor,omitempty" form:"-"`
	EISBN                   []string                `json:"eisbn,omitempty" form:"eisbn"`
	EISSN                   []string                `json:"eissn,omitempty" form:"eissn"`
	ESCIID                  string                  `json:"esci_id,omitempty" form:"esci_id"`
	Extern                  bool                    `json:"extern,omitempty" form:"extern"`
	File                    []*PublicationFile      `json:"file,omitempty" form:"-"`
	Handle                  string                  `json:"handle,omitempty" form:"-"`
	HasConfidentialData     string                  `json:"has_confidential_data,omitempty" form:"has_confidential_data"`
	HasPatentApplication    string                  `json:"has_patent_application,omitempty" form:"has_patent_application"`
	HasPublicationsPlanned  string                  `json:"has_publications_planned,omitempty" form:"has_publications_planned"`
	HasPublishedMaterial    string                  `json:"has_published_material,omitempty" form:"has_published_material"`
	ID                      string                  `json:"id,omitempty" form:"-"`
	ISBN                    []string                `json:"isbn,omitempty" form:"isbn"`
	ISSN                    []string                `json:"issn,omitempty" form:"issn"`
	Issue                   string                  `json:"issue,omitempty" form:"issue"`
	IssueTitle              string                  `json:"issue_title,omitempty" form:"issue_title"`
	JournalArticleType      string                  `json:"journal_article_type,omitempty" form:"journal_article_type"`
	Keyword                 []string                `json:"keyword,omitempty" form:"keyword"`
	Language                []string                `json:"language,omitempty" form:"language"`
	LaySummary              []Text                  `json:"lay_summary,omitempty" form:"lay_summary"`
	Link                    []PublicationLink       `json:"link,omitempty" form:"-"`
	Locked                  bool                    `json:"locked,omitempty" form:"-"`
	Message                 string                  `json:"message,omitempty" form:"-"`
	MiscellaneousType       string                  `json:"miscellaneous_type,omitempty" form:"miscellaneous_type"`
	ORCIDWork               []PublicationORCIDWork  `json:"orcid_work,omitempty" form:"-"`
	PageCount               string                  `json:"page_count,omitempty" form:"page_count"`
	PageFirst               string                  `json:"page_first,omitempty" form:"page_first"`
	PageLast                string                  `json:"page_last,omitempty" form:"page_last"`
	PlaceOfPublication      string                  `json:"place_of_publication,omitempty" form:"place_of_publication"`
	Project                 []PublicationProject    `json:"project,omitempty" form:"-"`
	Publication             string                  `json:"publication,omitempty" form:"publication"`
	PublicationAbbreviation string                  `json:"publication_abbreviation,omitempty" form:"publication_abbreviation"`
	PublicationStatus       string                  `json:"publication_status,omitempty" form:"publication_status"`
	Publisher               string                  `json:"publisher,omitempty" form:"publisher"`
	PubMedID                string                  `json:"pubmed_id,omitempty" form:"pubmed_id"`
	RelatedDataset          []RelatedDataset        `json:"related_dataset,omitempty" form:"-"`
	ReportNumber            string                  `json:"report_number,omitempty" form:"report_number"`
	ResearchField           []string                `json:"research_field,omitempty" form:"research_field"`
	ReviewerNote            string                  `json:"reviewer_note,omitempty" form:"-"`
	ReviewerTags            []string                `json:"reviewer_tags,omitempty" form:"-"`
	SeriesTitle             string                  `json:"series_title,omitempty" form:"series_title"`
	SnapshotID              string                  `json:"-" form:"-"`
	SourceDB                string                  `json:"source_db,omitempty" form:"-"`
	SourceID                string                  `json:"source_id,omitempty" form:"-"`
	SourceRecord            string                  `json:"source_record,omitempty" form:"-"`
	Status                  string                  `json:"status,omitempty" form:"-"`
	Supervisor              []*Contributor          `json:"supervisor,omitempty" form:"-"`
	Title                   string                  `json:"title,omitempty" form:"title"`
	Type                    string                  `json:"type,omitempty" form:"-"`
	URL                     string                  `json:"url,omitempty" form:"url"`
	UserID                  string                  `json:"user_id,omitempty" form:"-"`
	Volume                  string                  `json:"volume,omitempty" form:"volume"`
	WOSID                   string                  `json:"wos_id,omitempty" form:"wos_id"`
	WOSType                 string                  `json:"wos_type,omitempty" form:"-"`
	Year                    string                  `json:"year,omitempty" form:"year"`
}

func (p *Publication) Clone() *Publication {
	clone := *p
	clone.Abstract = nil
	clone.Abstract = append(clone.Abstract, p.Abstract...)
	clone.AlternativeTitle = nil
	clone.AlternativeTitle = append(clone.AlternativeTitle, p.AlternativeTitle...)
	clone.Author = nil
	for _, c := range p.Author {
		clone.Author = append(clone.Author, c.Clone())
	}
	clone.Department = nil
	clone.Department = append(clone.Department, p.Department...)
	clone.Editor = nil
	for _, c := range p.Editor {
		clone.Editor = append(clone.Editor, c.Clone())
	}
	clone.EISBN = nil
	clone.EISBN = append(clone.EISBN, p.EISBN...)
	clone.EISSN = nil
	clone.EISSN = append(clone.EISSN, p.EISSN...)
	clone.File = nil
	for _, f := range p.File {
		clone.File = append(clone.File, f.Clone())
	}
	clone.ISBN = nil
	clone.ISBN = append(clone.ISBN, p.ISBN...)
	clone.ISSN = nil
	clone.ISSN = append(clone.ISSN, p.ISSN...)
	clone.Keyword = nil
	clone.Keyword = append(clone.Keyword, p.Keyword...)
	clone.Language = nil
	clone.Language = append(clone.Language, p.Language...)
	clone.LaySummary = nil
	clone.LaySummary = append(clone.LaySummary, p.LaySummary...)
	clone.Link = nil
	clone.Link = append(clone.Link, p.Link...)
	clone.ORCIDWork = nil
	clone.ORCIDWork = append(clone.ORCIDWork, p.ORCIDWork...)
	clone.Project = nil
	clone.Project = append(clone.Project, p.Project...)
	clone.RelatedDataset = nil
	clone.RelatedDataset = append(clone.RelatedDataset, p.RelatedDataset...)
	clone.ResearchField = nil
	clone.ResearchField = append(clone.ResearchField, p.ResearchField...)
	clone.ReviewerTags = nil
	clone.ReviewerTags = append(clone.ReviewerTags, p.ReviewerTags...)
	clone.Supervisor = nil
	for _, c := range p.Supervisor {
		clone.Supervisor = append(clone.Supervisor, c.Clone())
	}
	return &clone
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

func (p *Publication) AddContributor(role string, i int, c *Contributor) {
	cc := p.Contributors(role)

	if len(cc) == i {
		p.SetContributors(role, append(cc, c))
		return
	}

	newCC := append(cc[:i+1], cc[i:]...)
	newCC[i] = c
	p.SetContributors(role, newCC)
}

func (p *Publication) RemoveContributor(role string, i int) {
	cc := p.Contributors(role)

	p.SetContributors(role, append(cc[:i], cc[i+1:]...))
}

func (p *Publication) GetAbstract(i int) (Text, error) {
	if i >= len(p.Abstract) {
		return Text{}, errors.New("index out of bounds")
	}

	return p.Abstract[i], nil
}

func (p *Publication) SetAbstract(i int, t Text) error {
	if i >= len(p.Abstract) {
		return errors.New("index out of bounds")
	}

	p.Abstract[i] = t

	return nil
}

func (p *Publication) RemoveAbstract(i int) error {
	if i >= len(p.Abstract) {
		return errors.New("index out of bounds")
	}

	p.Abstract = append(p.Abstract[:i], p.Abstract[i+1:]...)

	return nil
}

func (p *Publication) GetLaySummary(i int) (Text, error) {
	if i >= len(p.LaySummary) {
		return Text{}, errors.New("index out of bounds")
	}

	return p.LaySummary[i], nil
}

func (p *Publication) SetLaySummary(i int, t Text) error {
	if i >= len(p.LaySummary) {
		return errors.New("index out of bounds")
	}

	p.LaySummary[i] = t

	return nil
}

func (p *Publication) RemoveLaySummary(i int) error {
	if i >= len(p.LaySummary) {
		return errors.New("index out of bounds")
	}

	p.LaySummary = append(p.LaySummary[:i], p.LaySummary[i+1:]...)

	return nil
}

func (p *Publication) GetProject(i int) (PublicationProject, error) {
	if i >= len(p.Project) {
		return PublicationProject{}, errors.New("index out of bounds")
	}

	return p.Project[i], nil
}

func (p *Publication) RemoveProject(i int) error {
	if i >= len(p.Project) {
		return errors.New("index out of bounds")
	}

	p.Project = append(p.Project[:i], p.Project[i+1:]...)

	return nil
}

func (p *Publication) RemoveDepartment(i int) error {
	if i >= len(p.Department) {
		return errors.New("index out of bounds")
	}

	p.Department = append(p.Department[:i], p.Department[i+1:]...)

	return nil
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

func (d *Publication) Validate() error {
	var errs validation.Errors

	if d.ID == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/id",
			Code:    "publication.id.required",
		})
	}
	if d.Type == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/type",
			Code:    "publication.type.required",
		})
	} else if !validation.IsPublicationType(d.Type) {
		errs = append(errs, &validation.Error{
			Pointer: "/type",
			Code:    "publication.type.invalid",
		})
	}
	// TODO check classification validity
	if d.Classification == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/classification",
			Code:    "publication.classification.required",
		})
	}
	if d.Status == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "publication.status.required",
		})
	} else if !validation.IsStatus(d.Status) {
		errs = append(errs, &validation.Error{
			Pointer: "/status",
			Code:    "publication.status.invalid",
		})
	}

	if d.Status == "public" && d.Title == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/title",
			Code:    "publication.title.required",
		})
	}

	if d.Status == "public" {
		if d.Year == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "publication.year.required",
			})
		} else if !validation.IsYear(d.Year) {
			errs = append(errs, &validation.Error{
				Pointer: "/year",
				Code:    "publication.year.invalid",
			})
		}
	}

	for i, a := range d.Abstract {
		for _, err := range a.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/abstract/%d%s", i, err.Pointer),
				Code:    "publication.abstract." + err.Code,
			})
		}
	}

	for i, l := range d.LaySummary {
		for _, err := range l.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/lay_summary/%d%s", i, err.Pointer),
				Code:    "publication.lay_summary." + err.Code,
			})
		}
	}

	for i, c := range d.Author {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/author/%d%s", i, err.Pointer),
				Code:    "publication.author." + err.Code,
			})
		}
	}
	for i, c := range d.Editor {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/editor/%d%s", i, err.Pointer),
				Code:    "publication.editor." + err.Code,
			})
		}
	}
	for i, c := range d.Supervisor {
		for _, err := range c.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/supervisor/%d%s", i, err.Pointer),
				Code:    "publication.supervisor." + err.Code,
			})
		}
	}

	for i, pr := range d.Project {
		if pr.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/project/%d/id", i),
				Code:    "publication.project.id.required",
			})
		}
	}

	for i, dep := range d.Department {
		if dep.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/department/%d/id", i),
				Code:    "publication.department.id.required",
			})
		}
	}

	for i, rd := range d.RelatedDataset {
		if rd.ID == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/related_dataset/%d/id", i),
				Code:    "publication.related_dataset.id.required",
			})
		}
	}

	for i, f := range d.File {
		for _, err := range f.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/file/%d%s", i, err.Pointer),
				Code:    "publication.file" + err.Code,
			})
		}
	}

	for i, pl := range d.Link {
		for _, err := range pl.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/link/%d%s", i, err.Pointer),
				Code:    "publication.link" + err.Code,
			})
		}
	}

	// type specific validation
	switch d.Type {
	case "dissertation":
		errs = append(errs, d.validateDissertation()...)
	case "journal_article":
		errs = append(errs, d.validateJournalArticle()...)
	case "miscellaneous":
		errs = append(errs, d.validateMiscellaneous()...)
	case "book":
		errs = append(errs, d.validateBook()...)
	case "book_chapter":
		errs = append(errs, d.validateBookChapter()...)
	case "conference":
		errs = append(errs, d.validateConference()...)
	case "book_editor":
		errs = append(errs, d.validateBookEditor()...)
	case "issue_editor":
		errs = append(errs, d.validateIssueEditor()...)
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
			Pointer: "/file_id",
			Code:    "file_id.required",
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
