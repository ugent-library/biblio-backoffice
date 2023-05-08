package cmd

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(CleanupDatasetsCmd)
}

var CleanupDatasetsCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up inconsistencies across store datasets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CleanupDatasets(cmd, args)
	},
}

func CleanupDatasets(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.CleanupDatasetsRequest{}
		stream, err := c.CleanupDatasets(context.Background(), req)
		if err != nil {
			return err
		}

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

				return err
			}

			// Application level error
			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.Printf("%s\n", sre.Message())
			}

			if rr := res.GetMessage(); rr != "" {
				cmd.Printf("%s\n", rr)
			}
		}

		return nil
	})
}
