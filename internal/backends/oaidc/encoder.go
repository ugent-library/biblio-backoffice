package oaidc

import (
	"bytes"
	"encoding/xml"

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

func EncodePublication(p *models.Publication) ([]byte, error) {
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

	writeField(b, "date", p.Year)

	writeField(b, "publisher", p.Publisher)

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

	for _, val := range p.File {
		writeField(b, "format", val.ContentType)
		writeField(b, "rights", val.AccessLevel)
		if val.AccessLevel == "info:eu-repo/semantics/embargoedAccess" {
			writeField(b, "date", "info:eu-repo/date/embargoEnd/"+val.EmbargoDate)
		}

		break
	}

	b.WriteString(endTag)

	return b.Bytes(), nil
}

//     $dc->{rights}      = [ $pub->{copyright_statement} ] if $pub->{copyright_statement};

//     if (my $projects = $pub->{project}) {
//         for my $project (@$projects) {
//             if ($project->{eu_id} && $project->{eu_framework_programme}) {
//                 push @{$dc->{relation} ||= []}, "info:eu-repo/grantAgreement/EC/$project->{eu_framework_programme}/$project->{eu_id}";
//             }
//         }
//     }

//     $dc;
// }

// 1;
