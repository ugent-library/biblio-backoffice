package views

import (
	"fmt"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
)

func paginationCount(c *ctx.Ctx, searchHits *models.SearchHits) string {
	if searchHits.TotalPages() > 1 {
		return fmt.Sprintf("Showing %d-%d of %d", searchHits.FirstOnPage(), searchHits.LastOnPage(), searchHits.Total)
	} else {
		return fmt.Sprintf("Showing %d", searchHits.Total)
	}
}
