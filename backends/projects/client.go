package projects

import (
	"context"
	"fmt"
	"strings"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/projects-service/api/v1"
)

type Config struct {
	APIUrl string
	APIKey string
}

type securitySource struct {
	apiKey string
}

func (s *securitySource) ApiKey(ctx context.Context, operationName string) (api.ApiKey, error) {
	return api.ApiKey{APIKey: s.apiKey}, nil
}

type Client struct {
	config Config
	client *api.Client
}

func New(c Config) (*Client, error) {
	client, err := api.NewClient(c.APIUrl, &securitySource{c.APIKey})
	if err != nil {
		return nil, err
	}

	return &Client{
		config: c,
		client: client,
	}, nil
}

func (c *Client) GetProject(id string) (*models.Project, error) {
	ctx := context.TODO()

	res, err := c.client.GetProject(ctx, &api.GetProjectRequest{ID: fmt.Sprintf("urn:iweto:%s", id)})
	if err != nil {
		return nil, err
	}

	switch ap := res.(type) {
	case *api.GetProject:
		return c.mapper(ap), nil
	case *api.ErrorStatusCode:
		return nil, fmt.Errorf("projects suggestions: project service has thrown an error: code: %d, id: %s", ap.GetStatusCode(), id)
	}

	return nil, fmt.Errorf("projects suggestions: something went wrong (id: %s)", id)
}

func (c *Client) SuggestProjects(q string) ([]models.Project, error) {
	ctx := context.Background()

	res, err := c.client.SuggestProjects(ctx, &api.SuggestProjectsRequest{Query: q})
	if err != nil {
		return nil, err
	}

	projects := make([]models.Project, len(res.Data))

	for i, hit := range res.Data {
		p := c.mapper(&hit)
		projects[i] = *p
	}

	return projects, nil
}

func (c *Client) mapper(ap *api.GetProject) *models.Project {
	p := &models.Project{}

	if idx := strings.LastIndex(ap.GetID(), ":"); idx != -1 {
		p.ID = ap.GetID()[idx+1:]
	} else {
		p.ID = ""
	}

	for _, v := range ap.GetName() {
		if v.GetLanguage() == "und" {
			p.Title = v.GetValue()
		}
	}

	for _, v := range ap.GetDescription() {
		if v.GetLanguage() == "und" {
			p.Description = v.GetValue()
		}
	}

	acrs := ap.GetHasAcronym()
	if len(acrs) > 0 {
		p.Acronym = acrs[0]
	}

	if v, ok := ap.GetFoundingDate().Get(); ok {
		p.StartDate = v
	}

	if v, ok := ap.GetDissolutionDate().Get(); ok {
		p.EndDate = v
	}

	for _, v := range ap.GetIdentifier() {
		if v.GetPropertyID() == "CORDIS" {
			if p.EUProject == nil {
				p.EUProject = &models.EUProject{}
			}
			p.EUProject.ID = v.GetValue()
		}
	}

	if v, ok := ap.IsFundedBy.Get(); ok {
		if p.EUProject == nil {
			p.EUProject = &models.EUProject{}
		}

		if cid, ok := v.HasCallNumber.Get(); ok {
			p.EUProject.CallID = cid
		}

		if aw, ok := v.IsAwardedBy.Get(); ok {
			p.EUProject.FrameworkProgramme = aw.GetName()
		}
	}

	for _, v := range ap.GetIdentifier() {
		if v.GetPropertyID() == "GISMO" {
			p.GISMOID = v.GetValue()
		}
	}

	for _, v := range ap.GetIdentifier() {
		if v.GetPropertyID() == "IWETO" {
			p.IWETOID = v.GetValue()
		}
	}

	// iweto_id not filled in everywhere, but should be same as id for now
	if p.IWETOID == "" {
		p.IWETOID = p.ID
	}

	return p
}
