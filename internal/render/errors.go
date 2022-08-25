package render

import (
	"net/http"
)

var AuthURL string

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, AuthURL, http.StatusTemporaryRedirect)
	return
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	return
}
