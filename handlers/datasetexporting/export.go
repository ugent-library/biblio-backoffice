package datasetexporting

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

	exporterFactory, exporterFactoryFound := c.DatasetListExporters[format]
	if !exporterFactoryFound {
		c.HandleError(w, r, httperror.NotFound)
		return
	}
	exporter := exporterFactory(w)

	searchArgs := models.NewSearchArgs()
	if err := bind.Request(r, searchArgs); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	searcher := c.DatasetSearchIndex.WithScope("status", "private", "public", "returned")
	searcherErr := searcher.Each(searchArgs, 10000, func(dataset *models.Dataset) {
		exporter.Add(dataset)
	})

	if searcherErr != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("unable to execute search: %w", searcherErr)))
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
	fileName := fmt.Sprintf("datasets_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
	contentDisposition := fmt.Sprintf("attachment;filename=%s", fileName)
	w.Header().Set("Content-Type", exporter.GetContentType())
	w.Header().Set("Content-Disposition", contentDisposition)

	if err := exporter.Flush(); err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not export search: %w", err)))
		return
	}

}
