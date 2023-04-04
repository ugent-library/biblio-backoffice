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
	PublicationCmd.AddCommand(PurgePublicationCmd)
}

var PurgePublicationCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge publication",
	Long: `
	Purge a single stored publication.

	Outputs either a success message with the dataset ID or an error message.

		$ ./biblio-backoffice publication purge [ID]
		purged publication [ID]
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return PurgePublication(cmd, args)
	},
}

func PurgePublication(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		id := args[0]
		req := &api.PurgePublicationRequest{Id: id}
		res, err := c.PurgePublication(context.Background(), req)

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
			cmd.Printf("purged publication %s", id)
		}

		return nil
	})
}
