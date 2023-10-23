package cli

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/backends/ecoom"
	"github.com/ugent-library/biblio-backoffice/backends/jcr"
	"github.com/ugent-library/biblio-backoffice/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	rootCmd.AddCommand(updateExternalFields)
}

var updateExternalFields = &cobra.Command{
	Use:   "update-externalfields",
	Short: "Update publication external_fields",
	RunE: func(cmd *cobra.Command, args []string) error {
		mongoClient, err := mongo.Connect(
			context.Background(),
			options.Client().ApplyURI(config.MongoDBURL))

		if err != nil {
			return errors.Wrap(err, "unable to initialize connection to mongodb")
		}

		fixers := []func(context.Context, *models.Publication) error{
			ecoom.NewPublicationFixer(mongoClient),
			jcr.NewPublicationFixer(mongoClient),
		}
		repo := newServices().Repo

		var lastErr error
		ctx := context.TODO()

		repo.EachPublicationWithStatus("public", func(p *models.Publication) bool {
			for _, fixer := range fixers {
				if err := fixer(ctx, p); err != nil {
					lastErr = err
					return false
				}
			}

			if err := repo.UpdatePublication(p.SnapshotID, p, nil); err != nil {
				logger.Errorf("unable to update external fields in publication %s: %s", p.ID, err)
				lastErr = err
				return false
			}

			logger.Infof("successfully updated external fields in publication %s", p.ID)
			return true
		})

		return lastErr
	},
}
