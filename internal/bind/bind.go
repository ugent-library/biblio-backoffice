package bind

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

func PathValues(r *http.Request) url.Values {
	if PathValuesFunc == nil {
		return nil
	}
	return PathValuesFunc(r)
}

func RequestPath(r *http.Request, v interface{}) error {
	return Path(PathValues(r), v)
}

func Path(vals url.Values, v interface{}) error {
	return pathDecoder.Decode(v, vals)
}

func RequestQuery(r *http.Request, v interface{}) error {
	return Query(r.URL.Query(), v)
}

func Query(vals url.Values, v interface{}) error {
	return queryDecoder.Decode(v, vals)
}

func RequestForm(r *http.Request, v interface{}) error {
	r.ParseForm()
	return Form(r.Form, v)
}

func Form(vals url.Values, v interface{}) error {
	return formDecoder.Decode(v, vals)
}

func Request(r *http.Request, v interface{}) error {
	if err := RequestPath(r, v); err != nil {
		return err
	}

	if err := Query(r.URL.Query(), v); err != nil {
		return err
	}

	r.ParseForm()

	return Form(r.Form, v)
}
