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
		// Don't reuse connections between pages: processing a batch can take
		// long enough for an idle keep-alive connection to be dropped by an
		// intermediary, after which reusing it stalls until Client.Timeout
		// fires ("awaiting headers"). A fresh connection per request avoids this.
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
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

		body, err := s.fetchPage(ctx, &c, u.String())
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

// fetchPage requests a single page, retrying a few times on transient failures
// (network/timeout errors and 5xx responses) with a short backoff. It fails
// fast on context cancellation and permanent errors such as 4xx responses.
func (s *platoSource) fetchPage(ctx context.Context, c *http.Client, u string) ([]byte, error) {
	const maxAttempts = 3
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt-1) * time.Second):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(s.username, s.password)

		res, err := c.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = err
			continue
		}

		if res.StatusCode >= 500 {
			res.Body.Close()
			lastErr = fmt.Errorf("GET %q: %s", u, res.Status)
			continue
		}
		if res.StatusCode < 200 || res.StatusCode >= 400 {
			res.Body.Close()
			return nil, fmt.Errorf("GET %q: %s", u, res.Status)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = err
			continue
		}

		return body, nil
	}

	return nil, lastErr
}
