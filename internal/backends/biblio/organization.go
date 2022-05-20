package biblio

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetOrganization(id string) (*models.Organization, error) {
	org := &models.Organization{}
	qp := url.Values{}
	qp.Set("format", "json")
	if _, err := c.get(fmt.Sprintf("/organization/%s", url.PathEscape(id)), qp, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (c *Client) SuggestOrganizations(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if _, err := c.get("/organization/completion", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}
