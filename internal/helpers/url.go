package helpers

import (
	"html/template"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
)

var formEncoder = form.NewEncoder()

func URL(r *mux.Router) template.FuncMap {
	return template.FuncMap{
		"urlFor":     urlFor(r),
		"urlPathFor": urlPathFor(r),
		"urlSet":     urlSet,
		"urlQuery":   urlQuery,
	}
}

func urlFor(r *mux.Router) func(string, ...string) (*url.URL, error) {
	return func(name string, vars ...string) (*url.URL, error) {
		return r.Get(name).URL(vars...)
	}
}

func urlPathFor(r *mux.Router) func(string, ...string) (*url.URL, error) {
	return func(name string, vars ...string) (*url.URL, error) {
		return r.Get(name).URLPath(vars...)
	}
}

func urlSet(k, v interface{}, u *url.URL) (*url.URL, error) {
	q := u.Query()
	q.Set(cast.ToString(k), cast.ToString(v))
	u.RawQuery = q.Encode()
	return u, nil
}

func urlQuery(v interface{}, u *url.URL) (*url.URL, error) {
	vals, err := formEncoder.Encode(v)
	if err != nil {
		return u, err
	}

	q := u.Query()
	for k, vv := range vals {
		for i, v := range vv {
			if i == 0 {
				q.Set(k, v)
			} else {
				q.Add(k, v)
			}
		}
	}
	u.RawQuery = q.Encode()

	return u, nil
}
