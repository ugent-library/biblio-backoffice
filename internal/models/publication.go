package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

type PublicationUser struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PublicationFile struct {
	AccessLevel              string     `json:"access_level,omitempty"`
	License                  string     `json:"license,omitempty"`
	ContentType              string     `json:"content_type,omitempty"`
	DateCreated              *time.Time `json:"date_created,omitempty"`
	DateUpdated              *time.Time `json:"date_updated,omitempty"`
	EmbargoDate              string     `json:"embargo_date,omitempty"`
	AccessLevelDuringEmbargo string     `json:"access_level_during_embargo,omitempty"`
	AccessLevelAfterEmbargo  string     `json:"access_level_after_embargo,omitempty"`
	Name                     string     `json:"name,omitempty"`
	Size                     int        `json:"size,omitempty"`
	ID                       string     `json:"id,omitempty"`
	SHA256                   string     `json:"sha256,omitempty"`
	OtherLicense             string     `json:"other_license,omitempty"`
	PublicationVersion       string     `json:"publication_version,omitempty"`
	Relation                 string     `json:"relation,omitempty"`
	URL                      string     `json:"url,omitempty"`
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
	Abstract         []Text         `json:"abstract,omitempty"`
	AdditionalInfo   string         `json:"additional_info,omitempty"`
	AlternativeTitle []string       `json:"alternative_title,omitempty"`
	ArticleNumber    string         `json:"article_number,omitempty"`
	ArxivID          string         `json:"arxiv_id,omitempty"`
	Author           []*Contributor `json:"author,omitempty"`
	BatchID          string         `json:"batch_id,omitempty"`
	Classification   string         `json:"classification,omitempty"`
	// CompletenessScore       int                     `json:"completeness_score,omitempty"`
	ConferenceName          string                  `json:"conference_name,omitempty"`
	ConferenceLocation      string                  `json:"conference_location,omitempty"`
	ConferenceOrganizer     string                  `json:"conference_organizer,omitempty"`
	ConferenceStartDate     string                  `json:"conference_start_date,omitempty"`
	ConferenceEndDate       string                  `json:"conference_end_date,omitempty"`
	ConferenceType          string                  `json:"conference_type,omitempty"`
	Creator                 *PublicationUser        `json:"creator,omitempty"`
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
	Extern                  bool                    `json:"extern"`
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
	Legacy                  bool                    `json:"legacy"`
	Link                    []PublicationLink       `json:"link,omitempty"`
	Locked                  bool                    `json:"locked"`
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
	SnapshotID              string                  `json:"snapshot_id,omitempty"`
	SourceDB                string                  `json:"source_db,omitempty"`
	SourceID                string                  `json:"source_id,omitempty"`
	SourceRecord            string                  `json:"source_record,omitempty"`
	Status                  string                  `json:"status,omitempty"`
	Supervisor              []*Contributor          `json:"supervisor,omitempty"`
	Title                   string                  `json:"title,omitempty"`
	Type                    string                  `json:"type,omitempty"`
	User                    *PublicationUser        `json:"user,omitempty"`
	Volume                  string                  `json:"volume,omitempty"`
	VABBType                string                  `json:"vabb_type,omitempty"`
	VABBID                  string                  `json:"vabb_id,omitempty"`
	VABBApproved            bool                    `json:"vabb_approved"`
	VABBYear                []string                `json:"vabb_year,omitempty"`
	HasBeenPublic           bool                    `json:"has_been_public"`
	WOSID                   string                  `json:"wos_id,omitempty"`
	WOSType                 string                  `json:"wos_type,omitempty"`
	Year                    string                  `json:"year,omitempty"`
}

// TODO determine which file passes access level to top
func (p *Publication) AccessLevel() string {
	for _, a := range vocabularies.Map["publication_file_access_levels"] {
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

func (p *Publication) SetFile(f *PublicationFile) {
	for i, file := range p.File {
		if file.ID == f.ID {
			p.File[i] = f
		}
	}
}

// TODO determine which file will be the primary file for thumbnails, etc.
//
//	This isn't necessarily the same file as the primary file used to indicate
//	the access level across the entire publication (e.g. summaries)
func (p *Publication) PrimaryFile() *PublicationFile {
	for _, file := range p.File {
		if file.AccessLevel != "" && file.Relation == "main_file" {
			return file
		}
	}

	return nil
}

// format: c:vabb:419551 (VABB-1, not approved, 2017
func (p *Publication) VABB() string {
	VABBID := "-"
	VABBType := "-"
	VABBApproved := "not approved"
	VABBYear := "-"

	if p.VABBID != "" {
		VABBID = p.VABBID
	}

	if p.VABBType != "" {
		VABBType = p.VABBType
	}

	if p.VABBApproved {
		VABBApproved = "approved"
	}

	if len(p.VABBYear) > 0 {
		VABBYear = strings.Join(p.VABBYear, ", ")
	}

	return fmt.Sprintf("%s (%s, %s, %s)", VABBID, VABBType, VABBApproved, VABBYear)
}

// Citation
func (p *Publication) Reference() string {
	ref := ""

	ref_page := ""
	ref_publisher := ""
	ref_parent := ""
	year := ""
	volume := ""
	issue := ""

	if p.PublicationAbbreviation != "" {
		ref_parent = fmt.Sprintf(" %s ", p.PublicationAbbreviation)
	} else if p.Publication != "" {
		ref_parent = fmt.Sprintf(" %s ", p.Publication)
	}

	if p.Year != "" {
		year = fmt.Sprintf(" %s.", p.Year)
	}

	if p.Publisher != "" {
		ref_publisher = fmt.Sprintf(" %s.", p.Publisher)
	}

	if p.Volume != "" {
		volume = fmt.Sprintf(" %s", p.Volume)
	}

	if p.Issue != "" {
		issue = fmt.Sprintf(" (%s) ", p.Issue)
	}

	if p.PageFirst != "" || p.PageLast != "" {
		fp := ""
		lp := ""
		delim := ""

		if p.PageFirst != "" {
			fp = p.PageFirst
		}

		if p.PageLast != "" {
			lp = p.PageLast
			delim = "-"
		}

		ref_page = fmt.Sprintf(" %s%s%s ", fp, delim, lp)
	}

	ref = fmt.Sprintf("%s%s%s%s%s%s", ref_parent, year, ref_publisher, volume, issue, ref_page)

	var r *regexp.Regexp
	r = regexp.MustCompile(`^[ \.]+`)
	ref = r.ReplaceAllString(ref, "")
	r = regexp.MustCompile(`^\s*`)
	ref = r.ReplaceAllString(ref, "")
	r = regexp.MustCompile(`\s*$`)
	ref = r.ReplaceAllString(ref, "")
	r = regexp.MustCompile(`\.+`)
	ref = r.ReplaceAllString(ref, ".")
	r = regexp.MustCompile(`^\W*$`)

	if r.MatchString(ref) {
		ref = ""
	}

	return ref
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
	case "miscellaneous":
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

func (p *Publication) ResolveDOI() string {
	if p.DOI != "" {
		return "https://doi.org/" + p.DOI
	}
	return ""
}

func (p *Publication) ResolveWOSID() string {
	if p.WOSID != "" {
		return "https://www.webofscience.com/wos/woscc/full-record/" + p.WOSID
	}
	return ""
}

func (p *Publication) ResolvePubMedID() string {
	if p.PubMedID != "" {
		return "https://www.ncbi.nlm.nih.gov/pubmed/" + p.PubMedID
	}
	return ""
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

func (p *Publication) SetLink(l *PublicationLink) {
	for i, link := range p.Link {
		if link.ID == l.ID {
			p.Link[i] = *l
		}
	}
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

func (p *Publication) SetAbstract(t *Text) {
	for i, abstract := range p.Abstract {
		if abstract.ID == t.ID {
			p.Abstract[i] = *t
		}
	}
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

func (p *Publication) SetLaySummary(ls *Text) {
	for i, laySummary := range p.LaySummary {
		if laySummary.ID == ls.ID {
			p.LaySummary[i] = *ls
		}
	}
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
		file.AccessLevel = "info:eu-repo/semantics/restrictedAccess"
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

func (p *Publication) UsesAbstract() bool {
	return true
}

func (p *Publication) UsesAdditionalInfo() bool {
	return true
}

func (p *Publication) UsesAlternativeTitle() bool {
	return true
}

func (p *Publication) UsesArticleNumber() bool {
	switch p.Type {
	case "conference", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesArxivID() bool {
	switch p.Type {
	case "journal_article":
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

func (p *Publication) UsesDefense() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesDOI() bool {
	return true
}

func (p *Publication) UsesEdition() bool {
	switch p.Type {
	case "book_chapter", "book", "book_editor", "issue_editor", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesISBN() bool {
	return true
}

func (p *Publication) UsesISSN() bool {
	return true
}

func (p *Publication) UsesESCIID() bool {
	switch p.Type {
	case "journal_article":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesConfirmations() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesIssue() bool {
	switch p.Type {
	case "conference", "issue_editor", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesKeyword() bool {
	return true
}

func (p *Publication) UsesLanguage() bool {
	return true
}

func (p *Publication) UsesLaySummary() bool {
	switch p.Type {
	case "dissertation":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesLink() bool {
	return false
}

func (p *Publication) UsesPageCount() bool {
	return true
}

func (p *Publication) UsesPage() bool {
	switch p.Type {
	case "book_chapter", "conference", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesProject() bool {
	return true
}

func (p *Publication) UsesPublication() bool {
	switch p.Type {
	case "book_chapter", "conference", "issue_editor", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesPublicationAbbreviation() bool {
	switch p.Type {
	case "conference", "issue_editor", "journal_article", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesPublicationStatus() bool {
	return true
}

func (p *Publication) UsesPublisher() bool {
	return true
}

func (p *Publication) UsesPubMedID() bool {
	switch p.Type {
	case "journal_article":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesReportNumber() bool {
	switch p.Type {
	case "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesResearchField() bool {
	return true
}

func (p *Publication) UsesSeriesTitle() bool {
	switch p.Type {
	case "book_chapter", "book", "book_editor", "conference", "dissertation", "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesTitle() bool {
	return true
}

func (p *Publication) UsesURL() bool {
	return true
}

func (p *Publication) UsesVolume() bool {
	return true
}

func (p *Publication) UsesWOS() bool {
	return true
}

func (p *Publication) UsesYear() bool {
	return true
}

func (p *Publication) UsesJournalArticleType() bool {
	switch p.Type {
	case "journal_article":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesMiscellaneousType() bool {
	switch p.Type {
	case "miscellaneous":
		return true
	default:
		return false
	}
}

func (p *Publication) UsesConferenceType() bool {
	switch p.Type {
	case "conference":
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

func (p *Publication) ShowPublicationAsRequired() bool {
	return p.Type == "journal_article" || p.Type == "book_chapter"
}

func (p *Publication) ShowDefenseAsRequired() bool {
	return p.Type == "dissertation"
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
	} else if !validation.InArray(p.ClassificationChoices(), p.Classification) {
		errs = append(errs, &validation.Error{
			Pointer: "/classification",
			Code:    "publication.classification.invalid",
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

	if p.Status == "public" && !p.Legacy && p.Year == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/year",
			Code:    "publication.year.required",
		})
	}
	if p.Year != "" && !validation.IsYear(p.Year) {
		errs = append(errs, &validation.Error{
			Pointer: "/year",
			Code:    "publication.year.invalid",
		})
	}

	for i, k := range p.Keyword {
		if k == "" {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/keyword/%d", i),
				Code:    "publication.keyword.invalid",
			})
		}
	}

	for i, l := range p.Language {
		if !validation.InArray(vocabularies.Map["language_codes"], l) {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/language/%d", i),
				Code:    "publication.lang.invalid",
			})
		}
	}

	for i, a := range p.Abstract {
		for _, err := range a.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/abstract/%d%s", i, err.Pointer),
				Code:    "publication.abstract." + err.Code,
			})
		}
	}

	for i, l := range p.LaySummary {
		for _, err := range l.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/lay_summary/%d%s", i, err.Pointer),
				Code:    "publication.lay_summary." + err.Code,
			})
		}
	}

	if p.Status == "public" && p.UsesAuthor() && len(p.Author) == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/author",
			Code:    "publication.author.required",
		})

	}

	// at least one ugent author if not external
	if p.Status == "public" && !p.Legacy && p.UsesAuthor() && !p.Extern {
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
	}

	// at least one ugent editor for editor types if not external
	if p.Status == "public" && !p.Legacy && p.UsesEditor() && !p.UsesAuthor() && !p.Extern {
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

	if p.Status == "public" && !p.Legacy && p.UsesSupervisor() && len(p.Supervisor) == 0 {
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
				Code:    "publication.file." + err.Code,
			})
		}
	}

	for i, pl := range p.Link {
		for _, err := range pl.Validate() {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/link/%d%s", i, err.Pointer),
				Code:    "publication.link." + err.Code,
			})
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
	if p.Status == "public " && !p.Legacy && p.Publication == "" {
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
	if p.Status == "public" && !p.Legacy && p.DefensePlace == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_place",
			Code:    "publication.defense_place.required",
		})
	}
	if p.Status == "public" && !p.Legacy && p.DefenseDate == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_date",
			Code:    "publication.defense_date.required",
		})
	}
	if p.DefenseDate != "" && !validation.IsDate(p.DefenseDate) {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_date",
			Code:    "publication.defense_date.invalid",
		})
	}
	if p.Status == "public" && !p.Legacy && p.DefenseTime == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/defense_time",
			Code:    "publication.defense_time.required",
		})
	}
	if p.DefenseTime != "" && !validation.IsTime(p.DefenseTime) {
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

	if pf.Size == 0 {
		errs = append(errs, &validation.Error{
			Pointer: "/size",
			Code:    "size.empty",
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

	if pf.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
		if !validation.IsDate(pf.EmbargoDate) {
			errs = append(errs, &validation.Error{
				Pointer: "/embargo_date",
				Code:    "embargo_date.invalid",
			})
		}

		invalid := false
		if !validation.InArray(vocabularies.Map["publication_file_access_levels"], pf.AccessLevelAfterEmbargo) {
			invalid = true
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_after_embargo",
				Code:    "access_level_after_embargo.invalid",
			})
		}

		if !validation.InArray(vocabularies.Map["publication_file_access_levels"], pf.AccessLevelDuringEmbargo) {
			invalid = true
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_during_embargo",
				Code:    "access_level_during_embargo.invalid",
			})
		}

		if pf.AccessLevelAfterEmbargo == pf.AccessLevelDuringEmbargo && !invalid {
			errs = append(errs, &validation.Error{
				Pointer: "/access_level_after_embargo",
				Code:    "access_level_after_embargo.similar",
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

func (p *Publication) ChangeType(newType string) error {
	usedAuthor := p.UsesAuthor()
	usedEditor := p.UsesEditor()

	p.Type = newType

	if !validation.InArray(p.ClassificationChoices(), p.Classification) {
		p.Classification = "U"
	}

	if !p.UsesAbstract() {
		p.Abstract = nil
	}

	if !p.UsesAdditionalInfo() {
		p.AdditionalInfo = ""
	}

	if !p.UsesAlternativeTitle() {
		p.AlternativeTitle = nil
	}

	if !p.UsesArticleNumber() {
		p.ArticleNumber = ""
	}

	if !p.UsesArxivID() {
		p.ArxivID = ""
	}

	if !p.UsesConference() {
		p.ConferenceName = ""
		p.ConferenceLocation = ""
		p.ConferenceOrganizer = ""
		p.ConferenceStartDate = ""
		p.ConferenceEndDate = ""
	}

	if !p.UsesDefense() {
		p.DefenseDate = ""
		p.DefensePlace = ""
	}

	if !p.UsesDOI() {
		p.DOI = ""
	}

	if !p.UsesEdition() {
		p.Edition = ""
	}

	if !p.UsesISBN() {
		p.ISBN = nil
		p.EISBN = nil
	}

	if !p.UsesISSN() {
		p.ISSN = nil
		p.EISSN = nil
	}

	if !p.UsesESCIID() {
		p.ESCIID = ""
	}

	if !p.UsesConfirmations() {
		p.HasConfidentialData = ""
		p.HasPatentApplication = ""
		p.HasPublicationsPlanned = ""
		p.HasPublishedMaterial = ""
	}

	if !p.UsesIssue() {
		p.Issue = ""
		p.IssueTitle = ""
	}

	if !p.UsesKeyword() {
		p.Keyword = nil
	}

	if !p.UsesLanguage() {
		p.Language = nil
	}

	if !p.UsesLaySummary() {
		p.LaySummary = nil
	}

	if !p.UsesLink() {
		p.Link = nil
	}

	if !p.UsesPageCount() {
		p.PageCount = ""
	}

	if !p.UsesPage() {
		p.PageFirst = ""
		p.PageLast = ""
	}

	if !p.UsesProject() {
		p.Project = nil
	}

	if !p.UsesPublication() {
		p.Publication = ""
	}

	if !p.UsesPublicationAbbreviation() {
		p.PublicationAbbreviation = ""
	}

	if !p.UsesPublicationStatus() {
		p.PublicationStatus = ""
	}

	if !p.UsesPubMedID() {
		p.PubMedID = ""
	}

	if !p.UsesReportNumber() {
		p.ReportNumber = ""
	}

	if !p.UsesPublisher() {
		p.Publisher = ""
		p.PlaceOfPublication = ""
	}

	if !p.UsesResearchField() {
		p.ResearchField = nil
	}

	if !p.UsesSeriesTitle() {
		p.SeriesTitle = ""
	}

	if !p.UsesTitle() {
		p.Title = ""
	}

	if !p.UsesAuthor() {
		if usedAuthor && p.UsesEditor() && p.Editor == nil {
			p.Editor = p.Author
			for _, c := range p.Editor {
				c.CreditRole = nil
			}
		}
		p.Author = nil
	}

	if !p.UsesEditor() {
		if usedEditor && p.UsesAuthor() && p.Author == nil {
			p.Author = p.Editor
		}
		p.Editor = nil
	}

	if !p.UsesSupervisor() {
		p.Supervisor = nil
	}

	return nil
}

func (pf *PublicationFile) ClearEmbargo() {
	pf.AccessLevel = pf.AccessLevelAfterEmbargo
	pf.AccessLevelDuringEmbargo = ""
	pf.AccessLevelAfterEmbargo = ""
	pf.EmbargoDate = ""
}
