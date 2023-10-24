package plato

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/recordsources"
)

func init() {
	recordsources.Register("plato", New)
}

func New(conn string) (recordsources.Source, error) {
	return &platoSource{}, nil
}

type platoSource struct {
}

func (s *platoSource) GetRecords(ctx context.Context) ([]*models.CandidateRecord, error) {
	u := "https://plato.ea.ugent.be/service/dr/2biblio.jsp"
	c := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var recs []*models.CandidateRecord

	gjson.GetBytes(body, "list").ForEach(func(key, val gjson.Result) bool {
		p := &models.Publication{}

		if v := val.Get("titel.eng"); v.Exists() {
			p.Title = v.String()
		}

		j, _ := json.Marshal(p) // TODO handle error
		recs = append(recs, &models.CandidateRecord{
			SourceName: "plato",
			SourceID:   uuid.NewString(), // TODO
			Type:       "Publication",
			Metadata:   j,
		})
		return true
	})

	return recs, nil
}
