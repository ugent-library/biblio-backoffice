package oaidc

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/internal/models"
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

// TODO copied from frontoffice handler, DRY this
var licenses = map[string]string{
	"CC0-1.0":          "Creative Commons Public Domain Dedication (CC0 1.0)",
	"CC-BY-4.0":        "Creative Commons Attribution 4.0 International Public License (CC-BY 4.0)",
	"CC-BY-SA-4.0":     "Creative Commons Attribution-ShareAlike 4.0 International Public License (CC BY-SA 4.0)",
	"CC-BY-NC-4.0":     "Creative Commons Attribution-NonCommercial 4.0 International Public License (CC BY-NC 4.0)",
	"CC-BY-ND-4.0":     "Creative Commons Attribution-NoDerivatives 4.0 International Public License (CC BY-ND 4.0)",
	"CC-BY-NC-SA-4.0":  "Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International Public License (CC BY-NC-SA 4.0)",
	"CC-BY-NC-ND-4.0":  "Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International Public License (CC BY-NC-ND 4.0)",
	"InCopyright":      "No license (in copyright)",
	"LicenseNotListed": "A specific license has been chosen by the rights holder. Get in touch with the rights holder for reuse rights.",
	"CopyrightUnknown": "Information pending",
	"":                 "No license (in copyright)",
}

var openLicenses = map[string]struct{}{
	"CC0-1.0":         {},
	"CC-BY-4.0":       {},
	"CC-BY-SA-4.0":    {},
	"CC-BY-NC-4.0":    {},
	"CC-BY-ND-4.0":    {},
	"CC-BY-NC-SA-4.0": {},
	"CC-BY-NC-ND-4.0": {},
}

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
	baseURL string
}

func New(baseURL string) *Encoder {
	return &Encoder{
		baseURL: baseURL,
	}
}

func (e *Encoder) EncodePublication(p *models.Publication) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString(startTag)

	var t string
	switch p.Type {
	case "book", "book_editor":
		t = "book"
	case "book_chapter":
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
	switch p.MiscellaneousType {
	case "newsArticle", "newspaperPiece":
		t = "contributionToPeriodical"
	case "bookReview":
		t = "review"
	case "report":
		t = "report"
	}
	writeField(b, "type", t)
	writeField(b, "type", "info:eu-repo/semantics/"+t)

	writeField(b, "identifier", "oai:archive.ugent.be:"+p.ID)
	writeField(b, "identifier", p.Handle)
	writeField(b, "identifier", e.baseURL+"/publication/"+p.ID)
	if p.DOI != "" {
		writeField(b, "identifier", identifiers.DOI.Resolve(p.DOI))
	}

	switch p.PublicationStatus {
	case "unpublished":
		writeField(b, "type", "info:eu-repo/semantics/draft")
	case "accepted":
		writeField(b, "type", "info:eu-repo/semantics/acceptedVersion")
	default:
		writeField(b, "type", "info:eu-repo/semantics/publishedVersion")
	}

	writeField(b, "title", p.Title)

	writeField(b, "publisher", p.Publisher)

	writeField(b, "date", p.Year)

	if p.Publication != "" {
		writeField(b, "source", p.Publication)
	}

	if p.Publication != "" || p.PublicationAbbreviation != "" {
		for _, val := range p.ISBN {
			writeField(b, "source", "ISBN: "+val)
		}
		for _, val := range p.EISBN {
			writeField(b, "source", "ISBN: "+val)
		}
	} else {
		for _, val := range p.ISBN {
			writeField(b, "identifier", "urn:isbn:"+val)
		}
		for _, val := range p.EISBN {
			writeField(b, "identifier", "urn:isbn:"+val)
		}
	}

	for _, val := range p.ISSN {
		writeField(b, "source", "ISSN: "+val)
	}
	for _, val := range p.EISSN {
		writeField(b, "source", "ISSN: "+val)
	}

	for _, val := range p.Language {
		writeField(b, "language", val)
	}

	for _, val := range p.Abstract {
		writeField(b, "description", val.Text)
	}

	for _, val := range p.Keyword {
		writeField(b, "subject", val)
	}
	for _, val := range p.ResearchField {
		writeField(b, "subject", val)
	}

	for _, val := range p.Author {
		writeField(b, "creator", val.Name())
	}
	for _, val := range p.Supervisor {
		writeField(b, "contributor", val.Name())
	}
	if p.Type == "book_editor" || p.Type == "issue_editor" {
		for _, val := range p.Editor {
			writeField(b, "creator", val.Name())
		}
	} else {
		for _, val := range p.Editor {
			writeField(b, "contributor", val.Name())
		}
	}

	if len(p.File) > 0 {
		f := p.File[0]
		writeField(b, "identifier", e.baseURL+"/publication/"+p.ID+"/file/"+f.ID)
		writeField(b, "format", f.ContentType)
		writeField(b, "rights", f.AccessLevel)
		if f.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
			writeField(b, "date", "info:eu-repo/date/embargoEnd/"+f.EmbargoDate)
		}
	}

	if len(p.File) > 0 {
		bestLicense := ""
		for _, f := range p.File {
			if bestLicense == "" {
				if _, isLicense := licenses[f.License]; isLicense {
					bestLicense = f.License
				}
			}
			if _, isOpenLicense := openLicenses[f.License]; isOpenLicense {
				bestLicense = f.License
				break
			}
		}

		writeField(b, "rights", licenses[bestLicense])
	}

	for _, val := range p.RelatedProjects {
		eu := val.Project.EUProject
		if eu != nil && eu.ID != "" && eu.FrameworkProgramme != "" {
			writeField(b, "relation", fmt.Sprintf("info:eu-repo/grantAgreement/EC/%s/%s", eu.FrameworkProgramme, eu.ID))
		}
	}

	b.WriteString(endTag)

	return b.Bytes(), nil
}
