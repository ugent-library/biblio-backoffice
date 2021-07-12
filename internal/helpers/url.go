package helpers

import (
	"html/template"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/spf13/cast"
)

func URL(r *mux.Router) template.FuncMap {
	return template.FuncMap{
		"urlFor":     urlFor(r),
		"urlPathFor": urlPathFor(r),
		"urlSet":     urlSet,
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
