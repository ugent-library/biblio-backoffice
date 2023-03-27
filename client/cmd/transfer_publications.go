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
	err := client.Transmit(config, func(c api.BiblioClient) error {
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

		waitc := make(chan struct{})
		errorc := make(chan error)

		go func() {
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}

				// return gRPC level error
				if err != nil {
					errorc <- err
					return
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
		}()

		stream.CloseSend()

		select {
		case errc := <-errorc:
			if st, ok := status.FromError(errc); ok {
				return errors.New(st.Message())
			}
		case <-waitc:
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
