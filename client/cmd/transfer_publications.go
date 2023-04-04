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
	PublicationCmd.AddCommand(TransferPublicationsCmd)
}

var TransferPublicationsCmd = &cobra.Command{
	Use:   "transfer UID UID [PUBID]",
	Short: "Transfer publications between people",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return TransferPublications(cmd, args)
	},
}

func TransferPublications(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		source := args[0]
		dest := args[1]

		pubid := ""
		if len(args) > 2 {
			pubid = args[2]
		}

		req := &api.TransferPublicationsRequest{
			Src:           source,
			Dest:          dest,
			Publicationid: pubid,
		}

		stream, err := c.TransferPublications(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		stream.CloseSend()

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// read done.
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
