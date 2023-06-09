package publication

import (
	"fmt"
	"io"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/excel"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
)

const sep = " ; "

var headers = []string{
	"id",
	"status",
	"creator",
	"user",
	"date_created",
	"date_updated",
	"title",
	"access_level",
	"author",
	"ugent_author",
	"contributor",
	"ugent_contributor",
	"department",
	"doi",
	"embargo_date",
	"access_level_after_embargo",
	"format",
	"keyword",
	"license",
	"other_license",
	"locked",
	"project_id",
	"project_title",
	"publisher",
	"reviewer_tags",
	"year",
}

type xlsx struct {
	excel.BaseExporter
}

func NewExporter(writer io.Writer) backends.DatasetListExporter {
	baseExporter := excel.NewBaseExporter(writer)
	baseExporter.Headers = headers
	return &xlsx{
		BaseExporter: *baseExporter,
	}
}

func (x *xlsx) Add(dataset *models.Dataset) {
	x.BaseExporter.Add(x.datasetToRow(dataset))
}

func (x *xlsx) datasetToRow(d *models.Dataset) []string {
	//see also: biblio/lib/Catmandu/Fix/publication_to_csv.pm

	m := map[string]string{}
	m["id"] = d.ID
	m["status"] = d.Status
	if d.Creator != nil {
		m["creator"] = fmt.Sprintf("%s (%s)", d.Creator.Name, d.Creator.ID)
	}
	if d.User != nil {
		m["user"] = fmt.Sprintf("%s (%s)", d.User.Name, d.User.ID)
	}
	if d.DateCreated != nil {
		m["date_created"] = internal_time.FormatTimeUTC(d.DateCreated)
	}
	if d.DateUpdated != nil {
		m["date_updated"] = internal_time.FormatTimeUTC(d.DateUpdated)
	}
	m["title"] = d.Title
	m["access_level"] = d.AccessLevel

	//field: <role>
	//field: ugent_<role>
	for _, role := range []string{"author", "contributor"} {
		contributors := d.Contributors(role)
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
	if len(d.RelatedOrganizations) > 0 {
		depIds := make([]string, 0, len(d.RelatedOrganizations))
		for _, dep := range d.RelatedOrganizations {
			depIds = append(depIds, dep.OrganizationID)
		}
		m["department"] = strings.Join(depIds, sep)
	}

	m["doi"] = d.Identifiers.Get("DOI")
	m["embargo_date"] = d.EmbargoDate
	m["access_level_after_embargo"] = d.AccessLevelAfterEmbargo
	m["format"] = strings.Join(d.Format, sep)
	m["keyword"] = strings.Join(d.Keyword, sep)
	m["license"] = d.License
	m["other_license"] = d.OtherLicense
	m["locked"] = fmt.Sprintf("%t", d.Locked)
	//TODO: projects
	projectIds := make([]string, 0, len(d.RelatedProjects))
	projectTitles := make([]string, 0, len(d.RelatedProjects))
	for _, rel := range d.RelatedProjects {
		projectIds = append(projectIds, rel.ProjectID)
		projectTitles = append(projectTitles, rel.Project.Title)
	}
	m["project_id"] = strings.Join(projectIds, sep)
	m["project_title"] = strings.Join(projectTitles, sep)
	m["publisher"] = d.Publisher
	m["reviewer_tags"] = strings.Join(d.ReviewerTags, sep)
	m["year"] = d.Year

	//hash to ordered list
	row := make([]string, 0, len(headers))
	for _, h := range headers {
		row = append(row, m[h])
	}

	return row
}
