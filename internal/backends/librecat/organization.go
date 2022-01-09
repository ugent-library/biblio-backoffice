package librecat

import (
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) SuggestOrganizations(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if _, err := c.get("/completion/organization", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}
