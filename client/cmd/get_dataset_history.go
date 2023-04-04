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
	DatasetCmd.AddCommand(GetDatasetHistoryCmd)
}

var GetDatasetHistoryCmd = &cobra.Command{
	Use:   "get-history [id]",
	Short: "Get dataset history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetDatasetHistory(cmd, args)
	},
}

func GetDatasetHistory(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.GetDatasetHistoryRequest{Id: args[0]}
		stream, err := c.GetDatasetHistory(context.Background(), req)
		if err != nil {
			return err
		}

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

			if rr := res.GetDataset(); rr != nil {
				cmd.Printf("%s\n", rr.GetPayload())
			}
		}

		return nil
	})
}
