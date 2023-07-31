package ris

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

var (
	reTag        = regexp.MustCompile(`^([A-Z0-9]{2})\s+-?\s*(.*)$`)
	reSplit      = regexp.MustCompile(`\s*[,;]\s*`)
	wosLanguages = map[string]string{
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
		"Rumanian":       "rum",
		"Russian":        "rus",
		"Serbian":        "srp",
		// "Serbo-Croatian": "", // no current equivalent in iso 639-2
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

		switch line {
		case "ER", "ER  -":
			mapRecord(rec, p)
			return nil
		case "EF":
			return io.EOF
		}

		if match := reTag.FindStringSubmatch(line); match != nil {
			tag = match[1]
			rec[tag] = append(rec[tag], match[2])
		} else {
			rec[tag] = append(rec[tag], strings.TrimPrefix(line, "  "))
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
			// give preference to AF over AU
			if k == "AU" && p.Author != nil {
				continue
			}
			if k == "AF" {
				p.Author = nil
			}
			for _, val := range v {
				nameParts := reSplit.Split(val, -1)
				lastName := nameParts[0]
				firstName := "[missing]" // TODO
				if len(nameParts) > 1 {
					firstName = nameParts[1]
				}
				c := models.ContributorFromFirstLastName(firstName, lastName)
				p.Author = append(p.Author, c)
			}
		case "TI", "T1":
			p.Title = strings.Join(v, "")
		case "AB", "N2":
			p.AddAbstract(&models.Text{Text: strings.Join(v, "\n\n"), Lang: "eng"})
		case "KW", "DW", "ID", "DE":
			p.Keyword = append(p.Keyword, splitMultilineVals(v)...)
		case "DI":
			p.DOI = v[0]
		case "JF", "JO", "T2":
			p.Publication = v[0]
		case "SO":
			p.Publication = strings.Join(v, "")
		case "JA", "JI":
			p.PublicationAbbreviation = v[0]
		case "SN":
			p.ISSN = append(p.ISSN, splitVals(v)...)
		case "EI":
			p.EISSN = append(p.EISSN, splitVals(v)...)
		case "BN":
			p.EISBN = append(p.EISBN, splitVals(v)...)
		case "UT":
			p.WOSID = strings.TrimPrefix(v[0], "WOS:")
		case "AN":
			if strings.HasPrefix(v[0], "WOS:") {
				p.WOSID = strings.TrimPrefix(v[0], "WOS:")
			}
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
		case "BP", "SP":
			p.PageFirst = v[0]
		case "EP":
			p.PageLast = v[0]
		case "PG":
			p.PageCount = v[0]
		case "AP", "AR":
			p.ArticleNumber = v[0]
		case "PB", "PU":
			p.Publisher = v[0]
		case "PI":
			p.PlaceOfPublication = v[0]
		case "CT":
			p.ConferenceName = strings.Join(v, "")
		case "CL":
			p.ConferenceLocation = v[0]
		case "CY":
			date := parseConferenceDate(v[0])
			p.ConferenceStartDate = date[0]
			p.ConferenceEndDate = date[1]
		case "LA":
			for _, val := range v {
				if code, ok := wosLanguages[val]; ok {
					p.Language = append(p.Language, code)
				}
			}
		}
	}
}

func splitVals(vals []string) (newVals []string) {
	for _, v := range vals {
		for _, str := range reSplit.Split(v, -1) {
			if str != "" {
				newVals = append(newVals, str)
			}
		}
	}
	return
}

func splitMultilineVals(vals []string) (newVals []string) {
	for _, v := range reSplit.Split(strings.Join(vals, ""), -1) {
		if v != "" {
			newVals = append(newVals, v)
		}
	}
	return
}

func parseConferenceDate(date string) [2]string {
	var parts []string
	var result [2]string
	var r *regexp.Regexp

	result[0] = ""
	result[1] = ""

	// Match for a YYYY
	if ok, _ := regexp.MatchString("^[0-9]{4}$", date); ok {
		result[0] = date
		result[1] = date
	}

	// Match for a MMM, YYYY
	if ok, _ := regexp.MatchString("^[A-Z]{3}, [0-9]{4}$", date); ok {
		result[0] = date
		result[1] = date
	}

	// Match for a MMM DD, YYYY
	if ok, _ := regexp.MatchString("^[A-Z]{3}[\\s]*[0-9]{2}, [0-9]{4}$", date); ok {
		result[0] = date
		result[1] = date
	}

	// Match for MMM DD-DD, YYYY
	if ok, _ := regexp.MatchString("^[A-Z]{3} [0-9]{2}-[0-9]{2}, [0-9]{4}$", date); ok {
		r, _ = regexp.Compile("[0-9]{4}$")
		year := r.FindString(date)

		r, _ = regexp.Compile("^[A-Z]{3} [0-9]{2}-[0-9]{2}")
		daysMonth := r.FindString(date)

		parts = strings.Split(daysMonth, "-")
		r, _ = regexp.Compile("^[A-Z]{3}")
		month := r.FindString(parts[0])

		result[0] = fmt.Sprintf("%s, %s", parts[0], year)
		result[1] = fmt.Sprintf("%s %s, %s", month, parts[1], year)
	}

	// Match for a MMM DD-MMM DD, YYYY
	if ok, _ := regexp.MatchString("^[A-Z]{3} [0-9]{2}-[A-Z]{3} [0-9]{2}, [0-9]{4}$", date); ok {
		parts = strings.Split(date, "-")
		r, _ = regexp.Compile("[0-9]{4}$")
		year := r.FindString(parts[1])

		result[0] = fmt.Sprintf("%s, %s", parts[0], year)
		result[1] = parts[1]
	}

	return result
}
