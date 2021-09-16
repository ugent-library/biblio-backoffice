package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/unrolled/render"
)

type PublicationsDetails struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationsDetails(e *engine.Engine, r *render.Render) *PublicationsDetails {
	return &PublicationsDetails{
		engine: e,
		render: r,
	}
}

func (c *PublicationsDetails) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	buf := &bytes.Buffer{}
	if tmpl := c.render.TemplateLookup(fmt.Sprintf("publication/details/_%s", pub.Type)); tmpl != nil {
		tmpl.Execute(buf, pub)
	}

	fmt.Fprintf(w, buf.String())
}

func (c *PublicationsDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	buf := &bytes.Buffer{}

	if tmpl := c.render.TemplateLookup(fmt.Sprintf("publication/details/_%s_edit_form", pub.Type)); tmpl != nil {
		tmpl.Execute(buf, pub)
	}

	fmt.Fprintf(w, buf.String())
}

func (c *PublicationsDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	buf := &bytes.Buffer{}

	pub.Title = "Mock title"

	if tmpl := c.render.TemplateLookup(fmt.Sprintf("publication/details/_%s_edit_submit", pub.Type)); tmpl != nil {
		tmpl.Execute(buf, pub)
	}

	fmt.Fprintf(w, buf.String())
}
