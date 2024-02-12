package plato

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (s *platoSource) GetRecords(ctx context.Context, cb func(recordsources.Record) error) error {
	c := http.Client{
		Timeout: time.Second * 2,
	}

	baseURL, err := url.ParseRequestURI(s.url)
	if err != nil {
		return err
	}

	const count = 50
	from := 1

	for {
		u := *baseURL
		q := u.Query()
		q.Set("from", fmt.Sprintf("%d", from))
		q.Set("count", fmt.Sprintf("%d", count))
		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return err
		}
		res, err := c.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		listSize := 0
		var cbErr error
		gjson.GetBytes(body, "list").ForEach(func(key, val gjson.Result) bool {
			err = cb(recordsources.Record{
				SourceName:     "plato",
				SourceID:       val.Get("plato_id").String(),
				SourceMetadata: []byte(val.Raw),
			})
			if err != nil {
				cbErr = err
				return false
			}
			listSize++
			return true
		})
		if cbErr != nil {
			return cbErr
		}

		if listSize < count {
			break
		}

		from += count
	}

	return nil
}
