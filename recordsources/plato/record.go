package plato

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/recordsources"
)

type platoRecord struct {
	id   string
	data []byte
}

func NewRecord(id string, data []byte) *platoRecord {
	return &platoRecord{
		id:   id,
		data: data,
	}
}

func (r *platoRecord) SourceName() string {
	return "plato"
}

func (r *platoRecord) SourceID() string {
	return r.id
}

func (r *platoRecord) ToCandidateRecord(services *backends.Services) (*models.CandidateRecord, error) {
	p := &models.Publication{}
	md := gjson.ParseBytes(r.data)
	platoID := md.Get("plato_id").String()

	p.Type = "dissertation"
	p.Status = "private"
	p.Classification = "U"
	p.SourceDB = "plato"
	p.SourceID = platoID

	if v := md.Get("titel.eng"); v.Exists() {
		p.Title = v.String()
		if v := md.Get("titel.ned"); v.Exists() {
			p.AlternativeTitle = append(p.AlternativeTitle, v.String())
		}
	} else if v := md.Get("titel.ned"); v.Exists() {
		p.Title = v.String()
	}
	if v := md.Get("year"); v.Exists() {
		p.Year = v.String()
	}
	if v := md.Get("defence.date"); v.Exists() {
		p.DefenseDate = v.String()
	}
	p.DefensePlace = "Ghent, Belgium" // TODO

	ugentID := md.Get("student.ugentid").String()
	if ugentID == "" && md.Get("student.studid").String() != "" {
		ugentID = "0000" + md.Get("student.studid").String()
	}
	if ugentID != "" {
		hits, err := services.PersonSearchService.SuggestPeople(ugentID)
		if err != nil {
			return nil, err
		}
		if len(hits) != 1 {
			return nil, errors.New("multiple or no matches for ugent id " + ugentID)
		}
		c := models.ContributorFromPerson(hits[0])
		p.Author = append(p.Author, c)
	} else {
		c := models.ContributorFromFirstLastName(md.Get("student.first").String(), md.Get("student.last").String())
		c.ExternalPerson.Affiliation = md.Get("student.affil").String()
		c.ExternalPerson.HonorificPrefix = md.Get("student.title").String()
		p.Author = append(p.Author, c)
	}

	var cbErr error
	md.Get("supervisors").ForEach(func(key, val gjson.Result) bool {
		if v := val.Get("ugentid"); v.Exists() {
			hits, err := services.PersonSearchService.SuggestPeople(v.String())
			if err != nil {
				cbErr = err
				return false
			}
			if len(hits) != 1 {
				cbErr = errors.New("multiple or no matches for ugent id " + v.String())
				return false
			}
			p.Supervisor = append(p.Supervisor, models.ContributorFromPerson(hits[0]))
		} else {
			c := models.ContributorFromFirstLastName(val.Get("first").String(), val.Get("last").String())
			c.ExternalPerson.Affiliation = val.Get("affil").String()
			c.ExternalPerson.HonorificPrefix = val.Get("title").String()
			p.Supervisor = append(p.Supervisor, c)
		}
		return true
	})
	if cbErr != nil {
		return nil, cbErr
	}

	if v := md.Get("pdf.ISBN"); v.Exists() {
		p.ISBN = append(p.ISBN, v.String())
	}
	if v := md.Get("pdf.abstract"); v.Exists() {
		p.AddAbstract(&models.Text{Lang: "dut", Text: v.String()})
	}
	if v := md.Get("pdf.url"); v.Exists() {
		sha256, size, err := recordsources.StoreURL(context.TODO(), v.String(), services.FileStore)
		if err != nil {
			return nil, err
		}
		f := &models.PublicationFile{
			Relation:           "main_file",
			Name:               r.id + ".pdf",
			ContentType:        "application/pdf",
			Size:               size,
			SHA256:             sha256,
			PublicationVersion: "publishedVersion",
		}
		embargo := md.Get("pdf.embargo").String()
		access := md.Get("pdf.accesstype").String()
		if strings.HasPrefix(embargo, "9999") {
			f.AccessLevel = "info:eu-repo/semantics/closedAccess"
		} else if embargo != "" {
			f.AccessLevel = "info:eu-repo/semantics/embargoedAccess"
			f.AccessLevelDuringEmbargo = "info:eu-repo/semantics/closedAccess"
			f.EmbargoDate = embargo[:10]
			if access == "U" {
				f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/restrictedAccess"
			} else if access == "W" {
				f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/openAccess"
			}
		} else if access == "U" {
			f.AccessLevel = "info:eu-repo/semantics/restrictedAccess"
		} else if access == "W" {
			f.AccessLevel = "info:eu-repo/semantics/openAccess"
		}
		p.AddFile(f)
	}

	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return &models.CandidateRecord{
		SourceName:     r.SourceName(),
		SourceID:       r.SourceID(),
		SourceMetadata: r.data,
		Type:           "Publication",
		Metadata:       j,
	}, nil

}
