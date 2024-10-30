package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
)

func init() {
	createHandles.Flags().Bool("force", false, "force a recreation of all handle records")
	rootCmd.AddCommand(createHandles)
}

var createHandles = &cobra.Command{
	Use:   "create-handles",
	Short: "Create handles",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()

		if services.HandleService == nil {
			return errors.New("handle server updates are not enabled")
		}

		if force, _ := cmd.Flags().GetBool("force"); !force {
			recreatePublicationHandles(services)
			recreateDatasetHandles(services)
			return nil
		}

		createPublicationHandles(services)
		createDatasetHandles(services)

		return nil
	},
}

func createPublicationHandles(services *backends.Services) {
	repo := services.Repo

	var n int
	var err error

	repo.EachPublicationWithoutHandle(func(p *models.Publication) bool {
		h, e := services.HandleService.UpsertHandle(p.ID)
		if e != nil {
			err = fmt.Errorf("error adding handle for publication %s: %w", p.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for publication %s: %s", p.ID, h.Message)
			return false
		}

		p.Handle = h.GetFullHandleURL()
		if err = repo.UpdatePublication(p.SnapshotID, p, nil); err != nil {
			return false
		}

		n++

		return true
	})

	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("created %d publication handles", n))
}

func createDatasetHandles(services *backends.Services) {
	repo := services.Repo

	var n int
	var err error

	repo.EachDatasetWithoutHandle(func(d *models.Dataset) bool {
		h, e := services.HandleService.UpsertHandle(d.ID)
		if e != nil {
			err = fmt.Errorf("error adding handle for dataset %s: %w", d.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for dataset %s: %s", d.ID, h.Message)
			return false
		}

		d.Handle = h.GetFullHandleURL()
		if err = repo.UpdateDataset(d.SnapshotID, d, nil); err != nil {
			return false
		}

		n++

		return true
	})

	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("created %d dataset handles", n))
}

func recreatePublicationHandles(services *backends.Services) {
	repo := services.Repo

	var n int
	var err error

	repo.EachPublicationWithHandle(func(p *models.Publication) bool {
		h, e := services.HandleService.UpsertHandle(p.ID)
		if e != nil {
			err = fmt.Errorf("error adding handle for publication %s: %w", p.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for publication %s: %s", p.ID, h.Message)
			return false
		}
		n++
		return true
	})

	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("created %d publication handles", n))
}

func recreateDatasetHandles(services *backends.Services) {
	repo := services.Repo

	var n int
	var err error

	repo.EachDatasetWithHandle(func(d *models.Dataset) bool {
		h, e := services.HandleService.UpsertHandle(d.ID)
		if e != nil {
			err = fmt.Errorf("error adding handle for dataset %s: %w", d.ID, e)
			return false
		} else if !h.IsSuccess() {
			err = fmt.Errorf("error adding handle for dataset %s: %s", d.ID, h.Message)
			return false
		}
		n++
		return true
	})

	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("recreated %d dataset handles", n))
}
