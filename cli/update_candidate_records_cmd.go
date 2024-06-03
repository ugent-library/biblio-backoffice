package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/recordsources"
	_ "github.com/ugent-library/biblio-backoffice/recordsources/plato"
)

func init() {
	rootCmd.AddCommand(updateCandidateRecords)
}

var updateCandidateRecords = &cobra.Command{
	Use:   "update-candidate-records",
	Short: "Update candidate records",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()

		src, err := recordsources.New("plato")
		if err != nil {
			return err
		}

		err = src.GetRecords(context.Background(), func(srcRec recordsources.Record) error {
			oldCandidateRec, err := services.Repo.GetCandidateRecordBySource(context.TODO(), srcRec.SourceName(), srcRec.SourceID())
			if err != nil {
				if !errors.Is(err, models.ErrNotFound) {
					return err
				}
			}

			if oldCandidateRec != nil {
				logger.Warn(fmt.Sprintf("skipping duplicate candidate record from source %s/%s: already found in %s", srcRec.SourceName(), srcRec.SourceID(), oldCandidateRec.ID))
				return nil
			}

			candidateRec, err := srcRec.ToCandidateRecord(services)
			if err != nil {
				return err
			}

			if err := services.Repo.AddCandidateRecord(context.TODO(), candidateRec); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	},
}
