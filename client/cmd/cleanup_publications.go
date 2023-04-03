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
	PublicationCmd.AddCommand(CleanupPublicationsCmd)
}

var CleanupPublicationsCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Make publications consistent, clean up data anomalies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CleanupPublications(cmd, args)
	},
}

func CleanupPublications(cmd *cobra.Command, args []string) error {
	err := cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.CleanupPublicationsRequest{}
		stream, err := c.CleanupPublications(context.Background(), req)
		if err != nil {
			return err
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
