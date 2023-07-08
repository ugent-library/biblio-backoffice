package oaidc

import (
	"bytes"

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
		b.WriteString(val)
		b.WriteString("</dc:")
		b.WriteString(tag)
		b.WriteString(">")
	}
}

func EncodePublication(p *models.Publication) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString(startTag)

	writeField(b, "identifier", p.Handle)
	writeField(b, "title", p.Title)
	writeField(b, "date", p.Year)
	writeField(b, "publisher", p.Publisher)
	if p.DOI != "" {
		writeField(b, "identifier", identifiers.DOI.Resolve(p.DOI))
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
		writeField(b, "creator", val.Person.FullName)
	}
	for _, val := range p.Supervisor {
		writeField(b, "contributor", val.Person.FullName)
	}
	if p.Type == "book_editor" || p.Type == "issue_editor" {
		for _, val := range p.Editor {
			writeField(b, "creator", val.Person.FullName)
		}
	} else {
		for _, val := range p.Editor {
			writeField(b, "contributor", val.Person.FullName)
		}
	}

	b.WriteString(endTag)

	return b.Bytes(), nil
}

// my $TYPES = {
//     book             => 'book',
//     bookChapter      => 'bookPart',
//     bookEditor       => 'book',
//     conference       => 'conferenceObject',
//     dissertation     => 'doctoralThesis',
//     journalArticle   => 'article',
//     licentiateThesis => 'masterThesis',
// };

// my $VERSIONS = {
//     unsubmitted => 'draft',
//     inpress     => 'acceptedVersion',
//     accepted    => 'acceptedVersion',
//     published   => 'publishedVersion',
// };

// sub fix {
//     state $uri_base = Catmandu->config->{uri_base} . '/publication';

//     my ($self, $pub) = @_;

//     my $type = $TYPES->{$pub->{type}} || 'other';

//     if ($pub->{misc_type}) {
//         if ($pub->{misc_type} eq 'newsArticle' || $pub->{misc_type} eq 'newspaperPiece' || $pub->{misc_type} eq 'magazinePiece') {
//             $type = 'contributionToPeriodical';
//         } elsif ($pub->{misc_type} eq 'bookReview') {
//             $type = 'review';
//         } elsif ($pub->{misc_type} eq 'report') {
//             $type = 'report';
//         }
//     }

//     my $dc = {
//         identifier => [ "$uri_base/$pub->{_id}" ],
//         type => [ $pub->{type}, "info:eu-repo/semantics/$type" ],
//     };

//     if ($pub->{publication_status}) {
//         if (my $version = $VERSIONS->{$pub->{publication_status}}) {
//             push @{$dc->{type}}, "info:eu-repo/semantics/$version";
//         }
//     }
//     $dc->{rights}      = [ $pub->{copyright_statement} ] if $pub->{copyright_statement};
//     $dc->{source}      = [ $pub->{parent}{title} ]       if $pub->{parent} && $pub->{parent}{title};

//     if ($pub->{file}) {
//         if (my $file = $pub->{file}->[0]) {
//             push @{$dc->{identifier} ||= []}, "$uri_base/$pub->{_id}/file/$file->{_id}";
//             $dc->{format} = [ $file->{content_type} ];
//             if ($file->{change} && $file->{change}{to} eq 'open') {
//                 push @{$dc->{rights} ||= []}, "info:eu-repo/semantics/embargoedAccess";
//                 if ($file->{change}{on}) {
//                     push @{$dc->{date} ||= []}, "info:eu-repo/date/embargoEnd/" . substr($file->{change}{on}, 0, 10);
//                 }
//             } elsif ($file->{access}) {
//                 if ($file->{access} eq 'open')       { push @{$dc->{rights} ||= []}, "info:eu-repo/semantics/openAccess" }
//                 if ($file->{access} eq 'restricted') { push @{$dc->{rights} ||= []}, "info:eu-repo/semantics/restrictedAccess" }
//                 if ($file->{access} eq 'private')    { push @{$dc->{rights} ||= []}, "info:eu-repo/semantics/closedAccess" }
//             }
//         }
//     }

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
