package bibtex

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/dimchansky/utfbom"
	"github.com/nickng/bibtex"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	reSplit      = regexp.MustCompile(`\s*[,;]\s*`)
	reSplitPages = regexp.MustCompile(`\s*[-\x{2013}\x{2014}]+\s*`)
)

type Decoder struct {
	r      io.Reader
	bibtex *bibtex.BibTex
	i      int
}

func NewDecoder(r io.Reader) backends.PublicationDecoder {
	return &Decoder{r: r}
}

func (d *Decoder) parse() error {
	// cleanup
	var r io.Reader
	// remove utf8 bom
	r = utfbom.SkipOnly(d.r)
	// remove unicode non spacing marks
	// note that the parser doens't actually fail on combined grave, acute, circumflex, umlaut accents in field values
	r = transform.NewReader(r, transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC))
	// skip file preambles, comments, etc until we encounter the first entry
	b := bufio.NewReader(r)
	for {
		c, _, err := b.ReadRune()
		if err != nil {
			return err
		}
		if c == '@' {
			b.UnreadRune()
			break
		}
	}

	bib, err := bibtex.Parse(b)
	if err != nil {
		return err
	}
	d.bibtex = bib

	return nil
}

func (d *Decoder) Decode(p *models.Publication) error {
	if d.bibtex == nil {
		if err := d.parse(); err != nil {
			return err
		}
	}

	if len(d.bibtex.Entries) == 0 || d.i >= len(d.bibtex.Entries) {
		return io.EOF
	}

	entry := d.bibtex.Entries[d.i]
	d.i++

	mapEntry(entry, p)

	return nil
}

func mapEntry(e *bibtex.BibEntry, p *models.Publication) {
	p.Type = "journal_article"

	// field names may have capitals
	entries := map[string]string{}
	for key, val := range e.Fields {
		entries[strings.ToLower(key)] = val.String()
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

	if f, ok := entries["title"]; ok {
		p.Title = f
	}
	if f, ok := entries["year"]; ok {
		p.Year = f
	}
	if f, ok := entries["pages"]; ok {
		pageParts := reSplitPages.Split(f, -1)
		p.PageFirst = pageParts[0]
		if len(pageParts) > 1 {
			p.PageLast = pageParts[1]
		}
	}
	if f, ok := entries["keywords"]; ok {
		p.Keyword = reSplit.Split(f, -1)
	}
	if f, ok := entries["abstract"]; ok {
		p.AddAbstract(&models.Text{Text: f, Lang: "und"})
	}
	if f, ok := entries["volume"]; ok {
		p.Volume = f
	}
	if f, ok := entries["number"]; ok {
		p.Issue = f
	}
	if f, ok := entries["address"]; ok {
		p.PlaceOfPublication = f
	}
	if f, ok := entries["doi"]; ok {
		p.DOI = f
	}
	if f, ok := entries["issn"]; ok {
		p.ISSN = []string{f}
	}
	if f, ok := entries["isbn"]; ok {
		p.ISBN = []string{f}
	}
	if f, ok := entries["series"]; ok {
		p.SeriesTitle = f
	}
	if f, ok := entries["journal"]; ok {
		p.Publication = f
	}
	if f, ok := entries["booktitle"]; ok {
		p.Publication = f
	}
	if f, ok := entries["school"]; ok {
		p.Publisher = f
	}
	if f, ok := entries["publisher"]; ok {
		p.Publisher = f
	}
	if f, ok := entries["author"]; ok {
		for _, v := range strings.Split(f, " and ") {
			nameParts := reSplit.Split(v, -1)
			lastName := nameParts[0]
			firstName := "[missing]" // TODO
			if len(nameParts) > 1 {
				firstName = nameParts[1]
			}
			p.Author = append(p.Author, models.ContributorFromFirstLastName(firstName, lastName))
		}
	}
	if f, ok := entries["editor"]; ok {
		for _, v := range strings.Split(f, " and ") {
			nameParts := reSplit.Split(v, -1)
			lastName := nameParts[0]
			firstName := "[missing]" // TODO
			if len(nameParts) > 1 {
				firstName = nameParts[1]
			}
			p.Editor = append(p.Editor, models.ContributorFromFirstLastName(firstName, lastName))
		}
	}

	// WoS bibtex records
	if f, ok := entries["journal-iso"]; ok {
		p.PublicationAbbreviation = f
	}
	if f, ok := entries["keywords-plus"]; ok {
		p.Keyword = append(p.Keyword, reSplit.Split(f, -1)...)
	}
	if f, ok := entries["unique-id"]; ok {
		if strings.HasPrefix(f, "ISI:") {
			p.WOSID = strings.TrimPrefix(f, "ISI:")
		}
	}
}