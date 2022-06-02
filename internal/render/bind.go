package render

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

var formDecoder = schema.NewDecoder()
var pathDecoder = schema.NewDecoder()

func init() {
	formDecoder.SetAliasTag("form")
	formDecoder.IgnoreUnknownKeys(true)
	pathDecoder.SetAliasTag("path")
	pathDecoder.IgnoreUnknownKeys(true)
}

func MustBindForm(w http.ResponseWriter, vals url.Values, v interface{}) bool {
	if err := formDecoder.Decode(v, vals); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}

func MustBindPath(w http.ResponseWriter, vals url.Values, v interface{}) bool {
	if err := pathDecoder.Decode(v, vals); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}
