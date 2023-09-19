package publication

import (
	"fmt"
	"io"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/excel"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
	"github.com/ugent-library/biblio-backoffice/models"
)

const sep = " ; "

var headers = []string{
	"id",
	"type",
	"status",
	"author",
	"ugent_author",
	"editor",
	"ugent_editor",
	"supervisor",
	"ugent_supervisor",
	"department",
	"title",
	"alternative_title",
	"publisher",
	"place_of_publication",
	//"reference",
	"page_count",
	"page_first",
	"page_last",
	"isbn",
	"issn",
	"wos_id",
	"wos_type",
	"classification",
	"year",
	//"abstract",
	"research_field", //TODO
	"keyword",
	"publication",
	"publication_abbreviation",
	"series_title",
	"volume",
	"issue",
	"issue_title",
	"edition",
	"conference_name",
	"conference_organizer",
	"conference_start_date",
	"conference_end_date",
	"conference_location",
	"defense_date",
	"defense_location",
	"journal_article_type",
	"miscellaneous_type",
	"conference_type",
	"license",
	"other_license",
	"publication_status",
	"language",
	"doi",
	"date_created",
	"date_updated",
	//"date_submitted",
	//"date_approved",
	"external",
	"has_file",
	"has_open_access_file",
	"jcr_eigenfactor",
	"jcr_immediacy_index",
	"jcr_impact_factor",
	"jcr_impact_factor_5year",
	"jcr_total_cites",
	"jcr_category_quartile",
	"jcr_prev_impact_factor",
	"jcr_prev_category_quartile",
	"project_id",
	"project_title",
	//"project_abstract",
	"project_start_date",
	"project_end_date",
	"project_eu_id",
	"project_eu_call_id",
	"project_eu_acronym",
	"vabb_id",
	"vabb_type",
	"vabb_approved",
	"vabb_year",
	"article_number",
	"jcr_category",
	"jcr_category_rank",
	"jcr_category_decile",
	"jcr_prev_category_decile",
	"jcr_category_vigintile",
	"jcr_prev_category_vigintile",
	//"project_url",
	"esci_id",
	"locked",
}

type xlsx struct {
	excel.BaseExporter
}

func NewExporter(writer io.Writer) backends.PublicationListExporter {
	baseExporter := excel.NewBaseExporter(writer)
	baseExporter.Headers = headers
	return &xlsx{
		BaseExporter: *baseExporter,
	}
}

func (x *xlsx) Add(pub *models.Publication) {
	x.BaseExporter.Add(x.publicationToRow(pub))
}

func (x *xlsx) publicationToRow(pub *models.Publication) []string {
	//see also: biblio/lib/Catmandu/Fix/publication_to_csv.pm
	//see also: librecat/ugent/config/route.yml
	//see also: librecat/ugent/fixes/to_reviewer_xlsx.fix

	m := map[string]string{}
	m["id"] = pub.ID
	m["type"] = pub.Type
	m["status"] = pub.Status

	//field: <role>
	//field: ugent_<role>
	for _, role := range []string{"author", "editor", "supervisor"} {
		contributors := pub.Contributors(role)
		{
			values := []string{}
			for _, c := range contributors {
				values = append(values, c.Name())
			}
			m[role] = strings.Join(values, sep)
		}
		{
			values := []string{}
			for _, c := range contributors {
				if c.Person == nil || len(c.Person.UGentID) == 0 {
					continue
				}
				group := ""
				if len(c.Person.Affiliations) > 0 {
					group = "@" + c.Person.Affiliations[0].OrganizationID
				}
				//full_name (<ugent_id>)
				//full_name (<ugent_id>@<department.0.id>)
				val := fmt.Sprintf("%s (%s%s)", c.Name(), c.Person.UGentID[0], group)
				values = append(values, val)
			}
			m["ugent_"+role] = strings.Join(values, sep)
		}
	}

	//field: department
	if len(pub.RelatedOrganizations) > 0 {
		depIds := make([]string, 0, len(pub.RelatedOrganizations))
		for _, dep := range pub.RelatedOrganizations {
			//TODO: biblio skips department without id? Do they exist?
			depIds = append(depIds, dep.OrganizationID)
		}
		m["department"] = strings.Join(depIds, sep)
	}

	m["title"] = pub.Title
	m["alternative_title"] = strings.Join(pub.AlternativeTitle, sep)
	m["publisher"] = pub.Publisher
	m["place_of_publication"] = pub.PlaceOfPublication
	m["page_count"] = pub.PageCount
	m["page_first"] = pub.PageFirst
	m["page_last"] = pub.PageLast
	isbns := append([]string{}, pub.ISBN...)
	isbns = append(isbns, pub.EISBN...)
	m["isbn"] = strings.Join(isbns, sep)
	issns := append([]string{}, pub.ISSN...)
	issns = append(issns, pub.EISSN...)
	m["issn"] = strings.Join(issns, sep)
	m["wos_id"] = pub.WOSID
	m["wos_type"] = pub.WOSType
	m["classification"] = pub.Classification
	m["year"] = pub.Year
	m["research_field"] = strings.Join(pub.ResearchField, sep)
	m["keyword"] = strings.Join(pub.Keyword, sep)
	m["publication"] = pub.Publication
	m["publication_abbreviation"] = pub.PublicationAbbreviation
	m["series_title"] = pub.SeriesTitle
	m["volume"] = pub.Volume
	m["issue"] = pub.Issue
	m["issue_title"] = pub.IssueTitle
	m["edition"] = pub.Edition
	m["conference_name"] = pub.ConferenceName
	m["conference_organizer"] = pub.ConferenceOrganizer
	m["conference_start_date"] = pub.ConferenceStartDate
	m["conference_end_date"] = pub.ConferenceEndDate
	m["conference_location"] = pub.ConferenceLocation
	m["defense_date"] = pub.DefenseDate
	m["defense_location"] = pub.DefensePlace
	m["journal_article_type"] = pub.JournalArticleType
	m["miscellaneous_type"] = pub.MiscellaneousType
	m["conference_type"] = pub.ConferenceType

	if len(pub.File) > 0 {
		m["license"] = pub.File[0].License
		m["other_license"] = pub.File[0].OtherLicense
	}

	m["publication_status"] = pub.PublicationStatus
	m["language"] = strings.Join(pub.Language, sep)
	m["doi"] = pub.DOI
	m["date_created"] = internal_time.FormatTimeUTC(pub.DateCreated)
	m["date_updated"] = internal_time.FormatTimeUTC(pub.DateUpdated)
	m["external"] = fmt.Sprintf("%t", pub.Extern)
	m["has_file"] = fmt.Sprintf("%t", len(pub.File) > 0)
	m["has_open_access_file"] = "false"
	for _, file := range pub.File {
		if validation.InArray(vocabularies.Map["publication_file_access_levels"], file.AccessLevel) {
			m["has_open_access_file"] = "true"
			break
		}
	}

	//TODO: jcr fields
	m["jcr_eigenfactor"] = ""
	m["jcr_immediacy_index"] = ""
	m["jcr_impact_factor"] = ""
	m["jcr_impact_factor_5year"] = ""
	m["jcr_total_cites"] = ""
	m["jcr_category_quartile"] = ""
	m["jcr_prev_impact_factor"] = ""
	m["jcr_prev_category_quartile"] = ""

	//TODO: projects
	projectIds := make([]string, 0, len(pub.RelatedProjects))
	projectTitles := make([]string, 0, len(pub.RelatedProjects))
	for _, rel := range pub.RelatedProjects {
		projectIds = append(projectIds, rel.ProjectID)
		projectTitles = append(projectTitles, rel.Project.Title)
	}
	m["project_id"] = strings.Join(projectIds, sep)
	m["project_title"] = strings.Join(projectTitles, sep)
	m["project_abstract"] = ""
	m["project_start_date"] = ""
	m["project_end_date"] = ""
	m["project_eu_id"] = ""
	m["project_eu_call_id"] = ""
	m["project_eu_acronym"] = ""

	m["vabb_id"] = pub.VABBID
	m["vabb_type"] = pub.VABBType
	m["vabb_approved"] = fmt.Sprintf("%t", pub.VABBApproved)
	m["vabb_year"] = strings.Join(pub.VABBYear, sep)
	m["article_number"] = pub.ArticleNumber

	//TODO: jcr fields
	m["jcr_category"] = ""
	m["jcr_category_rank"] = ""
	m["jcr_category_decile"] = ""
	m["jcr_prev_category_decile"] = ""
	m["jcr_category_vigintile"] = ""
	m["jcr_prev_category_vigintile"] = ""

	m["esci_id"] = pub.ESCIID
	m["locked"] = fmt.Sprintf("%t", pub.Locked)

	//hash to ordered list
	row := make([]string, 0, len(headers))
	for _, h := range headers {
		row = append(row, m[h])
	}

	return row
}
