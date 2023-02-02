package biblio

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func (c *Client) SuggestProjects(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if err := c.get("/project/completion", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (c *Client) GetProject(id string) (*models.Project, error) {
	project := &models.Project{}
	qp := url.Values{}
	qp.Set("format", "json")
	if err := c.get(fmt.Sprintf("/project/%s", id), qp, project); err != nil {
		return nil, err
	}
	return project, nil
}
