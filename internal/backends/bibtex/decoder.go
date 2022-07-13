package bibtex

import (
	"io"
	"regexp"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/nickng/bibtex"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
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
	// bibtex fails on utf8 bom
	return &Decoder{r: utfbom.SkipOnly(r)}
}

func (d *Decoder) Decode(p *models.Publication) error {
	if d.bibtex == nil {
		b, err := bibtex.Parse(d.r)
		if err != nil {
			return err
		}
		d.bibtex = b
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
		p.Type = "disertation"
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
		p.Abstract = []models.Text{{Text: f, Lang: "und"}}
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
			c := &models.Contributor{FullName: v, LastName: nameParts[0]}
			if len(nameParts) > 1 {
				c.FirstName = nameParts[1]
			} else {
				c.FirstName = "[missing]" // TODO
			}
			p.Author = append(p.Author, c)
		}
	}
	if f, ok := entries["editor"]; ok {
		for _, v := range strings.Split(f, " and ") {
			nameParts := reSplit.Split(v, -1)
			c := &models.Contributor{FullName: v, LastName: nameParts[0]}
			if len(nameParts) > 1 {
				c.FirstName = nameParts[1]
			} else {
				c.FirstName = "[missing]" // TODO
			}
			p.Editor = append(p.Editor, c)
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
