package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/backends"
)

func init() {
	rootCmd.AddCommand(updateEmbargoes)
}

var updateEmbargoes = &cobra.Command{
	Use:   "update-embargoes",
	Short: "Update embargoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()
		if err := updatePublicationEmbargoes(services); err != nil {
			return err
		}
		if err := updateDatasetEmbargoes(services); err != nil {
			return err
		}
		return nil
	},
}

func updatePublicationEmbargoes(services *backends.Services) error {
	n, err := services.Repo.UpdatePublicationEmbargoes()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("updated %d publication embargoes", n))

	return nil
}

func updateDatasetEmbargoes(e *backends.Services) error {
	n, err := e.Repo.UpdateDatasetEmbargoes()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("updated %d dataset embargoes", n))

	return nil
}
