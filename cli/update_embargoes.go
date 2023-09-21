package cli

import (
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(updateEmbargoes)
}

var updateEmbargoes = &cobra.Command{
	Use:   "update-embargoes",
	Short: "Update embargoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		e := Services()
		logger := newLogger()

		if err := updatePublicationEmbargoes(e, logger); err != nil {
			return err
		}
		updateDatasetEmbargoes(e, logger)

		return nil
	},
}

func updatePublicationEmbargoes(e *backends.Services, logger *zap.SugaredLogger) error {
	n, err := e.Repo.UpdatePublicationEmbargoes()
	if err != nil {
		return err
	}

	logger.Infof("updated %d publication embargoes", n)

	return nil
}

func updateDatasetEmbargoes(e *backends.Services, logger *zap.SugaredLogger) error {
	n, err := e.Repo.UpdateDatasetEmbargoes()
	if err != nil {
		return err
	}

	logger.Infof("updated %d dataset embargoes", n)

	return nil
}
