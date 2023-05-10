package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(createHandles)
}

var createHandles = &cobra.Command{
	Use:   "create-handles",
	Short: "Create handles",
	Run: func(cmd *cobra.Command, args []string) {
		services := Services()
		logger := newLogger()

		if services.HandleService == nil {
			logger.Fatal("handle server updates are not enabled")
		}

		createPublicationHandles(services, logger)
		createDatasetHandles(services, logger)
	},
}

func createPublicationHandles(services *backends.Services, logger *zap.SugaredLogger) {
	repo := services.Repository

	repo.AddPublicationListener(func(p *models.Publication) {
		if p.DateUntil == nil {
			if err := services.PublicationSearchService.Index(p); err != nil {
				logger.Fatalf("error indexing publication %s: %v", p.ID, err)
			}
		}
	})

	var n int
	var err error

	repo.EachPublicationWithoutHandle(func(p *models.Publication) bool {
		h, e := services.HandleService.UpsertHandle(p.ID)
		if err != nil {
			err = fmt.Errorf("error adding handle for publication %s: %w", p.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for publication %s: %s", p.ID, h.Message)
			return false
		}

		logger.Infof("added handle url %s to publication %s", h.GetFullHandleURL(), p.ID)

		p.Handle = h.GetFullHandleURL()
		if err = repo.UpdatePublication(p.SnapshotID, p, nil); err != nil {
			return false
		}

		n++

		return true
	})

	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("created %d publication handles", n)
}

func createDatasetHandles(services *backends.Services, logger *zap.SugaredLogger) {
	repo := services.Repository

	repo.AddDatasetListener(func(p *models.Dataset) {
		if p.DateUntil == nil {
			if err := services.DatasetSearchService.Index(p); err != nil {
				logger.Fatalf("error indexing dataset %s: %v", p.ID, err)
			}
		}
	})

	var n int
	var err error

	repo.EachDatasetWithoutHandle(func(d *models.Dataset) bool {
		h, e := services.HandleService.UpsertHandle(d.ID)
		if err != nil {
			err = fmt.Errorf("error adding handle for dataset %s: %w", d.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for dataset %s: %s", d.ID, h.Message)
			return false
		}

		logger.Infof("added handle url %s to dataset %s", h.GetFullHandleURL(), d.ID)

		d.Handle = h.GetFullHandleURL()
		if err = repo.UpdateDataset(d.SnapshotID, d, nil); err != nil {
			return false
		}

		n++

		return true
	})

	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("created %d dataset handles", n)
}
