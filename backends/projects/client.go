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
	ctx := context.Background()

	res, err := c.client.GetProject(ctx, &api.GetProjectRequest{ID: fmt.Sprintf("urn:iweto:%s", id)})
	if err != nil {
		return nil, err
	}

	switch gp := res.(type) {
	case *api.GetProject:
		p := &models.Project{
			EUProject: &models.EUProject{},
		}

		if idx := strings.LastIndex(gp.GetID(), ":"); idx != -1 {
			p.ID = gp.GetID()[idx+1:]
		} else {
			p.ID = ""
		}

		for _, v := range gp.GetName() {
			if v.GetLanguage() == "und" {
				p.Title = v.GetValue()
			}
		}

		if v, ok := gp.GetFoundingDate().Get(); ok {
			p.StartDate = v
		}

		if v, ok := gp.GetDissolutionDate().Get(); ok {
			p.EndDate = v
		}

		for _, v := range gp.GetIdentifier() {
			if v.GetPropertyID() == "CORDIS" {
				p.EUProject.ID = v.GetValue()
			}
		}

		if v, ok := gp.IsFundedBy.Get(); ok {
			if cid, ok := v.HasCallNumber.Get(); ok {
				p.EUProject.CallID = cid
			}

			if aw, ok := v.IsAwardedBy.Get(); ok {
				p.EUProject.FrameworkProgramme = aw.GetName()
			}
		}

		for _, v := range gp.GetIdentifier() {
			if v.GetPropertyID() == "GISMO" {
				p.GISMOID = v.GetValue()
			}
		}

		for _, v := range gp.GetIdentifier() {
			if v.GetPropertyID() == "IWETO" {
				p.IWETOID = v.GetValue()
			}
		}

		// iweto_id not filled in everywhere, but should be same as id for now
		if p.IWETOID == "" {
			p.IWETOID = p.ID
		}

		return p, nil
	case *api.ErrorStatusCode:
		return nil, fmt.Errorf("projects suggestions: project service has thrown an error: code: %d, id: %s", gp.GetStatusCode(), id)
	}

	return nil, fmt.Errorf("projects suggestions: something went wrong (id: %s)", id)
}

func (c *Client) SuggestProjects(q string) ([]models.Completion, error) {
	ctx := context.Background()

	res, err := c.client.SuggestProjects(ctx, &api.SuggestProjectsRequest{Query: q})
	if err != nil {
		return nil, err
	}

	completions := make([]models.Completion, 0)

	for _, hit := range res.Data {
		c := models.Completion{}

		if idx := strings.LastIndex(hit.GetID(), ":"); idx != -1 {
			c.ID = hit.GetID()[idx+1:]
		} else {
			c.ID = ""
		}

		for _, v := range hit.GetName() {
			if v.GetLanguage() == "und" {
				c.Heading = v.GetValue()
			}
		}

		acrs := hit.GetHasAcronym()
		if len(acrs) > 0 {
			c.Description = acrs[0]
		}

		completions = append(completions, c)
	}

	return completions, nil
}
