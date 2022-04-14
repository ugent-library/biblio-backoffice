package datacite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/caltechlibrary/doitools"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backend/internal/models"
	"golang.org/x/text/language"
)

const ContentType = "application/vnd.datacite.datacite+json"

var whitespace = regexp.MustCompile(`\s*,\s*`)

type Client struct {
	url  string
	http *http.Client
}

func New() *Client {
	return &Client{
		url: "https://api.datacite.org/dois/",
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetDataset(id string) (*models.Dataset, error) {
	doi, err := doitools.NormalizeDOI(id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, c.url+url.PathEscape(doi), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ContentType)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	src, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't import dataset: %s", src)
	}

	// log.Printf("import dataset src: %s", src)

	attrs := gjson.ParseBytes(src)

	d := &models.Dataset{}

	if res := attrs.Get("doi"); res.Exists() {
		d.DOI = res.String()
	}
	if res := attrs.Get("publicationYear"); res.Exists() {
		d.Year = res.String()
	}
	if res := attrs.Get("titles.0.title"); res.Exists() {
		d.Title = res.String()
	}
	if res := attrs.Get("publisher"); res.Exists() {
		d.Publisher = res.String()
	}
	if res := attrs.Get("formats"); res.Exists() {
		for _, r := range res.Array() {
			d.Format = append(d.Format, r.String())
		}
	}
	if res := attrs.Get("subjects.#.subject"); res.Exists() {
		for _, r := range res.Array() {
			keywords := whitespace.Split(r.String(), -1)
			d.Keyword = append(d.Keyword, keywords...)
		}
	}
	if res := attrs.Get("creators"); res.Exists() {
		for _, r := range res.Array() {
			c := models.Contributor{}
			if res := r.Get("name"); res.Exists() {
				c.FullName = res.String()
			}
			if res := r.Get("givenName"); res.Exists() {
				c.FirstName = res.String()
			}
			if res := r.Get("familyName"); res.Exists() {
				c.LastName = res.String()
			}
			d.Author = append(d.Author, &c)
		}
	}
	if res := attrs.Get("descriptions"); res.Exists() {
		for _, r := range res.Array() {
			t := models.Text{Text: r.Get("description").String(), Lang: "und"}
			if res := r.Get("lang"); res.Exists() {
				if tag, err := language.Parse(res.String()); err == nil {
					t.Lang = tag.String()
				}

			}
			d.Abstract = append(d.Abstract, t)
		}
	}
	if res := attrs.Get(`rightsList.#(rightsIdentifierScheme="SPDX").rightsIdentifier`); res.Exists() {
		d.License = res.String()
	}
	if res := attrs.Get(`rightsList.#(rightsUri%"info:eu-repo/semantics/*").rightsUri`); res.Exists() {
		d.AccessLevel = res.String()
	}

	// log.Printf("import dataset: %+v", d)

	return d, nil
}
