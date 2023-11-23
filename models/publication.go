package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"slices"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/pagination"
	"github.com/ugent-library/biblio-backoffice/util"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
)

type PublicationHits struct {
	pagination.Pagination
	Hits   []*Publication         `json:"hits"`
	Facets map[string]FacetValues `json:"facets"`
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
}

type PublicationLink struct {
	ID          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Description string `json:"description,omitempty"`
}

type PublicationORCIDWork struct {
	ORCID   string `json:"orcid,omitempty"`
	PutCode int    `json:"put_code,omitempty"`
}

type RelatedDataset struct {
	ID string `json:"id,omitempty"`
}

type Publication struct {
	Abstract         []*Text        `json:"abstract,omitempty"`
	AdditionalInfo   string         `json:"additional_info,omitempty"`
	AlternativeTitle []string       `json:"alternative_title,omitempty"`
	ArticleNumber    string         `json:"article_number,omitempty"`
	ArxivID          string         `json:"arxiv_id,omitempty"`
	Author           []*Contributor `json:"author,omitempty"`
	BatchID          string         `json:"batch_id,omitempty"`
	Classification   string         `json:"classification,omitempty"`
	// CompletenessScore       int                     `json:"completeness_score,omitempty"`
	ConferenceName      string     `json:"conference_name,omitempty"`
	ConferenceLocation  string     `json:"conference_location,omitempty"`
	ConferenceOrganizer string     `json:"conference_organizer,omitempty"`
	ConferenceStartDate string     `json:"conference_start_date,omitempty"`
	ConferenceEndDate   string     `json:"conference_end_date,omitempty"`
	ConferenceType      string     `json:"conference_type,omitempty"`
	CreatorID           string     `json:"creator_id,omitempty"`
	Creator             *Person    `json:"-"`
	DateCreated         *time.Time `json:"date_created,omitempty"`
	DateUpdated         *time.Time `json:"date_updated,omitempty"`
	DateFrom            *time.Time `json:"date_from,omitempty"`
	DateUntil           *time.Time `json:"date_until,omitempty"`
	DefenseDate         string     `json:"defense_date,omitempty"`
	DefensePlace        string     `json:"defense_place,omitempty"`
	// DefenseTime is deprecated, see https://github.com/ugent-library/biblio-backoffice/issues/1058
	DefenseTime             string                 `json:"defense_time,omitempty"`
	DOI                     string                 `json:"doi,omitempty"`
	Edition                 string                 `json:"edition,omitempty"`
	Editor                  []*Contributor         `json:"editor,omitempty"`
	EISBN                   []string               `json:"eisbn,omitempty"`
	EISSN                   []string               `json:"eissn,omitempty"`
	ESCIID                  string                 `json:"esci_id,omitempty"`
	Extern                  bool                   `json:"extern"`
	ExternalFields          Values                 `json:"external_fields,omitempty"`
	File                    []*PublicationFile     `json:"file,omitempty"`
	Handle                  string                 `json:"handle,omitempty"`
	HasConfidentialData     string                 `json:"has_confidential_data,omitempty"`
	HasPatentApplication    string                 `json:"has_patent_application,omitempty"`
	HasPublicationsPlanned  string                 `json:"has_publications_planned,omitempty"`
	HasPublishedMaterial    string                 `json:"has_published_material,omitempty"`
	ID                      string                 `json:"id,omitempty"`
	ISBN                    []string               `json:"isbn,omitempty"`
	ISSN                    []string               `json:"issn,omitempty"`
	Issue                   string                 `json:"issue,omitempty"`
	IssueTitle              string                 `json:"issue_title,omitempty"`
	JournalArticleType      string                 `json:"journal_article_type,omitempty"`
	Keyword                 []string               `json:"keyword,omitempty"`
	Language                []string               `json:"language,omitempty"`
	LastUserID              string                 `json:"last_user_id,omitempty"`
	LastUser                *Person                `json:"-"`
	LaySummary              []*Text                `json:"lay_summary,omitempty"`
	Legacy                  bool                   `json:"legacy"`
	Link                    []*PublicationLink     `json:"link,omitempty"`
	Locked                  bool                   `json:"locked"`
	Message                 string                 `json:"message,omitempty"`
	MiscellaneousType       string                 `json:"miscellaneous_type,omitempty"`
	ORCIDWork               []PublicationORCIDWork `json:"orcid_work,omitempty"`
	PageCount               string                 `json:"page_count,omitempty"`
	PageFirst               string                 `json:"page_first,omitempty"`
	PageLast                string                 `json:"page_last,omitempty"`
	PlaceOfPublication      string                 `json:"place_of_publication,omitempty"`
	Publication             string                 `json:"publication,omitempty"`
	PublicationAbbreviation string                 `json:"publication_abbreviation,omitempty"`
	PublicationStatus       string                 `json:"publication_status,omitempty"`
	Publisher               string                 `json:"publisher,omitempty"`
	PubMedID                string                 `json:"pubmed_id,omitempty"`
	RelatedDataset          []RelatedDataset       `json:"related_dataset,omitempty"`
	RelatedOrganizations    []*RelatedOrganization `json:"related_organizations,omitempty"`
	RelatedProjects         []*RelatedProject      `json:"related_projects,omitempty"`
	ReportNumber            string                 `json:"report_number,omitempty"`
	ResearchField           []string               `json:"research_field,omitempty"`
	ReviewerNote            string                 `json:"reviewer_note,omitempty"`
	ReviewerTags            []string               `json:"reviewer_tags,omitempty"`
	SeriesTitle             string                 `json:"series_title,omitempty"`
	SnapshotID              string                 `json:"snapshot_id,omitempty"`
	SourceDB                string                 `json:"source_db,omitempty"`
	SourceID                string                 `json:"source_id,omitempty"`
	SourceRecord            string                 `json:"source_record,omitempty"`
	Status                  string                 `json:"status,omitempty"`
	Supervisor              []*Contributor         `json:"supervisor,omitempty"`
	Title                   string                 `json:"title,omitempty"`
	Type                    string                 `json:"type,omitempty"`
	UserID                  string                 `json:"user_id,omitempty"`
	User                    *Person                `json:"-"`
	Volume                  string                 `json:"volume,omitempty"`
	VABBType                string                 `json:"vabb_type,omitempty"`
	VABBID                  string                 `json:"vabb_id,omitempty"`
	VABBApproved            bool                   `json:"vabb_approved"`
	VABBYear                []string               `json:"vabb_year,omitempty"`
	HasBeenPublic           bool                   `json:"has_been_public"`
	WOSID                   string                 `json:"wos_id,omitempty"`
	WOSType                 string                 `json:"wos_type,omitempty"`
	Year                    string                 `json:"year,omitempty"`
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

func (p *Publication) FileIndex(id string) int {
	for i, file := range p.File {
		if file.ID == id {
			return i
		}
	}
	return 0
}

func (p *Publication) SetFile(f *PublicationFile) {
	for i, file := range p.File {
		if file.ID == f.ID {
			p.File[i] = f
		}
	}
}

// NOTE this assumes publication_file_access_levels are ordered from most to
// least accessible
func (p *Publication) MainFile() *PublicationFile {
	for _, a := range vocabularies.Map["publication_file_access_levels"] {
		for _, f := range p.File {
			if f.Relation == "main_file" && f.AccessLevel == a {
				return f
			}
		}
	}

	return nil
}

// format: c:vabb:419551, VABB-1, not approved, 2017
func (p *Publication) VABB() string {
	var parts []string

	if p.VABBID != "" {
		parts = append(parts, p.VABBID)

		if p.VABBType != "" {
			parts = append(parts, p.VABBType)
		}
		if p.VABBApproved {
			parts = append(parts, "approved")
		} else {
			parts = append(parts, "not approved")
		}
	}

	if len(p.VABBYear) > 0 {
		parts = append(parts, p.VABBYear...)
	}

	return strings.Join(parts, ", ")
}

// Citation
func (p *Publication) SummaryParts() []string {

	tmpParts := make([]string, 0)

	if p.Year != "" {
		tmpParts = append(tmpParts, p.Year)
	}

	if p.PublicationAbbreviation != "" {
		tmpParts = append(tmpParts, p.PublicationAbbreviation)
	} else if p.Publication != "" {
		tmpParts = append(tmpParts, p.Publication)
	}

	if p.Publisher != "" {
		tmpParts = append(tmpParts, p.Publisher)
	}

	if p.Volume != "" {
		tmpParts = append(tmpParts, p.Volume)
	}

	if p.Issue != "" {
		tmpParts = append(tmpParts, fmt.Sprintf("(%s)", p.Issue))
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

		tmpParts = append(tmpParts, fmt.Sprintf("%s%s%s", fp, delim, lp))
	}

	reTrimDot := regexp.MustCompile(`^[ \.]+`)
	reTrimSpaceStart := regexp.MustCompile(`^\s*`)
	reTrimSpaceEnd := regexp.MustCompile(`\s*$`)
	reMultiDot := regexp.MustCompile(`\.+`)
	reNonAlpha := regexp.MustCompile(`^\W*$`)

	summaryParts := make([]string, 0, len(tmpParts))

	for _, v := range tmpParts {
		v = reTrimDot.ReplaceAllString(v, "")
		v = reTrimSpaceStart.ReplaceAllString(v, "")
		v = reTrimSpaceEnd.ReplaceAllString(v, "")
		v = reMultiDot.ReplaceAllString(v, ".")
		if reNonAlpha.MatchString(v) {
			v = ""
		}
		if v != "" {
			summaryParts = append(summaryParts, v)
		}
	}

	return summaryParts
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
			return pl
		}
	}
	return nil
}

func (p *Publication) SetLink(l *PublicationLink) {
	for i, link := range p.Link {
		if link.ID == l.ID {
			p.Link[i] = l
		}
	}
}

func (p *Publication) AddLink(l *PublicationLink) {
	l.ID = ulid.Make().String()
	p.Link = append(p.Link, l)
}

func (p *Publication) RemoveLink(id string) {
	links := make([]*PublicationLink, 0)
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
			return abstract
		}
	}
	return nil
}

func (p *Publication) SetAbstract(t *Text) {
	for i, abstract := range p.Abstract {
		if abstract.ID == t.ID {
			p.Abstract[i] = t
		}
	}
}

func (p *Publication) AddAbstract(t *Text) {
	t.ID = ulid.Make().String()
	p.Abstract = append(p.Abstract, t)
}

func (p *Publication) RemoveAbstract(id string) {
	abstracts := make([]*Text, 0)
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
			return ls
		}
	}
	return nil
}

func (p *Publication) SetLaySummary(ls *Text) {
	for i, laySummary := range p.LaySummary {
		if laySummary.ID == ls.ID {
			p.LaySummary[i] = ls
		}
	}
}

func (p *Publication) AddLaySummary(t *Text) {
	t.ID = ulid.Make().String()
	p.LaySummary = append(p.LaySummary, t)
}

func (p *Publication) RemoveLaySummary(id string) {
	lay_summaries := make([]*Text, 0)
	for _, ls := range p.LaySummary {
		if ls.ID != id {
			lay_summaries = append(lay_summaries, ls)
		}
	}
	p.LaySummary = lay_summaries
}

func (p *Publication) AddProject(project *Project) {
	p.RemoveProject(project.ID)
	p.RelatedProjects = append(p.RelatedProjects, &RelatedProject{
		ProjectID: project.ID,
		Project:   project,
	})
}

func (p *Publication) RemoveProject(id string) {
	rels := make([]*RelatedProject, 0)
	for _, rel := range p.RelatedProjects {
		if rel.ProjectID != id {
			rels = append(rels, rel)
		}
	}
	p.RelatedProjects = rels
}

func (p *Publication) AddOrganization(org *Organization) {
	p.RemoveOrganization(org.ID)
	p.RelatedOrganizations = append(p.RelatedOrganizations, &RelatedOrganization{
		OrganizationID: org.ID,
		Organization:   org,
	})
}

func (p *Publication) RemoveOrganization(id string) {
	rels := make([]*RelatedOrganization, 0)
	for _, rel := range p.RelatedOrganizations {
		if rel.OrganizationID != id {
			rels = append(rels, rel)
		}
	}
	p.RelatedOrganizations = rels
}

func (p *Publication) AddFile(file *PublicationFile) {
	file.ID = ulid.Make().String()
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
	case "journal_article", "miscellaneous":
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
	case "journal_article", "miscellaneous":
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
	case "journal_article", "miscellaneous":
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
	errs := okay.NewErrors()

	if p.ID == "" {
		errs.Add(okay.NewError("/id", "publication.id.required"))
	}
	if p.Type == "" {
		errs.Add(okay.NewError("/type", "publication.type.required"))
	} else if !util.IsPublicationType(p.Type) {
		errs.Add(okay.NewError("/type", "publication.type.invalid"))
	}
	// TODO check classification validity
	if p.Classification == "" {
		errs.Add(okay.NewError("/classification", "publication.classification.required"))
	} else if !slices.Contains(p.ClassificationChoices(), p.Classification) {
		errs.Add(okay.NewError("/classification", "publication.classification.invalid"))
	}
	if p.Status == "" {
		errs.Add(okay.NewError("/status", "publication.status.required"))
	} else if !util.IsStatus(p.Status) {
		errs.Add(okay.NewError("/status", "publication.status.invalid"))
	}

	if p.Status == "public" && p.Title == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/title",
			Rule: "publication.title.required",
		})
	}

	if p.Status == "public" && !p.Legacy && p.Year == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/year",
			Rule: "publication.year.required",
		})
	}
	if p.Year != "" && !util.IsYear(p.Year) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/year",
			Rule: "publication.year.invalid",
		})
	}

	for i, l := range p.Language {
		if !slices.Contains(vocabularies.Map["language_codes"], l) {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  fmt.Sprintf("/language/%d", i),
				Rule: "publication.lang.invalid",
			})
		}
	}

	for i, a := range p.Abstract {
		var e *okay.Errors
		if errors.As(a.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/abstract/%d%s", i, err.Key),
					Rule: "publication.abstract." + err.Rule,
				})
			}
		}
	}

	for i, l := range p.LaySummary {
		var e *okay.Errors
		if errors.As(l.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/lay_summary/%d%s", i, err.Key),
					Rule: "publication.lay_summary." + err.Rule,
				})
			}
		}
	}

	if p.Status == "public" && p.UsesAuthor() && len(p.Author) == 0 {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/author",
			Rule: "publication.author.required",
		})

	}

	// at least one ugent author if not external
	if p.Status == "public" && !p.Legacy && p.UsesAuthor() && !p.Extern {
		var hasUgentAuthors bool = false
		for _, a := range p.Author {
			if a.PersonID != "" {
				hasUgentAuthors = true
				break
			}
		}
		if !hasUgentAuthors {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/author",
				Rule: "publication.author.min_ugent_authors",
			})
		}
	}

	if p.Status == "public" && p.UsesEditor() && !p.UsesAuthor() && len(p.Editor) == 0 {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/editor",
			Rule: "publication.editor.required",
		})
	}

	// at least one ugent editor for editor types if not external
	if p.Status == "public" && !p.Legacy && p.UsesEditor() && !p.UsesAuthor() && !p.Extern {
		var hasUgentEditors bool = false
		for _, a := range p.Editor {
			if a.PersonID != "" {
				hasUgentEditors = true
				break
			}
		}
		if !hasUgentEditors {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/editor",
				Rule: "publication.editor.min_ugent_editors",
			})
		}
	}

	if p.Status == "public" && !p.Legacy && p.UsesSupervisor() && len(p.Supervisor) == 0 {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/supervisor",
			Rule: "publication.supervisor.required",
		})
	}

	for i, c := range p.Author {
		var e *okay.Errors
		if errors.As(c.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/author/%d%s", i, err.Key),
					Rule: "publication.author." + err.Rule,
				})
			}
		}
	}
	for i, c := range p.Editor {
		var e *okay.Errors
		if errors.As(c.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/editor/%d%s", i, err.Key),
					Rule: "publication.editor." + err.Rule,
				})
			}
		}
	}
	for i, c := range p.Supervisor {
		var e *okay.Errors
		if errors.As(c.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/supervisor/%d%s", i, err.Key),
					Rule: "publication.supervisor." + err.Rule,
				})
			}
		}
	}

	for i, rel := range p.RelatedProjects {
		if rel.ProjectID == "" {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  fmt.Sprintf("/project/%d/id", i),
				Rule: "publication.project.id.required",
			})
		}
	}

	for i, rel := range p.RelatedOrganizations {
		if rel.OrganizationID == "" {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  fmt.Sprintf("/department/%d/id", i),
				Rule: "publication.department.id.required",
			})
		}
	}

	for i, rel := range p.RelatedDataset {
		if rel.ID == "" {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  fmt.Sprintf("/related_dataset/%d/id", i),
				Rule: "publication.related_dataset.id.required",
			})
		}
	}

	for i, f := range p.File {
		var e *okay.Errors
		if errors.As(f.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/file/%d%s", i, err.Key),
					Rule: "publication.file." + err.Rule,
				})
			}
		}
	}

	for i, l := range p.Link {
		var e *okay.Errors
		if errors.As(l.Validate(), &e) {
			for _, err := range e.Errors {
				errs.Errors = append(errs.Errors, &okay.Error{
					Key:  fmt.Sprintf("/link/%d%s", i, err.Key),
					Rule: "publication.link." + err.Rule,
				})
			}
		}
	}

	// type specific validation
	switch p.Type {
	case "dissertation":
		okay.Add(errs, p.validateDissertation())
	case "journal_article":
		okay.Add(errs, p.validateJournalArticle())
	case "miscellaneous":
		okay.Add(errs, p.validateMiscellaneous())
	case "book":
		okay.Add(errs, p.validateBook())
	case "book_chapter":
		okay.Add(errs, p.validateBookChapter())
	case "conference":
		okay.Add(errs, p.validateConference())
	case "book_editor":
		okay.Add(errs, p.validateBookEditor())
	case "issue_editor":
		okay.Add(errs, p.validateIssueEditor())
	}

	return errs.ErrorOrNil()
}

func (p *Publication) validateBookEditor() error {
	return nil
}

func (p *Publication) validateIssueEditor() error {
	return nil
}

func (p *Publication) validateJournalArticle() error {
	errs := okay.NewErrors()

	// TODO: confusing: gui shows select without empty element
	// but first creation sets this value to empty
	if p.JournalArticleType != "" && !slices.Contains(vocabularies.Map["journal_article_types"], p.JournalArticleType) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/journal_article_type",
			Rule: "publication.journal_article_type.invalid",
		})
	}
	if p.Status == "public " && !p.Legacy && p.Publication == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/publication",
			Rule: "publication.journal_article.publication.required",
		})
	}

	return errs.ErrorOrNil()
}

func (p *Publication) validateBook() error {
	return nil
}

func (p *Publication) validateConference() error {
	return nil
}

func (p *Publication) validateBookChapter() error {
	return nil
}

func (p *Publication) validateDissertation() error {
	errs := okay.NewErrors()

	if p.Status == "public" && !p.Legacy && p.DefensePlace == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/defense_place",
			Rule: "publication.defense_place.required",
		})
	}
	if p.Status == "public" && !p.Legacy && p.DefenseDate == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/defense_date",
			Rule: "publication.defense_date.required",
		})
	}
	if p.DefenseDate != "" && !util.IsDate(p.DefenseDate) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/defense_date",
			Rule: "publication.defense_date.invalid",
		})
	}

	return errs.ErrorOrNil()
}

func (p *Publication) validateMiscellaneous() error {
	errs := okay.NewErrors()

	// TODO confusing: gui shows select without empty element
	// but first creation sets this value to empty
	if p.MiscellaneousType != "" && !slices.Contains(vocabularies.Map["miscellaneous_types"], p.MiscellaneousType) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/miscellaneous_type",
			Rule: "publication.miscellaneous_type.invalid",
		})
	}

	return errs.ErrorOrNil()
}

func (pf *PublicationFile) Validate() error {
	errs := okay.NewErrors()

	if !slices.Contains(vocabularies.Map["publication_file_access_levels"], pf.AccessLevel) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/access_level",
			Rule: "access_level.invalid",
		})
	}

	if pf.Size == 0 {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/size",
			Rule: "size.empty",
		})
	}

	if pf.ContentType == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/content_type",
			Rule: "content_type.required",
		})
	}

	if pf.ID == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/id",
			Rule: "id.required",
		})
	}

	if pf.Relation != "" && !slices.Contains(vocabularies.Map["publication_file_relations"], pf.Relation) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/relation",
			Rule: "relation.invalid",
		})
	}

	if pf.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
		if !util.IsDate(pf.EmbargoDate) {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/embargo_date",
				Rule: "embargo_date.invalid",
			})
		}

		invalid := false
		if !slices.Contains(vocabularies.Map["publication_file_access_levels"], pf.AccessLevelAfterEmbargo) {
			invalid = true
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/access_level_after_embargo",
				Rule: "access_level_after_embargo.invalid",
			})
		}

		if !slices.Contains(vocabularies.Map["publication_file_access_levels"], pf.AccessLevelDuringEmbargo) {
			invalid = true
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/access_level_during_embargo",
				Rule: "access_level_during_embargo.invalid",
			})
		}

		if pf.AccessLevelAfterEmbargo == pf.AccessLevelDuringEmbargo && !invalid {
			errs.Errors = append(errs.Errors, &okay.Error{
				Key:  "/access_level_after_embargo",
				Rule: "access_level_after_embargo.similar",
			})
		}
	}

	if pf.PublicationVersion != "" && !slices.Contains(vocabularies.Map["publication_versions"], pf.PublicationVersion) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/publication_version",
			Rule: "publication_version.invalid",
		})
	}

	return errs.ErrorOrNil()
}

func (pl *PublicationLink) Validate() error {
	errs := okay.NewErrors()

	if pl.ID == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/id",
			Rule: "id.required",
		})
	}
	if pl.URL == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/url",
			Rule: "url.required",
		})
	}
	if !slices.Contains(vocabularies.Map["publication_link_relations"], pl.Relation) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/relation",
			Rule: "relation.invalid",
		})
	}

	return errs.ErrorOrNil()
}

func (p *Publication) CleanupUnusedFields() bool {
	changed := false

	if !p.UsesJournalArticleType() && p.JournalArticleType != "" {
		changed = true
		p.JournalArticleType = ""
	}

	if !p.UsesConferenceType() && p.ConferenceType != "" {
		changed = true
		p.ConferenceType = ""
	}

	if !p.UsesMiscellaneousType() && p.MiscellaneousType != "" {
		changed = true
		p.MiscellaneousType = ""
	}

	if !slices.Contains(p.ClassificationChoices(), p.Classification) {
		changed = true
		p.Classification = "U"
	}

	if !p.UsesAbstract() && p.Abstract != nil {
		changed = true
		p.Abstract = nil
	}

	if !p.UsesAdditionalInfo() && p.AdditionalInfo != "" {
		changed = true
		p.AdditionalInfo = ""
	}

	if !p.UsesAlternativeTitle() && p.AlternativeTitle != nil {
		changed = true
		p.AlternativeTitle = nil
	}

	if !p.UsesArticleNumber() && p.ArticleNumber != "" {
		changed = true
		p.ArticleNumber = ""
	}

	if !p.UsesArxivID() && p.ArxivID != "" {
		changed = true
		p.ArxivID = ""
	}

	if !p.UsesConference() {
		if p.ConferenceName != "" {
			changed = true
			p.ConferenceName = ""
		}
		if p.ConferenceLocation != "" {
			changed = true
			p.ConferenceLocation = ""
		}
		if p.ConferenceOrganizer != "" {
			changed = true
			p.ConferenceOrganizer = ""
		}
		if p.ConferenceStartDate != "" {
			changed = true
			p.ConferenceStartDate = ""
		}
		if p.ConferenceEndDate != "" {
			changed = true
			p.ConferenceEndDate = ""
		}
	}

	if !p.UsesDefense() {
		if p.DefenseDate != "" {
			changed = true
			p.DefenseDate = ""
		}
		if p.DefensePlace != "" {
			changed = true
			p.DefensePlace = ""
		}
	}

	if !p.UsesDOI() && p.DOI != "" {
		changed = true
		p.DOI = ""
	}

	if !p.UsesEdition() && p.Edition != "" {
		changed = true
		p.Edition = ""
	}

	if !p.UsesISBN() {
		if p.ISBN != nil {
			changed = true
			p.ISBN = nil
		}
		if p.EISBN != nil {
			changed = true
			p.EISBN = nil
		}
	}

	if !p.UsesISSN() {
		if p.ISSN != nil {
			changed = true
			p.ISSN = nil
		}
		if p.EISSN != nil {
			changed = true
			p.EISSN = nil
		}
	}

	if !p.UsesESCIID() && p.ESCIID != "" {
		changed = true
		p.ESCIID = ""
	}

	if !p.UsesConfirmations() {
		if p.HasConfidentialData != "" {
			changed = true
			p.HasConfidentialData = ""
		}
		if p.HasPatentApplication != "" {
			changed = true
			p.HasPatentApplication = ""
		}
		if p.HasPublicationsPlanned != "" {
			changed = true
			p.HasPublicationsPlanned = ""
		}
		if p.HasPublishedMaterial != "" {
			changed = true
			p.HasPublishedMaterial = ""
		}
	}

	if !p.UsesIssue() {
		if p.Issue != "" {
			changed = true
			p.Issue = ""
		}
		if p.IssueTitle != "" {
			changed = true
			p.IssueTitle = ""
		}
	}

	if !p.UsesKeyword() && p.Keyword != nil {
		changed = true
		p.Keyword = nil
	}

	if !p.UsesLanguage() && p.Language != nil {
		changed = true
		p.Language = nil
	}

	if !p.UsesLaySummary() && p.LaySummary != nil {
		changed = true
		p.LaySummary = nil
	}

	if !p.UsesLink() && p.Link != nil {
		changed = true
		p.Link = nil
	}

	if !p.UsesPageCount() && p.PageCount != "" {
		changed = true
		p.PageCount = ""
	}

	if !p.UsesPage() {
		if p.PageFirst != "" {
			changed = true
			p.PageFirst = ""
		}
		if p.PageLast != "" {
			changed = true
			p.PageLast = ""
		}
	}

	if !p.UsesProject() && p.RelatedProjects != nil {
		changed = true
		p.RelatedProjects = nil
	}

	if !p.UsesPublication() && p.Publication != "" {
		changed = true
		p.Publication = ""
	}

	if !p.UsesPublicationAbbreviation() && p.PublicationAbbreviation != "" {
		changed = true
		p.PublicationAbbreviation = ""
	}

	if !p.UsesPublicationStatus() && p.PublicationAbbreviation != "" {
		changed = true
		p.PublicationStatus = ""
	}

	if !p.UsesPubMedID() && p.PubMedID != "" {
		changed = true
		p.PubMedID = ""
	}

	if !p.UsesReportNumber() && p.ReportNumber != "" {
		changed = true
		p.ReportNumber = ""
	}

	if !p.UsesPublisher() {
		if p.Publisher != "" {
			changed = true
			p.Publisher = ""
		}
		if p.PlaceOfPublication != "" {
			changed = true
			p.PlaceOfPublication = ""
		}
	}

	if !p.UsesResearchField() && p.ResearchField != nil {
		changed = true
		p.ResearchField = nil
	}

	if !p.UsesSeriesTitle() && p.SeriesTitle != "" {
		changed = true
		p.SeriesTitle = ""
	}

	if !p.UsesTitle() && p.Title != "" {
		changed = true
		p.Title = ""
	}

	if !p.UsesAuthor() && p.Author != nil {
		changed = true
		p.Author = nil
	}

	if !p.UsesEditor() && p.Editor != nil {
		changed = true
		p.Editor = nil
	}

	if !p.UsesSupervisor() && p.Supervisor != nil {
		changed = true
		p.Supervisor = nil
	}

	return changed
}

func (p *Publication) ChangeType(newType string) {
	usedAuthor := p.UsesAuthor()
	usedEditor := p.UsesEditor()

	p.Type = newType

	// transfer authors to Editor
	if !p.UsesAuthor() && usedAuthor && p.UsesEditor() && p.Editor == nil {
		p.Editor = p.Author
		for _, c := range p.Editor {
			c.CreditRole = nil
		}
	}

	// transfer editors to Author
	if !p.UsesEditor() && usedEditor && p.UsesAuthor() && p.Author == nil {
		p.Author = p.Editor
	}

	p.CleanupUnusedFields()
}

func (pf *PublicationFile) ClearEmbargo() {
	pf.AccessLevel = pf.AccessLevelAfterEmbargo
	pf.AccessLevelDuringEmbargo = ""
	pf.AccessLevelAfterEmbargo = ""
	pf.EmbargoDate = ""
}
