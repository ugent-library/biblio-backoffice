package cmd

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(SyncPublicationContributorsCmd)
}

var SyncPublicationContributorsCmd = &cobra.Command{
	Use:   "sync-publication-contributors",
	Short: "Synchronize publication contributor attributes with person service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return SyncPublicationContributors(cmd, args)
	},
}

func SyncPublicationContributors(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}
	req := &api.SyncPublicationContributorsRequest{}
	stream, err := c.SyncPublicationContributors(context.Background(), req)
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
				cmd.PrintErrf("%s\n", sre.Message())
				continue
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
}
