package render

import (
	"net/http"

	"github.com/gorilla/schema"
)

var formDecoder = schema.NewDecoder()

func init() {
	formDecoder.SetAliasTag("form")
	formDecoder.IgnoreUnknownKeys(true)
}

func MustBindForm(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	r.ParseForm()

	if err := formDecoder.Decode(v, r.Form); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
}
