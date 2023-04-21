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
	PublicationCmd.AddCommand(PurgeAllPublicationsCmd)
}

var PurgeAllPublicationsCmd = &cobra.Command{
	Use:   "purge-all",
	Short: "Purge all publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return PurgeAllPublications(cmd, args)
	},
}

func init() {
	PurgeAllPublicationsCmd.Flags().BoolP("yes", "y", false, "are you sure?")
}

func PurgeAllPublications(cmd *cobra.Command, args []string) error {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		cmd.Printf("no confirmation flag set. you need to set the --yes flag")
		return nil
	}

	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.PurgeAllPublicationsRequest{
			Confirm: true,
		}
		res, err := c.PurgeAllPublications(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s\n", sre.Message())
		}

		if res.GetOk() {
			cmd.Printf("purged all publications")
		}

		return nil
	})
}
