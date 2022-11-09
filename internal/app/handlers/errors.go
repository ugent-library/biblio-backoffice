package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/render"
)

var AuthURL string

type Context struct {
	BaseContext
}

type YieldNotFound struct {
	Context
	PageTitle string
	ActiveNav string
}

// TODO Make these user friendly pages with a nice error message informing the user on
//    a. What went wrong
//    b. How to proceed

// HTTP 500 error
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// HTTP 404 error
func NotFound(w http.ResponseWriter, r *http.Request, ctx BaseContext, err error) {
	// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	render.NotFoundLayout(w, "layouts/default", "pages/notfound", YieldNotFound{
		Context: Context{
			BaseContext: ctx,
		},
		PageTitle: "Biblio",
	})
}

// HTTP 400 error
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// HTTP 401 error, redirects the user to the authentication url
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, AuthURL, http.StatusTemporaryRedirect)
}

// HTTP 403 error
func Forbidden(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
