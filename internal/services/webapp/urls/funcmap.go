package urls

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
)

var formEncoder = form.NewEncoder()

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
			return route.URL(vars...)
		}
		return nil, fmt.Errorf("route %s not found", name)
	}
}

func pathFor(r *mux.Router) func(string, ...string) (*url.URL, error) {
	return func(name string, vars ...string) (*url.URL, error) {
		if route := r.Get(name); route != nil {
			return route.URLPath(vars...)
		}
		return nil, fmt.Errorf("route %s not found", name)
	}
}

func query(v interface{}, u *url.URL) (*url.URL, error) {
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

func querySet(k, v string, u *url.URL) (*url.URL, error) {
	q := u.Query()
	q.Set(k, v)
	u.RawQuery = q.Encode()
	return u, nil
}
