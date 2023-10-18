package oaidc

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/ugent-library/biblio-backoffice/frontoffice"
	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

const startTag = `<oai_dc:dc xmlns="http://www.openarchives.org/OAI/2.0/oai_dc/"
xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/"
xmlns:dc="http://purl.org/dc/elements/1.1/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd">
`

const endTag = `
</oai_dc:dc>
`

func writeField(b *bytes.Buffer, tag, val string) {
	if val != "" {
		b.WriteString("<dc:")
		b.WriteString(tag)
		b.WriteString(">")
		xml.EscapeText(b, []byte(val))
		b.WriteString("</dc:")
		b.WriteString(tag)
		b.WriteString(">")
	}
}

type Encoder struct {
	repo    *repositories.Repo
	baseURL string
}

func New(repo *repositories.Repo, baseURL string) *Encoder {
	return &Encoder{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (e *Encoder) encode(r *frontoffice.Record) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString(startTag)

	var t string
	switch r.Type {
	case "book", "bookEditor":
		t = "book"
	case "bookChapter":
		t = "bookPart"
	case "conference":
		t = "conferenceObject"
	case "dissertation":
		t = "doctoralThesis"
	case "journalArticle":
		t = "article"
	default:
		t = "other"
	}
	switch r.MiscType {
	case "newsArticle", "newspaperPiece", "magazinePiece":
		t = "contributionToPeriodical"
	case "bookReview":
		t = "review"
	case "report":
		t = "report"
	}

	writeField(b, "type", t)
	writeField(b, "type", "info:eu-repo/semantics/"+t)

	writeField(b, "identifier", "oai:archive.ugent.be:"+r.ID)
	writeField(b, "identifier", r.Handle)
	writeField(b, "identifier", e.baseURL+"/publication/"+r.ID)
	for _, val := range r.DOI {
		writeField(b, "identifier", identifiers.DOI.Resolve(val))
	}

	switch r.PublicationStatus {
	case "unsubmitted":
		writeField(b, "type", "info:eu-repo/semantics/draft")
	case "inpress":
		writeField(b, "type", "info:eu-repo/semantics/acceptedVersion")
	default:
		writeField(b, "type", "info:eu-repo/semantics/publishedVersion")
	}

	for _, val := range r.Author {
		writeField(b, "creator", val.Name)
	}
	for _, val := range r.Promoter {
		writeField(b, "contributor", val.Name)
	}
	if r.Type == "bookEditor" || r.Type == "issueEditor" {
		for _, val := range r.Editor {
			writeField(b, "creator", val.Name)
		}
	} else {
		for _, val := range r.Editor {
			writeField(b, "contributor", val.Name)
		}
	}

	for _, val := range r.Keyword {
		writeField(b, "subject", val)
	}
	for _, val := range r.Subject {
		writeField(b, "subject", val)
	}

	writeField(b, "title", r.Title)

	if r.Publisher != nil {
		writeField(b, "publisher", r.Publisher.Name)
	}

	writeField(b, "date", r.Year)

	for _, val := range r.Abstract {
		writeField(b, "description", val)
	}

	for _, val := range r.Language {
		writeField(b, "language", val)
	}

	writeField(b, "rights", r.CopyrightStatement)

	if r.Parent != nil {
		writeField(b, "source", r.Parent.Title)
	}

	for _, val := range r.ISSN {
		writeField(b, "source", "ISSN: "+val)
	}

	if r.Parent != nil {
		for _, val := range r.ISBN {
			writeField(b, "source", "ISBN: "+val)
		}
	} else {
		for _, val := range r.ISBN {
			writeField(b, "identifier", "urn:isbn:"+val)
		}
	}

	if len(r.File) > 0 {
		f := r.File[0]
		writeField(b, "identifier", e.baseURL+"/publication/"+r.ID+"/file/"+f.ID)
		writeField(b, "format", f.ContentType)
		if f.Change != nil && f.Change.To == "open" {
			writeField(b, "rights", "info:eu-repo/semantics/embargoedAccess")
			if f.Change.On != "" {
				writeField(b, "date", "info:eu-repo/date/embargoEnd/"+f.Change.On[0:10])
			}
		} else if f.Access == "open" {
			writeField(b, "rights", "info:eu-repo/semantics/openAccess")
		} else if f.Access == "restricted" {
			writeField(b, "rights", "info:eu-repo/semantics/restrictedAccess")
		} else if f.Access == "private" {
			writeField(b, "rights", "info:eu-repo/semantics/closedAccess")
		}
	}

	for _, val := range r.Project {
		if val.EUFrameworkProgramme != "" && val.EUID != "" {
			writeField(b, "relation", fmt.Sprintf("info:eu-repo/grantAgreement/EC/%s/%s", val.EUFrameworkProgramme, val.EUID))
		}
	}

	b.WriteString(endTag)

	return b.Bytes(), nil
}

func (e *Encoder) EncodePublication(p *models.Publication) ([]byte, error) {
	return e.encode(frontoffice.MapPublication(p, e.repo))
}

func (e *Encoder) EncodeDataset(d *models.Dataset) ([]byte, error) {
	return e.encode(frontoffice.MapDataset(d, e.repo))
}
