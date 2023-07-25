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
	"github.com/ugent-library/biblio-backoffice/internal/backends"
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
		return nil, fmt.Errorf("%w : %w", backends.ErrBaddConn, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, backends.ErrNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, backends.ErrInvalidContent
	}

	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !gjson.ValidBytes(src) {
		return nil, backends.ErrInvalidContent
	}

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
			name := r.Get("name").String()
			firstName := r.Get("givenName").String()
			lastName := r.Get("familyName").String()
			if firstName == "" {
				firstName = "[missing]" // TODO
			}
			if lastName == "" {
				lastName = name
			}
			d.Author = append(d.Author, models.ContributorFromFirstLastName(firstName, lastName))
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
