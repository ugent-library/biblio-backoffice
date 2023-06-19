package commands

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
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		logger := newLogger()

		updatePublicationEmbargoes(e, logger)
		updateDatasetEmbargoes(e, logger)
	},
}

func updatePublicationEmbargoes(e *backends.Services, logger *zap.SugaredLogger) {
	n, err := e.Repository.UpdatePublicationEmbargoes()

	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("updated %d publication embargoes", n)
}

func updateDatasetEmbargoes(e *backends.Services, logger *zap.SugaredLogger) {
	n, err := e.Repository.UpdateDatasetEmbargoes()

	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("updated %d dataset embargoes", n)
}
