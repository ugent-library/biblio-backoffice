package plato

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/recordsources"
)

const limit = 50

func init() {
	recordsources.Register("plato", NewSource)
}

type Config struct {
	URL      string `env:"URL"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

func NewSource() (recordsources.Source, error) {
	c := &Config{}
	env.ParseWithOptions(c, env.Options{
		Prefix: "BIBLIO_BACKOFFICE_PLATO_",
	})

	return &platoSource{
		url:      c.URL,
		username: c.Username,
		password: c.Password,
	}, nil
}

type platoSource struct {
	url      string
	username string
	password string
}

func (s *platoSource) GetRecords(ctx context.Context, cb func(recordsources.Record) error) error {
	c := http.Client{
		Timeout: 30 * time.Second,
	}

	baseURL, err := url.ParseRequestURI(s.url)
	if err != nil {
		return fmt.Errorf("plato: %w", err)
	}

	for from := 1; ; from += limit {
		u := *baseURL
		q := u.Query()
		q.Set("from", fmt.Sprint(from))
		q.Set("count", fmt.Sprint(limit))
		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return fmt.Errorf("plato: %w", err)
		}

		req.SetBasicAuth(s.username, s.password)

		res, err := c.Do(req)
		if err != nil {
			return fmt.Errorf("plato: %w", err)
		}

		if res.StatusCode < 200 || res.StatusCode >= 400 {
			return fmt.Errorf("plato: GET %q: %s", u.String(), res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("plato: %w", err)
		}

		recs := gjson.GetBytes(body, "list").Array()

		for _, rec := range recs {
			err = cb(NewRecord(rec.Get("plato_id").String(), []byte(rec.Raw)))
			if err != nil {
				return fmt.Errorf("plato: %w", err)
			}
		}

		if len(recs) < limit {
			break
		}
	}

	return nil
}
