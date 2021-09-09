package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/DAO"
	"github.com/ugent-library/biblio-backend/internal/ctx"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/presenters"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type Publication struct {
	engine *engine.Engine
	render *render.Render
}

type PublicationListVars struct {
	SearchArgs *engine.SearchArgs
	Hits       *engine.PublicationHits
}

type PublicationShowVars struct {
	Pub *engine.Publication
}

type PublicationNewVars struct {
}

func NewPublication(e *engine.Engine, r *render.Render) *Publication {
	return &Publication{engine: e, render: r}
}

func (c *Publication) List(w http.ResponseWriter, r *http.Request) {
	args := engine.NewSearchArgs()
	if err := forms.Decode(args, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hits, err := c.engine.UserPublications(ctx.GetUser(r).ID, args)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/list", PublicationListVars{SearchArgs: args, Hits: hits})
}

func (c *Publication) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, http.StatusOK, "publication/show", PublicationShowVars{Pub: pub})
}

func (c *Publication) Description(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := pub.Data()
	fmt.Println(data)
	description := &DAO.DescriptionDAO{
		PublicationType:     data["type"].(string),
		DOI:                 data["handle"].(string),
		ISXN:                "ISSN:XXX-XXX-XXX",
		Title:               data["title"].(string),
		AlternativeTitle:    "An alternative journal article title",
		Classification:      "A1 publication",
		Conference:          "Summit on journal articles 2021",
		ConferenceLocation:  "Ghent",
		ConferenceOrganiser: "University Library Ghent",
		ConferenceDate:      "22.05.2021-25.05.2021",
	}

	// Feed it into a presenter
	var presenter = &presenters.DescriptionPresenter{Description: description}

	// Create a View
	var view = &views.DescriptionView{Presenter: presenter}

	content := view.Render()

	c.render.HTML(w, http.StatusOK, "publication/description", struct {
		Pub     *engine.Publication
		Content string
	}{Content: content, Pub: pub})
}

func (c *Publication) New(w http.ResponseWriter, r *http.Request) {
	c.render.HTML(w, http.StatusOK, "publication/new", PublicationNewVars{})
}
