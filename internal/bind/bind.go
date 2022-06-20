package bind

import (
	"net/http"
	"net/url"

	"github.com/monoculum/formam/v3"
)

var (
	PathValuesFunc func(r *http.Request) url.Values
	// TODO the go-playground/form decoder doesn't seem to handle structs with multiple tags
	// for a single slice field (like in SearchArgs) well so use formam for now
	pathDecoder  = formam.NewDecoder(&formam.DecoderOptions{TagName: "path", IgnoreUnknownKeys: true})
	formDecoder  = formam.NewDecoder(&formam.DecoderOptions{TagName: "form", IgnoreUnknownKeys: true})
	queryDecoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "query", IgnoreUnknownKeys: true})
)

func PathValues(r *http.Request) url.Values {
	if PathValuesFunc != nil {
		return PathValuesFunc(r)
	}
	return nil
}

func RequestPath(r *http.Request, v interface{}) error {
	return Path(PathValues(r), v)
}

func Path(vals url.Values, v interface{}) error {
	return pathDecoder.Decode(vals, v)
}

func RequestForm(r *http.Request, v interface{}) error {
	r.ParseForm()
	return Form(r.Form, v)
}

func Form(vals url.Values, v interface{}) error {
	return formDecoder.Decode(vals, v)
}

func RequestQuery(r *http.Request, v interface{}) error {
	return Query(r.URL.Query(), v)
}

func Query(vals url.Values, v interface{}) error {
	return queryDecoder.Decode(vals, v)
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
