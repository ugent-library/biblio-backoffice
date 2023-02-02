package biblio

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func (c *Client) GetPerson(id string) (*models.Person, error) {
	p := &models.Person{}
	qp := url.Values{}
	qp.Set("format", "json")
	if err := c.get(fmt.Sprintf("/person/%s", url.PathEscape(id)), qp, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Client) SuggestPeople(q string) ([]models.Person, error) {
	hits := make([]models.Person, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if err := c.get("/person/completion", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}
