package publication

import (
	"fmt"
	"io"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/excel"
	"github.com/ugent-library/biblio-backend/internal/models"
	internal_time "github.com/ugent-library/biblio-backend/internal/time"
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

func (x *xlsx) datasetToRow(ds *models.Dataset) []string {
	//see also: biblio/lib/Catmandu/Fix/publication_to_csv.pm

	m := map[string]string{}
	m["id"] = ds.ID
	m["status"] = ds.Status
	if ds.Creator != nil {
		m["creator"] = fmt.Sprintf("%s (%s)", ds.Creator.Name, ds.Creator.ID)
	}
	if ds.User != nil {
		m["user"] = fmt.Sprintf("%s (%s)", ds.User.Name, ds.User.ID)
	}
	if ds.DateCreated != nil {
		m["date_created"] = internal_time.FormatTimeUTC(ds.DateCreated)
	}
	if ds.DateUpdated != nil {
		m["date_updated"] = internal_time.FormatTimeUTC(ds.DateUpdated)
	}
	m["title"] = ds.Title
	m["access_level"] = ds.AccessLevel

	//field: <role>
	//field: ugent_<role>
	for _, role := range []string{"author", "contributor"} {
		contributors := ds.Contributors(role)
		{
			values := []string{}
			for _, contributor := range contributors {
				fullName := contributor.FullName
				if fullName == "" {
					fullName = contributor.FirstName + " " + contributor.LastName
				}
				values = append(values, fullName)
			}
			m[role] = strings.Join(values, sep)
		}
		{
			values := []string{}
			for _, contributor := range contributors {
				if len(contributor.UGentID) == 0 {
					continue
				}
				group := ""
				if len(contributor.Department) > 0 {
					group = "@" + contributor.Department[0].ID
				}
				//full_name (<ugent_id>)
				//full_name (<ugent_id>@<department.0.id>)
				fullName := contributor.FullName
				if fullName == "" {
					fullName = contributor.FirstName + " " + contributor.LastName
				}
				val := fmt.Sprintf("%s (%s%s)", fullName, contributor.UGentID[0], group)
				values = append(values, val)
			}
			m["ugent_"+role] = strings.Join(values, sep)
		}
	}

	//field: department
	if len(ds.Department) > 0 {
		depIds := make([]string, 0, len(ds.Department))
		for _, dep := range ds.Department {
			depIds = append(depIds, dep.ID)
		}
		m["department"] = strings.Join(depIds, sep)
	}

	m["doi"] = ds.DOI
	m["embargo_date"] = ds.EmbargoDate
	m["access_level_after_embargo"] = ds.AccessLevelAfterEmbargo
	m["format"] = strings.Join(ds.Format, sep)
	m["keyword"] = strings.Join(ds.Keyword, sep)
	m["license"] = ds.License
	m["other_license"] = ds.OtherLicense
	m["locked"] = fmt.Sprintf("%t", ds.Locked)
	//TODO: projects
	projectIds := make([]string, 0, len(ds.Project))
	projectTitles := make([]string, 0, len(ds.Project))
	for _, project := range ds.Project {
		projectIds = append(projectIds, project.ID)
		projectTitles = append(projectTitles, project.Name)
	}
	m["project_id"] = strings.Join(projectIds, sep)
	m["project_title"] = strings.Join(projectTitles, sep)
	m["publisher"] = ds.Publisher
	m["reviewer_tags"] = strings.Join(ds.ReviewerTags, sep)
	m["year"] = ds.Year

	//hash to ordered list
	row := make([]string, 0, len(headers))
	for _, h := range headers {
		row = append(row, m[h])
	}

	return row
}
