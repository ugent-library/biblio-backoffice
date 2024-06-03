package arxiv

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backoffice/models"
)

var reNormalizeID = regexp.MustCompile(`(?i)^arxiv:`)

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	TotalResults int      `xml:"totalResults"`
	Entry        Entry    `xml:"entry"`
}

type Entry struct {
	Title      string `xml:"title"`
	Summary    string `xml:"summary"`
	Published  string `xml:"published"`
	DOI        string `xml:"doi"`
	JournalRef string `xml:"journal_ref"`
	Author     []struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Comment string `xml:"comment"`
}

type Client struct {
	url  string
	http *http.Client
}

func New() *Client {
	return &Client{
		url: "https://export.arxiv.org/api/query",
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetPublication(id string) (*models.Publication, error) {
	id = reNormalizeID.ReplaceAllString(id, "")

	u, _ := url.Parse(c.url)
	q := u.Query()
	q.Set("id_list", id)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("arxiv: request failed: %w", err)
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("arxiv: request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arxiv: request failed with status %d", res.StatusCode)
	}

	defer res.Body.Close()
	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("arxiv: reading response failed: %w", err)
	}

	feed := Feed{}

	if err := xml.Unmarshal(src, &feed); err != nil {
		return nil, fmt.Errorf("arxiv: unmarshalling response failed: %w", err)
	}

	if feed.TotalResults != 1 {
		return nil, fmt.Errorf("arxiv: expected 1 entry, but found %d", feed.TotalResults)
	}

	p := &models.Publication{
		Type:               "journal_article",
		PublicationStatus:  "unpublished",
		JournalArticleType: "original",
		ArxivID:            id,
		Title:              feed.Entry.Title,
		DOI:                feed.Entry.DOI,
		Publication:        feed.Entry.JournalRef,
		AdditionalInfo:     feed.Entry.Comment,
	}

	if feed.Entry.Summary != "" {
		p.AddAbstract(&models.Text{
			Text: feed.Entry.Summary,
			Lang: "und",
		})
	}

	if len(feed.Entry.Published) > 4 {
		p.Year = feed.Entry.Published[0:4]
	} else if len(feed.Entry.Published) == 4 {
		p.Year = feed.Entry.Published
	}

	for _, a := range feed.Entry.Author {
		nameParts := strings.Split(a.Name, " ")
		firstName := nameParts[0]
		lastName := nameParts[0]
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}
		p.Author = append(p.Author, models.ContributorFromFirstLastName(firstName, lastName))
	}

	return p, nil
}
