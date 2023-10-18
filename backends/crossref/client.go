package crossref

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caltechlibrary/doitools"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/validation"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"golang.org/x/text/language"
)

type Client struct {
	url  string
	http *http.Client
}

func New() *Client {
	return &Client{
		url: "https://api.crossref.org/works/",
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetPublication(id string) (*models.Publication, error) {
	doi, err := doitools.NormalizeDOI(id)
	if err != nil {
		return nil, err
	}

	// log.Printf("import publication doi: %s", doi)

	u, _ := url.Parse(c.url + url.PathEscape(doi))
	q := u.Query()
	// TODO remove hardcoded email
	q.Set("mailto", "bib-infra@lists.ugent.be")
	u.RawQuery = q.Encode()
	// log.Printf("import publication url: %s", u.String())
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	// log.Printf("%+v", res)
	defer res.Body.Close()
	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't import publication: %s", src)
	}

	// log.Printf("import publication src: %s", src)

	attrs := gjson.ParseBytes(src).Get("message")

	p := &models.Publication{
		Type:              "miscellaneous",
		PublicationStatus: "published",
	}

	if res := attrs.Get("type"); res.Exists() {
		switch res.String() {
		case "book-chapter", "book-section", "book-track", "book-part", "other":
			p.Type = "book_chapter"
		case "book", "book-set", "monograph", "reference-book":
			p.Type = "book"
		case "book-series", "edited-book", "standard-series":
			p.Type = "book_editor"
		case "journal-article", "reference-entry":
			p.Type = "journal_article"
		case "journal-volume", "journal-issue":
			p.Type = "issue_editor"
		case "proceedings", "proceedings-article":
			p.Type = "conference"
		case "dissertation":
			p.Type = "dissertation"
		case "preprint", "posted-content":
			p.MiscellaneousType = "preprint"
		case "report-series":
			p.MiscellaneousType = "report"
		case "standard":
			p.MiscellaneousType = "technicalStandard"
		}
	}
	if res := attrs.Get("DOI"); res.Exists() {
		p.DOI = res.String()
	}
	if res := attrs.Get("issued.date-parts.0.0"); res.Exists() {
		p.Year = res.String()
	} else if res := attrs.Get("created.date-parts.0.0"); res.Exists() {
		p.Year = res.String()
	}
	if res := attrs.Get("title.0"); res.Exists() {
		p.Title = res.String()
	}
	if res := attrs.Get("subtitle"); res.Exists() {
		for _, r := range res.Array() {
			p.AlternativeTitle = append(p.AlternativeTitle, r.String())
		}
	}
	if res := attrs.Get("isbn-type"); res.IsArray() {
		for _, r := range res.Array() {
			switch r.Get("type").String() {
			case "print":
				p.ISBN = append(p.ISBN, r.Get("value").String())
			case "electronic":
				p.EISBN = append(p.EISBN, r.Get("value").String())
			}
		}
	} else if res := attrs.Get("ISBN"); res.Exists() { // TODO: keep fallback?
		for _, r := range res.Array() {
			p.ISBN = append(p.ISBN, r.String())
		}
	}
	if res := attrs.Get("issn-type"); res.IsArray() {
		for _, r := range res.Array() {
			switch r.Get("type").String() {
			case "print":
				p.ISSN = append(p.ISSN, r.Get("value").String())
			case "electronic":
				p.EISSN = append(p.EISSN, r.Get("value").String())
			}
		}
	} else if res := attrs.Get("ISSN"); res.Exists() { // TODO: keep fallback?
		for _, r := range res.Array() {
			p.ISSN = append(p.ISSN, r.String())
		}
	}
	if res := attrs.Get("author"); res.Exists() {
		for _, r := range res.Array() {
			name := r.Get("name").String()
			firstName := r.Get("given").String()
			lastName := r.Get("family").String()
			if firstName == "" {
				firstName = "[missing]" // TODO
			}
			if lastName == "" {
				lastName = name
			}
			p.Author = append(p.Author, models.ContributorFromFirstLastName(firstName, lastName))
		}
	}
	if res := attrs.Get("editor"); res.Exists() {
		for _, r := range res.Array() {
			name := r.Get("name").String()
			firstName := r.Get("given").String()
			lastName := r.Get("family").String()
			if firstName == "" {
				firstName = "[missing]" // TODO
			}
			if lastName == "" {
				lastName = name
			}
			p.Editor = append(p.Editor, models.ContributorFromFirstLastName(firstName, lastName))
		}
	}
	if res := attrs.Get("subject"); res.Exists() {
		for _, r := range res.Array() {
			p.Keyword = append(p.Keyword, r.String())
		}
	}
	if res := attrs.Get("abstract"); res.Exists() {
		p.AddAbstract(&models.Text{
			Text: res.String(),
			Lang: "und",
		})
	}
	if res := attrs.Get("language"); res.Exists() {
		if base, err := language.ParseBase(res.String()); err == nil {
			if validation.InArray(vocabularies.Map["language_codes"], base.ISO3()) {
				p.Language = []string{base.ISO3()}
			}
		}
	}
	if res := attrs.Get("volume"); res.Exists() {
		p.Volume = res.String()
	}
	if res := attrs.Get("issue"); res.Exists() {
		p.Issue = res.String()
	}
	if res := attrs.Get("article-number"); res.Exists() {
		p.ArticleNumber = res.String()
	}
	if res := attrs.Get("publisher"); res.Exists() {
		p.Publisher = res.String()
	}
	if res := attrs.Get("publisher-location"); res.Exists() {
		p.PlaceOfPublication = res.String()
	}
	if res := attrs.Get("page"); res.Exists() {
		pages := strings.Split(res.String(), "-")
		p.PageFirst = pages[0]
		if len(pages) > 1 {
			p.PageLast = pages[1]
		}
	}
	switch p.Type {
	case "book":
		if res := attrs.Get("container-title.0"); res.Exists() {
			p.SeriesTitle = res.String()
		}
	case "book_chapter":
		if res := attrs.Get("container-title.0"); res.Exists() {
			p.Publication = res.String()
		}
		if res := attrs.Get("container-title.1"); res.Exists() {
			p.SeriesTitle = res.String()
		}
	case "journal_article", "conference":
		if res := attrs.Get("container-title.0"); res.Exists() {
			p.Publication = res.String()
		}
	}
	if res := attrs.Get("short-container-title.0"); res.Exists() {
		p.PublicationAbbreviation = res.String()
	}

	// log.Printf("import publication: %+v", p)

	return p, nil
}
