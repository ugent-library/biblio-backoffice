package plato

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/recordsources"
)

func init() {
	recordsources.Register("plato", New)
}

func New(conn string) (recordsources.Source, error) {
	return &platoSource{
		url: conn,
	}, nil
}

type platoSource struct {
	url string
}

func (s *platoSource) GetRecords(ctx context.Context) ([]recordsources.Record, error) {
	c := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
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

	var recs []recordsources.Record

	gjson.GetBytes(body, "list").ForEach(func(key, val gjson.Result) bool {
		recs = append(recs, recordsources.Record{
			SourceName:     "plato",
			SourceID:       val.Get("plato_id").String(),
			SourceMetadata: []byte(val.Raw),
		})
		return true
	})

	return recs, nil
}
