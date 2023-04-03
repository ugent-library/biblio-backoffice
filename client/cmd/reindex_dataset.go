package cmd

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(ReindexDatasetCmd)
}

var ReindexDatasetCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all datasets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ReindexDatasets(cmd, args)
	},
}

func ReindexDatasets(cmd *cobra.Command, args []string) error {
	err := cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.ReindexDatasetsRequest{}
		stream, err := c.ReindexDatasets(context.Background(), req)
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

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
