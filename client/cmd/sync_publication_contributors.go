package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	SyncPublicationContributorsCmd.Flags().BoolP("noop", "n", false, "inspect changes. Do not execute them.")
	PublicationCmd.AddCommand(SyncPublicationContributorsCmd)
}

var SyncPublicationContributorsCmd = &cobra.Command{
	Use:   "sync-publication-contributors",
	Short: "Synchronize publication contributor attributes with person service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return SyncPublicationContributors(cmd, args)
	},
}

func SyncPublicationContributors(cmd *cobra.Command, args []string) error {

	// --noop : dry mode
	noop, _ := cmd.Flags().GetBool("noop")

	err := cnx.Handle(config, func(c api.BiblioClient) error {

		req := &api.SyncPublicationContributorsRequest{Noop: noop}
		stream, err := c.SyncPublicationContributors(context.Background(), req)
		if err != nil {
			return err
		}

		// close send handle (only read handle needed)
		stream.CloseSend()

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			// return gRPC level error
			if err != nil {
				if st, ok := status.FromError(err); ok {
					return errors.New(st.Message())
				}
				break
			}

			// Application level error
			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.PrintErrf("%s\n", sre.Message())
				continue
			}

			if rr := res.GetContributorChange(); rr != nil {
				bytes, _ := json.Marshal(rr)
				cmd.Printf("%s\n", string(bytes))
			}
		}

		return nil
	})

	return err
}
