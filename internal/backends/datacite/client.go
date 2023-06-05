package datacite

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/caltechlibrary/doitools"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
	"golang.org/x/text/language"
)

const ContentType = "application/vnd.datacite.datacite+json"

var reSplit = regexp.MustCompile(`\s*,\s*`)

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
	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't import dataset: %s", src)
	}

	// log.Printf("import dataset src: %s", src)

	attrs := gjson.ParseBytes(src)

	d := &models.Dataset{}

	if res := attrs.Get("language"); res.Exists() {
		if base, err := language.ParseBase(res.String()); err == nil {
			if validation.InArray(vocabularies.Map["language_codes"], base.ISO3()) {
				d.Language = append(d.Language, base.ISO3())
			}
		}
	}
	if res := attrs.Get("doi"); res.Exists() {
		d.Identifiers = models.Identifiers{"DOI": []string{res.String()}}
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
	if res := attrs.Get("subjects"); res.Exists() {
		for _, r := range res.Array() {
			if r.Get("subjectScheme").Exists() {
				continue
			}
			if keywords := reSplit.Split(r.Get("subject").String(), -1); len(keywords) > 0 {
				d.Keyword = append(d.Keyword, keywords...)
			}
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
			if c.FirstName == "" {
				c.FirstName = "[missing]" // TODO
			}
			if c.LastName == "" {
				c.LastName = c.FullName
			}
			if c.FullName != "" && c.FirstName != "" && c.LastName != "" {
				d.Author = append(d.Author, &c)
			}
		}
	}
	if res := attrs.Get("descriptions"); res.Exists() {
		for _, r := range res.Array() {
			t := models.Text{Text: r.Get("description").String(), Lang: "und"}
			if res := r.Get("lang"); res.Exists() {
				if base, err := language.ParseBase(res.String()); err == nil {
					if validation.InArray(vocabularies.Map["language_codes"], base.ISO3()) {
						t.Lang = base.ISO3()
					}
				}
			}
			d.AddAbstract(&t)
		}
	}
	if res := attrs.Get(`rightsList.#(rightsIdentifierScheme="SPDX").rightsIdentifier`); res.Exists() {
		license := strings.ToUpper(res.String())

		// @todo Clean IsDatasetLicense() up.
		if validation.IsDatasetLicense(license) {
			d.License = license
		} else {
			d.License = "LicenseNotListed"
			d.OtherLicense = license
		}
	}
	if res := attrs.Get(`rightsList.#(rightsUri%"info:eu-repo/semantics/*").rightsUri`); res.Exists() {
		d.AccessLevel = res.String()
	}

	// log.Printf("import dataset: %+v", d)

	return d, nil
}
