package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(PurgeDatasetCmd)
}

var PurgeDatasetCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge dataset",
	Long: `
	Purge a single stored dataset.

	Outputs either a success message with the dataset ID or an error message.

		$ ./biblio-backoffice dataset purge [ID]
		purged dataset [ID]
	`,
	Args: cobra.ExactArgs(1),
	RunE: PurgeDataset,
}

func PurgeDataset(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		id := args[0]
		req := &api.PurgeDatasetRequest{Id: id}
		res, err := c.PurgeDataset(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}

			return err
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		}

		if res.GetOk() {
			cmd.Printf("purged dataset %s", id)
		}

		return nil
	})
}
