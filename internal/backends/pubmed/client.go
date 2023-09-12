package pubmed

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

// api reference: https://europepmc.org/RestfulWebService

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
	q.Set("resultType", "core")
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
	src, err := io.ReadAll(res.Body)
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
	if journalInfo := attrs.Get("journalInfo"); journalInfo.IsObject() {
		if journalVolume := journalInfo.Get("volume"); journalVolume.Exists() {
			p.Volume = journalVolume.String()
		}
		if journalIssue := journalInfo.Get("issue"); journalIssue.Exists() {
			p.Issue = journalIssue.String()
		}
		if journal := journalInfo.Get("journal"); journal.IsObject() {
			if journalTitle := journal.Get("title"); journalTitle.Exists() {
				p.Publication = journalTitle.String()
			}
			if journalAbbr := journal.Get("isoabbreviation"); journalAbbr.Exists() {
				p.PublicationAbbreviation = journalAbbr.String()
			}
			if journalISSN := journal.Get("issn"); journalISSN.Exists() {
				p.ISSN = append(p.ISSN, journalISSN.String())
			}
			if journalEISSN := journal.Get("essn"); journalEISSN.Exists() {
				p.EISSN = append(p.EISSN, journalEISSN.String())
			}
		}
	}
	if res := attrs.Get("pubYear"); res.Exists() {
		p.Year = res.String()
	}
	if authorList := attrs.Get("authorList.author"); authorList.IsArray() {
		for _, author := range authorList.Array() {
			firstName := author.Get("firstName").String()
			if firstName == "" {
				firstName = "[missing]"
			}
			lastName := author.Get("lastName").String()
			if lastName == "" {
				lastName = "[missing]"
			}
			c := models.ContributorFromFirstLastName(firstName, lastName)
			p.Author = append(p.Author, c)
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
	// TODO: language of abstract always in English?
	if res := attrs.Get("abstractText"); res.Exists() {
		p.AddAbstract(&models.Text{Text: res.String(), Lang: "eng"})
	}
	if res := attrs.Get("language"); res.Exists() {
		p.Language = append(p.Language, res.String())
	}

	// log.Printf("import publication: %+v", p)

	return p, nil
}
