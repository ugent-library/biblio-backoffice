package cmd

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(GetPublicationHistoryCmd)
}

var GetPublicationHistoryCmd = &cobra.Command{
	Use:   "get-history [id]",
	Short: "Get publication history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetPublicationHistory(cmd, args)
	},
}

func GetPublicationHistory(cmd *cobra.Command, args []string) error {
	err := client.Transmit(config, func(c api.BiblioClient) error {
		req := &api.GetPublicationHistoryRequest{Id: args[0]}
		stream, err := c.GetPublicationHistory(context.Background(), req)
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

			if rr := res.GetPublication(); rr != nil {
				cmd.Printf("%s\n", rr.GetPayload())
			}
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
