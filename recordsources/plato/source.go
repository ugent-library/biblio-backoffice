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

func init() {
	recordsources.Register("plato", NewSource)
}

type Config struct {
	Plato struct {
		URL      string `env:"URL"`
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
	} `envPrefix:"PLATO_"`
}

func NewSource() (recordsources.Source, error) {
	c := &Config{}
	env.ParseWithOptions(c, env.Options{
		Prefix: "BIBLIO_BACKOFFICE_",
	})

	return &platoSource{
		url:      c.Plato.URL,
		username: c.Plato.Username,
		password: c.Plato.Password,
	}, nil
}

type platoSource struct {
	url      string
	username string
	password string
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

		req.SetBasicAuth(s.username, s.password)

		res, err := c.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		statusOK := res.StatusCode >= http.StatusOK && res.StatusCode <= http.StatusPermanentRedirect
		if !statusOK {
			var err error
			switch res.StatusCode {
			// 4xx
			case http.StatusBadRequest:
				err = fmt.Errorf("[plato] a bad http request was sent to the server: %s", res.Status)
			case http.StatusUnauthorized:
				err = fmt.Errorf("[plato] the client wasn't authorized by the server, check credentials: %s", res.Status)
			case http.StatusNotFound:
				err = fmt.Errorf("[plato] the API endpoint URL could not be found: %s", res.Status)
			case http.StatusForbidden:
				err = fmt.Errorf("[plato] authorized but access forbidden: %s", res.Status)
			// 5xx
			case http.StatusInternalServerError:
				err = fmt.Errorf("[plato] internal server error: %s", res.Status)
			case http.StatusBadGateway:
				err = fmt.Errorf("[plato] bad gateway error: %s", res.Status)
			case http.StatusServiceUnavailable:
				err = fmt.Errorf("[plato] service unavailable error: %s", res.Status)
			default:
				err = fmt.Errorf("[plato] server returned an error: %s", res.Status)
			}

			return err
		}

		listSize := 0
		var cbErr error
		gjson.GetBytes(body, "list").ForEach(func(key, val gjson.Result) bool {
			err = cb(NewRecord(val.Get("plato_id").String(), []byte(val.Raw)))
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
