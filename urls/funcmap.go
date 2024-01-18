package urls

import (
	"html/template"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/nics/ich"
)

var queryEncoder = form.NewEncoder()

func init() {
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}

// TODO split into mux and query packages
func FuncMap(r *ich.Mux, scheme, host string) template.FuncMap {
	return template.FuncMap{
		"urlFor":     urlFor(r, scheme, host),
		"pathFor":    pathFor(r),
		"query":      query,
		"querySet":   querySet,
		"queryAdd":   queryAdd,
		"queryDel":   queryDel,
		"queryClear": queryClear,
	}
}

func urlFor(r *ich.Mux, scheme, host string) func(string, ...string) *url.URL {
	return func(name string, pairs ...string) *url.URL {
		u := r.PathTo(name, pairs...)
		u.Host = host
		u.Scheme = scheme
		return u
	}
}

func pathFor(r *ich.Mux) func(string, ...string) *url.URL {
	return r.PathTo
}

func query(v any, u *url.URL) (*url.URL, error) {
	vals, err := queryEncoder.Encode(v)
	if err != nil {
		return u, err
	}

	newU := *u
	newU.RawQuery = vals.Encode()

	return &newU, nil
}

func querySet(k, v string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Set(k, v)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryAdd(k, v string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Add(k, v)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryDel(k string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Del(k)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryClear(u *url.URL) (*url.URL, error) {
	newU := *u
	newU.RawQuery = ""

	return &newU, nil
}
