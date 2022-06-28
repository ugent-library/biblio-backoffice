package urls

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
)

var queryEncoder = form.NewEncoder()

func init() {
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}

func FuncMap(r *mux.Router) template.FuncMap {
	return template.FuncMap{
		"urlFor":   urlFor(r),
		"pathFor":  pathFor(r),
		"query":    query,
		"querySet": querySet,
	}
}

func urlFor(r *mux.Router) func(string, ...string) (*url.URL, error) {
	return func(name string, vars ...string) (*url.URL, error) {
		if route := r.Get(name); route != nil {
			u, err := route.URL(vars...)
			if err != nil {
				return nil, fmt.Errorf("can't reverse route %s: %w", name, err)
			}
			return u, nil
		}
		return nil, fmt.Errorf("route %s not found", name)
	}
}

func pathFor(r *mux.Router) func(string, ...string) (*url.URL, error) {
	return func(name string, vars ...string) (*url.URL, error) {
		if route := r.Get(name); route != nil {
			u, err := route.URLPath(vars...)
			if err != nil {
				return nil, fmt.Errorf("can't reverse route %s: %w", name, err)
			}
			return u, nil
		}
		return nil, fmt.Errorf("route %s not found", name)
	}
}

func query(v interface{}, u *url.URL) (*url.URL, error) {
	vals, err := queryEncoder.Encode(v)
	if err != nil {
		return u, err
	}

	u.RawQuery = vals.Encode()

	return u, nil
}

func querySet(k, v string, u *url.URL) (*url.URL, error) {
	q := u.Query()
	q.Set(k, v)
	u.RawQuery = q.Encode()

	return u, nil
}
