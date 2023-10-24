package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ugent-library/biblio-backoffice/recordsources"
	_ "github.com/ugent-library/biblio-backoffice/recordsources/plato"
)

func init() {
	rootCmd.AddCommand(updateRecordCandidates)
}

var updateRecordCandidates = &cobra.Command{
	Use:   "update-record-candidates",
	Short: "Update record candidates",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range []string{"plato"} {
			src, err := recordsources.New(name, "")
			if err != nil {
				return err
			}
			recs, err := src.GetRecords(context.Background())
			if err != nil {
				return err
			}

			for _, rec := range recs {
				logger.Infof("rec metadata: %s", rec.Metadata)
			}
		}

		return nil
	},
}
