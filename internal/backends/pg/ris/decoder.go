package ris

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

var (
	reTag   = regexp.MustCompile(`^([A-Z0-9]{2})\s+-?\s*(.*)$`)
	reSplit = regexp.MustCompile(`\s*[,;]\s*`)
)

type Record map[string][]string

type Decoder struct {
	scanner *bufio.Scanner
}

func NewDecoder(r io.Reader) backends.PublicationDecoder {
	return &Decoder{scanner: bufio.NewScanner(r)}
}

func (d *Decoder) Decode(p *models.Publication) error {
	rec := make(Record)
	tag := ""

	for d.scanner.Scan() {
		line := d.scanner.Text()
		match := reTag.FindStringSubmatch(line)

		if match == nil {
			rec[tag] = append(rec[tag], strings.TrimPrefix(line, "  "))
		} else if match[1] == "ER" {
			mapRecord(rec, p)
			return nil
		} else {
			tag = match[1]
			rec[tag] = append(rec[tag], match[2])
		}
	}

	if d.scanner.Err() != nil {
		return d.scanner.Err()
	}

	return io.EOF
}

// TODO map_wos_conference_dates(), clean_wos_type() fixes
func mapRecord(r Record, p *models.Publication) {
	p.Type = "journal_article"
	p.PublicationStatus = "published"

	for k, v := range r {
		switch k {
		case "TY", "DT", "PT":
			p.WOSType = v[0]
			switch v[0] {
			case "Art Exhibit Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "exhibitionReview"
			case "Book Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "bookReview"
			case "Dance Performance Review", "Theater Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "theatreReview"
			case "Database Review", "Hardware Review", "Software Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "productReview"
			case "Editorial Material":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "editorialMaterial"
			case "Fiction, Creative Prose", "Poetry", "Script":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "artisticWork"
			case "Film Review", "TV Review, Radio Review", "TV Review, Radio Review, Video Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "filmReview"
			case "Music Score Review", "Music Performance Review", "Record Review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "musicReview"
			case "Music Score":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "musicEdition"
			case "News Item":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "newsArticle"
			case "Data Paper":
				p.Type = "miscellaneous"
			case "BOOK", "Book, Whole":
				p.Type = "book"
			case "Book chapter", "CHAP":
				p.Type = "book_chapter"
			case "Meeting Abstract":
				p.Type = "conference"
				p.ConferenceType = "abstract"
			case "C", "CONF", "Conference proceeding", "Proceedings Paper", "S":
				p.Type = "conference"
			}
		case "AF", "AU":
			for _, val := range v {
				nameParts := reSplit.Split(val, -1)
				c := &models.Contributor{FullName: val, LastName: nameParts[0]}
				if len(nameParts) > 1 {
					c.FirstName = nameParts[1]
				}
				p.Author = append(p.Author, c)
			}
		case "TI", "T1":
			p.Title = v[0]
		case "AB", "N2":
			p.Abstract = []models.Text{{Text: v[0], Lang: "und"}}
		case "KW", "DW", "ID", "DE":
			for _, val := range v {
				p.Keyword = append(p.Keyword, reSplit.Split(val, -1)...)
			}
		case "DI":
			p.DOI = v[0]
		case "JF", "JO", "T2", "SO":
			p.Publication = v[0]
		case "JA", "JI":
			p.PublicationAbbreviation = v[0]
		case "SN":
			for _, val := range v {
				p.ISSN = append(p.ISSN, reSplit.Split(val, -1)...)
			}
		case "EI":
			for _, val := range v {
				p.EISSN = append(p.EISSN, reSplit.Split(val, -1)...)
			}
		case "BN":
			for _, val := range v {
				p.ISBN = append(p.ISBN, reSplit.Split(val, -1)...)
			}
		case "UT":
			p.WOSID = strings.TrimPrefix(v[0], "WOS:")
		case "PM":
			p.PubMedID = v[0]
		case "VL":
			p.Volume = v[0]
		case "CP", "IS":
			p.Issue = v[0]
		case "T3", "SE":
			p.SeriesTitle = v[0]
		case "Y1", "PY":
			p.Year = v[0]
		case "BP":
			p.PageFirst = v[0]
		case "EP":
			p.PageLast = v[0]
		case "PG":
			p.PageCount = v[0]
		case "AP", "AR":
			p.ArticleNumber = v[0]
		case "PB", "PU":
			p.Publisher = v[0]
		case "CY", "PI":
			p.PlaceOfPublication = v[0]
		case "CT":
			p.Conference.Name = v[0]
		case "CL":
			p.Conference.Location = v[0]
		}
	}
}
