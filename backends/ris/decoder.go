package ris

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"slices"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
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

		switch strings.TrimSpace(line) {
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

	if err := d.scanner.Err(); err != nil {
		return fmt.Errorf("ris: line scanner: %w", err)
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
			types := reSplit.Split(p.WOSType, -1)
			for i, t := range types {
				types[i] = strings.ToLower(t)
			}
			firstType := types[0]
			switch {
			case slices.Contains(types, "article") && slices.Contains(types, "proceedings paper"):
				p.JournalArticleType = "proceedingsPaper"
			case firstType == "journal article" || firstType == "article" || firstType == "journal paper":
				p.JournalArticleType = "original"
			case firstType == "review":
				p.JournalArticleType = "review"
			case firstType == "letter" || firstType == "note" || firstType == "letter/note":
				p.JournalArticleType = "letterNote"
			case firstType == "book":
				p.Type = "book"
			case firstType == "book chapter":
				p.Type = "book_chapter"
			case firstType == "meeting abstract":
				p.Type = "conference"
				p.ConferenceType = "abstract"
			case firstType == "conference proceeding" || firstType == "proceedings paper" || firstType == "conference paper":
				p.Type = "conference"
				p.ConferenceType = "proceedingsPaper"
			case firstType == "poster":
				p.Type = "conference"
				p.ConferenceType = "poster"
			case firstType == "art exhibit review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "exhibitionReview"
			case firstType == "book review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "bookReview"
			case firstType == "dance performance review" || firstType == "theatre review" || firstType == "theater review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "theatreReview"
			case firstType == "database review" || firstType == "hardware review" || firstType == "software review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "productReview"
			case firstType == "editorial material" || firstType == "editorial":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "editorialMaterial"
			case firstType == "fiction" || firstType == "creative prose" || firstType == "poetry" || firstType == "script":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "artisticWork"
			case firstType == "film review" || firstType == "tv review" || firstType == "radio review" || firstType == "video review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "filmReview"
			case firstType == "music score review" || firstType == "music performance review" || firstType == "record review":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "musicReview"
			case firstType == "music score":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "musicEdition"
			case firstType == "news item":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "newsArticle"
			case firstType == "correction":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "correction"
			case firstType == "biographical-item" || firstType == "biographical item" || firstType == "item about an individual":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "biography"
			case firstType == "bibliography":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "bibliography"
			case firstType == "preprint":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "preprint"
			case firstType == "data paper":
				p.Type = "miscellaneous"
			case firstType == "other" || firstType == "discussion" || firstType == "slide":
				p.Type = "miscellaneous"
				p.MiscellaneousType = "other"
			}
		case "AF":
			for _, val := range v {
				p.Author = append(p.Author, extractContributor(val))
			}
		case "AU":
			// give preference to AF over AU
			if _, ok := r["AF"]; !ok {
				for _, val := range v {
					p.Author = append(p.Author, extractContributor(val))
				}
			}
		case "BF", "ED":
			for _, val := range v {
				p.Editor = append(p.Editor, extractContributor(val))
			}
		case "BE":
			// give preference to BF over BE
			if _, ok := r["BF"]; !ok {
				for _, val := range v {
					p.Editor = append(p.Editor, extractContributor(val))
				}
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

func extractContributor(val string) *models.Contributor {
	nameParts := reSplit.Split(val, -1)

	lastName := strings.TrimSpace(nameParts[0])
	firstName := "[missing]" // TODO
	if len(nameParts) > 1 {
		firstName = strings.TrimSpace(nameParts[1])
	}

	return models.ContributorFromFirstLastName(firstName, lastName)
}
