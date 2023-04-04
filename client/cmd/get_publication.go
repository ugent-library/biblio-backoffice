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
	PublicationCmd.AddCommand(GetPublicationCmd)
}

var GetPublicationCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Long: `
	Retrieve the a single publication as a JSONL formatted record.
	The record will be outputted to stdout.

		$ ./biblio-backoffice publication get [ID] > publication.jsonl
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetPublication(cmd, args)
	},
}

func GetPublication(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		id := args[0]
		req := &api.GetPublicationRequest{Id: id}
		res, err := c.GetPublication(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		} else {
			cmd.Printf("%s", res.GetPublication().GetPayload())
		}

		return nil
	})
}
