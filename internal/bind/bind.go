package bind

import (
	"net/http"
	"net/url"

	"github.com/go-playground/form/v4"
)

var (
	PathValuesFunc func(r *http.Request) url.Values

	pathDecoder  = form.NewDecoder()
	formDecoder  = form.NewDecoder()
	queryDecoder = form.NewDecoder()
	pathEncoder  = form.NewEncoder()
	formEncoder  = form.NewEncoder()
	queryEncoder = form.NewEncoder()
)

func init() {
	pathDecoder.SetTagName("path")
	pathDecoder.SetMode(form.ModeExplicit)
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	pathEncoder.SetTagName("path")
	pathEncoder.SetMode(form.ModeExplicit)
	formEncoder.SetTagName("form")
	formEncoder.SetMode(form.ModeExplicit)
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}

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

// include encoding helpers as a convenience

func EncodePath(v interface{}) (url.Values, error) {
	return pathEncoder.Encode(v)
}

func EncodeForm(v interface{}) (url.Values, error) {
	return formEncoder.Encode(v)
}

func EncodeQuery(v interface{}) (url.Values, error) {
	return queryEncoder.Encode(v)
}
