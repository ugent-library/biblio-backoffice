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
	DatasetCmd.AddCommand(GetDatasetCmd)
}

var GetDatasetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get dataset by id",
	Long: `
	Retrieve the a single dataset as a JSONL formatted record.
	The record will be outputted to stdout.

		$ ./biblio-backoffice dataset get [ID] > dataset.jsonl
	`,
	Args: cobra.ExactArgs(1),
	RunE: GetDataset,
}

func GetDataset(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		id := args[0]
		req := &api.GetDatasetRequest{Id: id}
		res, err := c.GetDataset(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		} else {
			cmd.Printf("%s", res.GetDataset().GetPayload())
		}

		return nil
	})
}
