package engine

import (
	"fmt"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) SuggestProjects(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if _, err := e.get("/completion/project", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) GetProject(id string) (*models.Project, error) {
	project := &models.Project{}
	if _, err := e.get(fmt.Sprintf("/project/%s", id), nil, project); err != nil {
		return nil, err
	}
	return project, nil
}
