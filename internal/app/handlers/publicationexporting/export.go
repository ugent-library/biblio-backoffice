package publicationexporting

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
)

func (h *Handler) ExportByCurationSearch(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanCurate() {
		render.Forbidden(w, r)
		return
	}

	exporterFactory, exporterFactoryFound := h.PublicationListExporters[ctx.ExportArgs.Format]
	if !exporterFactoryFound {
		e := fmt.Errorf("unable to find exporter %s", ctx.ExportArgs.Format)
		h.Logger.Errorw(
			"publication search: could not find exporter",
			"errors", e,
			"user", ctx.User.ID,
		)
		render.InternalServerError(w, r, e)
		return
	}
	exporter := exporterFactory(w)

	searcher := h.PublicationSearchIndex.WithScope("status", "private", "public", "returned")
	searcherErr := searcher.Each(ctx.SearchArgs, 10000, func(pub *models.Publication) {
		exporter.Add(pub)
	})

	if searcherErr != nil {
		h.Logger.Errorw(
			"publication search: unable to execute search",
			"errors", searcherErr,
			"user", ctx.User.ID,
		)
		render.InternalServerError(w, r, fmt.Errorf("unable to execute search"))
		return
	}

	/*
		TODO:
			- move to /tasks. Then we need to store this temporary file somewhere
			For now keep maxSize
			- use stream writer. Now the exporters assume a temporary file,
			and then write everything at once to the http writer. Possible caveats
			using a stream writer is that the headers need to proceed the content,
			during which an error may occur. In which case you would get a non working
			file, with the right file name, containing the server error
	*/
	fileName := fmt.Sprintf("publications_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
	contentDisposition := fmt.Sprintf("attachment;filename=%s", fileName)
	w.Header().Set("Content-Type", exporter.GetContentType())
	w.Header().Set("Content-Disposition", contentDisposition)

	if err := exporter.Flush(); err != nil {
		h.Logger.Errorw("publication search: could not export search", "errors", err, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
}
