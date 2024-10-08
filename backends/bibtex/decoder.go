package bibtex

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/bibtex"
)

var (
	reSplit      = regexp.MustCompile(`\s*[,;]\s*`)
	reSplitPages = regexp.MustCompile(`\s*[-\x{2013}\x{2014}]+\s*`)
)

type Decoder struct {
	parser *bibtex.Parser
}

func NewDecoder(r io.Reader) backends.PublicationDecoder {
	return &Decoder{parser: bibtex.NewParser(r)}
}

func (d *Decoder) Decode(p *models.Publication) error {
	entry, err := d.parser.Next()
	if err != nil {
		return fmt.Errorf("bibtex: parse: %w", err)
	}
	if entry == nil {
		return io.EOF
	}

	mapEntry(entry, p)

	return nil
}

func mapEntry(e *bibtex.Entry, p *models.Publication) {
	p.Type = "journal_article"

	fields := make(map[string]string, len(e.Fields))
	for _, f := range e.Fields {
		fields[f.Name] = f.Value
	}

	switch e.Type {
	case "article":
		p.Type = "journal_article"
	case "book", "booklet":
		p.Type = "book"
	case "inbook", "incollection":
		p.Type = "book_chapter"
	case "conference", "proceedings", "inproceedings":
		p.Type = "conference"
	case "phdthesis":
		p.Type = "dissertation"
	case "unpublished":
		p.Type = "miscellaneous"
		p.MiscellaneousType = "preprint"
	case "techreport":
		p.Type = "miscellaneous"
		p.MiscellaneousType = "report"
	}

	for _, name := range e.Authors {
		nameParts := reSplit.Split(name, -1)
		lastName := nameParts[0]
		firstName := "[missing]" // TODO
		if len(nameParts) > 1 {
			firstName = nameParts[1]
		}
		p.Author = append(p.Author, models.ContributorFromFirstLastName(firstName, lastName))
	}

	for _, name := range e.Editors {
		nameParts := reSplit.Split(name, -1)
		lastName := nameParts[0]
		firstName := "[missing]" // TODO
		if len(nameParts) > 1 {
			firstName = nameParts[1]
		}
		p.Editor = append(p.Editor, models.ContributorFromFirstLastName(firstName, lastName))
	}

	if f, ok := fields["title"]; ok {
		p.Title = f
	}
	if f, ok := fields["year"]; ok {
		p.Year = f
	}
	if f, ok := fields["pages"]; ok {
		pageParts := reSplitPages.Split(f, -1)
		p.PageFirst = pageParts[0]
		if len(pageParts) > 1 {
			p.PageLast = pageParts[1]
		}
	}
	if f, ok := fields["keywords"]; ok {
		p.Keyword = reSplit.Split(f, -1)
	}
	if f, ok := fields["abstract"]; ok {
		p.AddAbstract(&models.Text{Text: f, Lang: "und"})
	}
	if f, ok := fields["volume"]; ok {
		p.Volume = f
	}
	if f, ok := fields["number"]; ok {
		p.Issue = f
	}
	if f, ok := fields["address"]; ok {
		p.PlaceOfPublication = f
	}
	if f, ok := fields["doi"]; ok {
		p.DOI = f
	}
	if f, ok := fields["issn"]; ok {
		p.ISSN = []string{f}
	}
	if f, ok := fields["isbn"]; ok {
		p.ISBN = []string{f}
	}
	if f, ok := fields["series"]; ok {
		p.SeriesTitle = f
	}
	if f, ok := fields["journal"]; ok {
		p.Publication = f
	}
	if f, ok := fields["booktitle"]; ok {
		p.Publication = f
	}
	if f, ok := fields["school"]; ok {
		p.Publisher = f
	}
	if f, ok := fields["publisher"]; ok {
		p.Publisher = f
	}

	// WoS bibtex records
	if f, ok := fields["journal-iso"]; ok {
		p.PublicationAbbreviation = f
	}
	if f, ok := fields["keywords-plus"]; ok {
		p.Keyword = append(p.Keyword, reSplit.Split(f, -1)...)
	}
	if f, ok := fields["unique-id"]; ok {
		if strings.HasPrefix(f, "ISI:") {
			p.WOSID = strings.TrimPrefix(f, "ISI:")
		}
	}

	// WOS language
	if f, ok := fields["language"]; ok {
		if lang, ok := wosLangs[f]; ok {
			p.Language = []string{lang}
		}
	}
}

// Mapping between value in WOS Bibtex Language field and ISO-396-2
// code is made based on https://www.loc.gov/standards/iso639-2/php/code_list.php
// Languages which differentiate between a bibliographic and terminology code,
// preference was given to the B code.
// WOS Language list: https://images.webofknowledge.com/images/help/WOS/hs_language.html
var wosLangs = map[string]string{
	"English":        "eng",
	"Afrikaans":      "afr",
	"Arabic":         "ara",
	"Basque":         "baq",
	"Bengali":        "ben",
	"Bulgarian":      "bul",
	"Byelorussian":   "bel",
	"Catalan":        "cat",
	"Chinese":        "chi",
	"Croatian":       "hrv",
	"Czech":          "cze",
	"Danish":         "dan",
	"Dutch":          "dut",
	"Estonian":       "est",
	"Finnish":        "fin",
	"Flemish":        "dut",
	"French":         "fre",
	"Gaelic":         "gla",
	"Galician":       "glg",
	"Georgian":       "geo",
	"German":         "ger",
	"Greek":          "gre",
	"Hebrew":         "heb",
	"Hungarian":      "hun",
	"Icelandic":      "ice",
	"Italian":        "ita",
	"Japanese":       "jpn",
	"Korean":         "kor",
	"Latin":          "lat",
	"Latvian":        "lav",
	"Lithuanian":     "lit",
	"Macedonian":     "mac",
	"Malay":          "may",
	"Multi-Language": "mul",
	"Norwegian":      "nor",
	"Persian":        "per",
	"Polish":         "pol",
	"Portuguese":     "por",
	"Provencal":      "pro",
	"Romansch":       "roh",
	"Rumanian":       "rup",
	"Russian":        "rus",
	"Serbian":        "srp",
	// "Serbo-Croatian": "" not supported in ISO 369-2
	"Slovak":      "slo",
	"Slovenian":   "slv",
	"Spanish":     "spa",
	"Swedish":     "swe",
	"Thai":        "tha",
	"Turkish":     "tur",
	"Ukrainian":   "ukr",
	"Unspecified": "und",
	"Welsh":       "wel",
}
