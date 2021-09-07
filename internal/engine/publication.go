package engine

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/ugent-library/go-web/forms"
)

type PublicationHits struct {
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	LastPage     int           `json:"last_page"`
	PreviousPage bool          `json:"previous_page"`
	NextPage     bool          `json:"next_page"`
	Hits         []Publication `json:"hits"`
}

type Publication struct {
	*gabs.Container
}

func (p *Publication) UnmarshalJSON(b []byte) error {
	c, err := gabs.ParseJSON(b)
	if err != nil {
		return err
	}

	p.Container = c

	return nil
}

func (p *Publication) Data() map[string]interface{} {
	return p.Container.Data().(map[string]interface{})
}

func (p *Publication) IsOpenAccess() bool {
	for _, accessLevel := range p.Container.S("file.*.access_level").Children() {
		if accessLevel.Data().(string) == "open_access" {
			return true
		}
	}
	return false
}

func (e *Engine) UserPublications(userID string, args *SearchArgs) (*PublicationHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &PublicationHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/publication", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) GetPublication(id string) (*Publication, error) {
	pub := &Publication{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, pub); err != nil {
		return nil, err
	}
	return pub, nil
}
