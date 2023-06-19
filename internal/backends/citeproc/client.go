package citeproc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type Client struct {
	style string
	url   string
	http  *http.Client
}

func New(baseURL, style string) *Client {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("style", style)
	u.RawQuery = q.Encode()
	return &Client{
		style: style,
		url:   u.String(),
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) EncodePublication(p *models.Publication) ([]byte, error) {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(&RequestBody{Items: []Item{publicationToItem(p)}})

	req, err := http.NewRequest(http.MethodPost, c.url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
		return nil, fmt.Errorf("can't generate citation: %s", src)
	}

	// log.Printf("citation src: %s", src)

	cite := gjson.GetBytes(src, "bibliography.1.0").String()

	return []byte(strings.TrimSpace(cite)), nil
}

type RequestBody struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID             string   `json:"id"`
	Type           string   `json:"type,omitempty"`
	Title          string   `json:"title,omitempty"`
	Author         []Person `json:"author,omitempty"`
	Edition        string   `json:"edition,omitempty"`
	Issued         Issued   `json:"issued,omitempty"`
	Publisher      string   `json:"publisher,omitempty"`
	PublisherPlace string   `json:"publisher-place,omitempty"`
	DOI            string   `json:"DOI,omitempty"`
	ISBN           string   `json:"ISBN,omitempty"`
}

type Issued struct {
	Raw string `json:"raw,omitempty"`
}

type Person struct {
	Family string `json:"family,omitempty"`
	Given  string `json:"given,omitempty"`
}

func publicationToItem(p *models.Publication) Item {
	item := Item{
		ID:             p.ID,
		Title:          p.Title,
		Edition:        p.Edition,
		Publisher:      p.Publisher,
		PublisherPlace: p.PlaceOfPublication,
		DOI:            p.DOI,
	}

	switch p.Type {
	case "book":
		item.Type = "book"
	case "journal_article":
		item.Type = "article-journal"
	case "chapter", "book_chapter":
		item.Type = "chapter"
	case "dissertation":
		item.Type = "thesis"
	}
	item.Issued.Raw = p.Year
	for _, a := range p.Author {
		item.Author = append(item.Author, Person{Family: a.LastName(), Given: a.FirstName()})
	}
	if len(p.ISBN) > 0 {
		item.ISBN = p.ISBN[0]
	}

	return item
}
