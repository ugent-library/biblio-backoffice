package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type PublicationAuthors struct {
	Context
}

func NewPublicationAuthors(c Context) *PublicationAuthors {
	return &PublicationAuthors{c}
}

func (c *PublicationAuthors) PromoteSearchAuthor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	person := &struct {
		FirstName string
		LastName  string
	}{}

	if err := forms.Decode(person, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := person.FirstName + " " + person.LastName
	people, _ := c.Engine.SuggestPersons(q)

	length := strconv.Itoa(len(people))

	w.Header().Set("HX-Trigger", "ITPromoteModal")
	w.Header().Set("HX-Trigger-After-Swap", "ITPromoteModalAfterSwap")
	w.Header().Set("HX-Trigger-After-Settle", "ITPromoteModalAfterSettle")

	c.Render.HTML(w, 200,
		"publication/authors/_modal_promote_author",
		views.NewData(c.Render, r, struct {
			ID       string
			People   []models.Person
			Length   string
			RowDelta string
		}{
			id,
			people,
			length,
			muxRowDelta,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
