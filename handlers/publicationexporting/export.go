package publicationexporting

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

func ExportByCurationSearch(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	format := bind.PathValue(r, "format")
	if format == "" {
		format = "xlsx"
	}

	exporterFactory, exporterFactoryFound := c.PublicationListExporters[format]
	if !exporterFactoryFound {
		e := fmt.Errorf("unable to find exporter %s", format)
		c.Log.Errorw(
			"publication search: could not find exporter",
			"errors", e,
			"user", c.User.ID,
		)
		c.HandleError(w, r, httperror.NotFound)
		return
	}
	exporter := exporterFactory(w)

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.Log.Warnw("publication search: could not bind search arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, err)
		return
	}
	searcher := c.PublicationSearchIndex.WithScope("status", "private", "public", "returned")
	searcherErr := searcher.Each(searchArgs, 10000, func(pub *models.Publication) {
		exporter.Add(pub)
	})

	if searcherErr != nil {
		c.Log.Errorw(
			"publication search: unable to execute search",
			"errors", searcherErr,
			"user", c.User.ID,
		)
		c.HandleError(w, r, fmt.Errorf("unable to execute search: %w", searcherErr))
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
		c.Log.Errorw("publication search: could not export search", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, err)
		return
	}
}
