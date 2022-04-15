package pubmed

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backend/internal/models"
)

var reSplit = regexp.MustCompile(`\s*[,;]\s*`)

type Client struct {
	url  string
	http *http.Client
}

func New() *Client {
	return &Client{
		url: "https://www.ebi.ac.uk/europepmc/webservices/rest/search",
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetPublication(id string) (*models.Publication, error) {
	// log.Printf("import publication pubmed: %s", id)

	u, _ := url.Parse(c.url)
	q := u.Query()
	q.Set("format", "json")
	q.Set("query", id)
	u.RawQuery = q.Encode()
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
	src, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't import publication: %s", src)
	}

	// log.Printf("import publication src: %s", src)

	p := &models.Publication{
		Type: "journal_article",
	}

	attrs := gjson.ParseBytes(src)

	if attrs.Get("hitCount").Int() != 1 {
		return nil, fmt.Errorf("no publication found")
	}

	attrs = attrs.Get("resultList.result.0")

	if res := attrs.Get("pmid"); res.Exists() {
		p.PubMedID = res.String()
	}
	if res := attrs.Get("doi"); res.Exists() {
		p.DOI = res.String()
	}
	if res := attrs.Get("title"); res.Exists() {
		p.Title = res.String()
	}
	if res := attrs.Get("journalTitle"); res.Exists() {
		p.Publication = res.String()
	}
	if res := attrs.Get("journalVolume"); res.Exists() {
		p.Volume = res.String()
	}
	if res := attrs.Get("journalIssn"); res.Exists() {
		p.ISSN = reSplit.Split(res.String(), -1)
	}
	if res := attrs.Get("pubYear"); res.Exists() {
		p.Year = res.String()
	}
	if res := attrs.Get("authorString"); res.Exists() {
		for _, r := range reSplit.Split(res.String(), -1) {
			nameParts := strings.Split(strings.ReplaceAll(r, ".", ""), " ")
			firstName := nameParts[len(nameParts)-1]
			lastName := strings.Join(nameParts[:len(nameParts)-1], " ")
			c := models.Contributor{
				FirstName: firstName,
				LastName:  lastName,
				FullName:  firstName + " " + lastName,
			}
			p.Author = append(p.Author, &c)
		}
	}
	if res := attrs.Get("pageInfo"); res.Exists() {
		pages := strings.Split(res.String(), "-")
		if len(pages) > 1 {
			p.PageFirst = pages[0]
			p.PageLast = pages[1]
		} else {
			p.ArticleNumber = pages[0]
		}
	}

	// log.Printf("import publication: %+v", p)

	return p, nil
}
