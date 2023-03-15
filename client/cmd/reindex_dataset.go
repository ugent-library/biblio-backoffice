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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.ReindexDatasetsRequest{}
	stream, err := c.ReindexDatasets(context.Background(), req)
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
}
