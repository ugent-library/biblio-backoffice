package peopleservice

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/people-service/api/v1"
	pmodels "github.com/ugent-library/people-service/models"
)

type Config struct {
	APIUrl string
	APIKey string
}

type Client struct {
	config Config
	client *api.Client
}

type securitySource struct {
	apiKey string
}

func (s *securitySource) ApiKey(ctx context.Context, operationName string) (api.ApiKey, error) {
	return api.ApiKey{APIKey: s.apiKey}, nil
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

func (c *Client) GetPerson(biblioID string) (*models.Person, error) {
	if biblioID == "" {
		return nil, models.ErrNotFound
	}

	ctx := context.TODO()

	res, err := c.client.GetPeopleByIdentifier(ctx, &api.GetPeopleByIdentifierRequest{
		Identifier: []string{"urn:biblio_id:" + biblioID},
	})
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, models.ErrNotFound
	}

	people, err := c.mapPeople(ctx, res.Data[0])
	if err != nil {
		return nil, err
	}
	return people[0], nil
}

func (c *Client) SuggestPeople(query string) ([]*models.Person, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*models.Person{}, nil
	}

	ctx := context.TODO()
	limit := 20

	res, err := c.client.SuggestPeople(ctx, &api.SuggestPeopleRequest{
		Limit:  api.NewOptInt(limit),
		Active: []bool{true, false},
		Query:  query,
	})
	if err != nil {
		return nil, err
	}

	people, err := c.mapPeople(ctx, res.Data...)
	if err != nil {
		return nil, err
	}
	return people, nil
}

func (c *Client) GetUser(biblioID string) (*models.Person, error) {
	if biblioID == "" {
		return nil, models.ErrNotFound
	}

	ctx := context.TODO()

	res, err := c.client.GetPeopleByIdentifier(ctx, &api.GetPeopleByIdentifierRequest{
		Identifier: []string{"urn:biblio_id:" + biblioID},
	})
	if err != nil {
		return nil, err
	}

	res.Data = lo.Filter(res.Data, func(ap api.Person, idx int) bool {
		return ap.Active.Value
	})
	if len(res.Data) == 0 {
		return nil, models.ErrNotFound
	}

	people, err := c.mapPeople(ctx, res.Data[0])
	if err != nil {
		return nil, err
	}
	return people[0], nil
}

func (c *Client) GetUserByUsername(username string) (*models.Person, error) {
	if username == "" {
		return nil, models.ErrNotFound
	}

	ctx := context.TODO()
	res, err := c.client.GetPeopleByIdentifier(ctx, &api.GetPeopleByIdentifierRequest{
		Identifier: []string{"urn:ugent_username:" + username},
	})
	if err != nil {
		return nil, err
	}
	res.Data = lo.Filter(res.Data, func(ap api.Person, idx int) bool {
		return ap.Active.Value
	})
	if len(res.Data) == 0 {
		return nil, models.ErrNotFound
	}
	people, err := c.mapPeople(ctx, res.Data[0])
	if err != nil {
		return nil, err
	}
	return people[0], nil
}

func (c *Client) SuggestUsers(query string) ([]*models.Person, error) {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return []*models.Person{}, nil
	}

	ctx := context.TODO()
	limit := 25

	res, err := c.client.SuggestPeople(ctx, &api.SuggestPeopleRequest{
		Limit:  api.NewOptInt(limit),
		Active: []bool{true},
		Query:  query,
	})
	if err != nil {
		return nil, err
	}

	people, err := c.mapPeople(ctx, res.Data...)
	if err != nil {
		return nil, err
	}
	return people, nil
}

func (c *Client) mapPeople(ctx context.Context, apiPeople ...api.Person) ([]*models.Person, error) {
	orgExternalIDs := []string{}
	for _, ap := range apiPeople {
		for _, orgMember := range ap.Organization {
			orgExternalIDs = append(orgExternalIDs, orgMember.ID)
		}
	}
	orgExternalIDs = lo.Uniq(orgExternalIDs)
	res, err := c.client.GetOrganizationsById(ctx, &api.GetOrganizationsByIdRequest{
		ID: orgExternalIDs,
	})
	if err != nil {
		return nil, err
	}
	externalOrgs := map[string]*models.Organization{}
	for _, ao := range res.Data {
		org, err := c.mapOrganization(ctx, &ao)
		if err != nil {
			return nil, err
		}
		externalOrgs[ao.ID.Value] = org
	}

	people := make([]*models.Person, 0, len(apiPeople))
	for _, ap := range apiPeople {
		p := &models.Person{
			Active:      ap.Active.Value,
			DateCreated: &ap.DateCreated.Value,
			DateUpdated: &ap.DateUpdated.Value,
			Email:       ap.Email.Value,
			FirstName:   ap.GivenName.Value,
			FullName:    ap.Name.Value,
			LastName:    ap.FamilyName.Value,
		}
		for _, role := range ap.Role {
			if role == "biblio-admin" {
				p.Role = "admin"
			}
		}
		for _, id := range ap.Identifier {
			urn, _ := pmodels.ParseURN(id)
			switch urn.Namespace {
			case "orcid":
				p.ORCID = urn.Value
			case "biblio_id":
				p.ID = urn.Value
			case "historic_ugent_id":
				p.UGentID = append(p.UGentID, urn.Value)
			case "ugent_username":
				p.Username = urn.Value
			}
		}
		if tokens, ok := ap.Token.Get(); ok {
			if orcidToken, ok := tokens["orcid"]; ok {
				p.ORCIDToken = orcidToken
			}
		}
		for _, orgMember := range ap.Organization {
			if org, ok := externalOrgs[orgMember.ID]; ok {
				p.Affiliations = append(p.Affiliations, &models.Affiliation{
					OrganizationID: org.ID,
					Organization:   org,
				})
			}
		}
		people = append(people, p)
	}

	return people, nil
}

func (c *Client) GetOrganization(biblioID string) (*models.Organization, error) {
	if biblioID == "" {
		return nil, models.ErrNotFound
	}

	ctx := context.TODO()

	res, err := c.client.GetOrganizationsByIdentifier(ctx, &api.GetOrganizationsByIdentifierRequest{
		Identifier: []string{"urn:biblio_id:" + biblioID},
	})
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, models.ErrNotFound
	}

	return c.mapOrganization(ctx, &res.Data[0])
}

func (c *Client) SuggestOrganizations(query string) ([]models.Completion, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []models.Completion{}, nil
	}

	ctx := context.TODO()
	limit := 20
	completions := make([]models.Completion, 0, limit)

	res, err := c.client.SuggestOrganizations(ctx, &api.SuggestOrganizationsRequest{
		Limit: api.NewOptInt(limit),
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	for _, apiOrg := range res.Data {
		completion := models.Completion{
			Heading: apiOrg.NameEng.Value,
		}
		for _, id := range apiOrg.Identifier {
			urn, _ := pmodels.ParseURN(id)
			if urn.Namespace == "biblio_id" {
				completion.ID = urn.Value
			}
		}
		completions = append(completions, completion)
	}

	return completions, nil
}

func (c *Client) mapOrganization(ctx context.Context, ao *api.Organization) (*models.Organization, error) {
	o := &models.Organization{
		Name: ao.NameEng.Value,
	}
	for _, id := range ao.Identifier {
		urn, _ := pmodels.ParseURN(id)
		if urn.Namespace == "biblio_id" {
			o.ID = urn.Value
			break
		}
	}

	o.Tree = append(o.Tree, models.OrganizationTreeElement{ID: o.ID})

	// IMPORTANT: only one parent processed now
	if len(ao.Parent) > 0 {
		parentID := ao.Parent[0].ID
		for parentID != "" {
			apiParentOrg, err := c.client.GetOrganization(ctx, &api.GetOrganizationRequest{
				ID: parentID,
			})
			if err != nil {
				return nil, err
			}
			parentBiblioID := ""
			for _, id := range apiParentOrg.Identifier {
				urn, _ := pmodels.ParseURN(id)
				if urn.Namespace == "biblio_id" {
					parentBiblioID = urn.Value
				}
			}
			o.Tree = append([]models.OrganizationTreeElement{{ID: parentBiblioID}}, o.Tree...)
			if len(apiParentOrg.Parent) > 0 {
				parentID = apiParentOrg.Parent[0].ID
			} else {
				parentID = ""
			}
		}
	}

	return o, nil
}
