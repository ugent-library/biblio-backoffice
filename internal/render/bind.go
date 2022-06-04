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

func PathValues(r *http.Request) url.Values {
	if PathValuesFunc == nil {
		return nil
	}
	return PathValuesFunc(r)
}

func MustBindPath(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	return MustBindPathValues(w, PathValues(r), v)
}

func MustBindPathValues(w http.ResponseWriter, vals url.Values, v interface{}) bool {
	if err := pathDecoder.Decode(v, vals); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}

func MustBindQuery(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	return MustBindQueryValues(w, r.URL.Query(), v)
}

func MustBindQueryValues(w http.ResponseWriter, vals url.Values, v interface{}) bool {
	if err := queryDecoder.Decode(v, vals); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}

func MustBindForm(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	r.ParseForm()
	return MustBindFormValues(w, r.Form, v)
}

func MustBindFormValues(w http.ResponseWriter, vals url.Values, v interface{}) bool {
	if err := formDecoder.Decode(v, vals); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}

func MustBind(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if !MustBindPath(w, r, v) {
		return false
	}

	if !MustBindQueryValues(w, r.URL.Query(), v) {
		return false
	}

	r.ParseForm()

	return MustBindFormValues(w, r.Form, v)
}
