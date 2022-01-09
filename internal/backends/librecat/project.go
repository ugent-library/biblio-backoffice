package librecat

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) SuggestProjects(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if _, err := c.get("/completion/project", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (c *Client) GetProject(id string) (*models.Project, error) {
	project := &models.Project{}
	if _, err := c.get(fmt.Sprintf("/project/%s", id), nil, project); err != nil {
		return nil, err
	}
	return project, nil
}
