package render

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

var (
	PathValuesFunc func(r *http.Request) url.Values

	pathDecoder  = schema.NewDecoder()
	formDecoder  = schema.NewDecoder()
	queryDecoder = schema.NewDecoder()
)

func init() {
	pathDecoder.SetAliasTag("path")
	pathDecoder.IgnoreUnknownKeys(true)
	formDecoder.SetAliasTag("form")
	formDecoder.IgnoreUnknownKeys(true)
	queryDecoder.SetAliasTag("query")
	queryDecoder.IgnoreUnknownKeys(true)
}

func GetPathValues(r *http.Request) url.Values {
	if PathValuesFunc == nil {
		return nil
	}
	return PathValuesFunc(r)
}

func BindPath(r *http.Request, v interface{}) error {
	return BindPathValues(GetPathValues(r), v)
}

func BindPathValues(vals url.Values, v interface{}) error {
	return pathDecoder.Decode(v, vals)
}

func BindQuery(r *http.Request, v interface{}) error {
	return BindQueryValues(r.URL.Query(), v)
}

func BindQueryValues(vals url.Values, v interface{}) error {
	return queryDecoder.Decode(v, vals)
}

func BindForm(r *http.Request, v interface{}) error {
	r.ParseForm()
	return BindFormValues(r.Form, v)
}

func BindFormValues(vals url.Values, v interface{}) error {
	return formDecoder.Decode(v, vals)
}

func Bind(r *http.Request, v interface{}) error {
	if err := BindPath(r, v); err != nil {
		return err
	}

	if err := BindQueryValues(r.URL.Query(), v); err != nil {
		return err
	}

	r.ParseForm()

	return BindFormValues(r.Form, v)
}
