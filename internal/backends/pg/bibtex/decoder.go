package bibtex

import (
	"io"
	"regexp"
	"strings"

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
	return &Decoder{r: r}
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

	if f, ok := e.Fields["title"]; ok {
		p.Title = f.String()
	}
	if f, ok := e.Fields["year"]; ok {
		p.Year = f.String()
	}
	if f, ok := e.Fields["pages"]; ok {
		pageParts := reSplitPages.Split(f.String(), -1)
		p.PageFirst = pageParts[0]
		if len(pageParts) > 1 {
			p.PageLast = pageParts[1]
		}
	}
	if f, ok := e.Fields["keywords"]; ok {
		p.Keyword = reSplit.Split(f.String(), -1)
	}
	if f, ok := e.Fields["abstract"]; ok {
		p.Abstract = []models.Text{{Text: f.String(), Lang: "und"}}
	}
	if f, ok := e.Fields["volume"]; ok {
		p.Volume = f.String()
	}
	if f, ok := e.Fields["number"]; ok {
		p.Issue = f.String()
	}
	if f, ok := e.Fields["address"]; ok {
		p.PlaceOfPublication = f.String()
	}
	if f, ok := e.Fields["doi"]; ok {
		p.DOI = f.String()
	}
	if f, ok := e.Fields["issn"]; ok {
		p.ISSN = []string{f.String()}
	}
	if f, ok := e.Fields["isbn"]; ok {
		p.ISBN = []string{f.String()}
	}
	if f, ok := e.Fields["series"]; ok {
		p.SeriesTitle = f.String()
	}
	if f, ok := e.Fields["journal"]; ok {
		p.Publication = f.String()
	}
	if f, ok := e.Fields["booktitle"]; ok {
		p.Publication = f.String()
	}
	if f, ok := e.Fields["school"]; ok {
		p.Publisher = f.String()
	}
	if f, ok := e.Fields["publisher"]; ok {
		p.Publisher = f.String()
	}
	if f, ok := e.Fields["author"]; ok {
		for _, v := range strings.Split(f.String(), " and ") {
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
	if f, ok := e.Fields["editor"]; ok {
		for _, v := range strings.Split(f.String(), " and ") {
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
	if f, ok := e.Fields["journal-iso"]; ok {
		p.PublicationAbbreviation = f.String()
	}
	if f, ok := e.Fields["keywords-plus"]; ok {
		p.Keyword = append(p.Keyword, reSplit.Split(f.String(), -1)...)
	}
	if f, ok := e.Fields["unique-id"]; ok {
		if strings.HasPrefix(f.String(), "ISI:") {
			p.WOSID = strings.TrimPrefix(f.String(), "ISI:")
		}
	}
}
