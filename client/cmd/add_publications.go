package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(AddPublicationsCmd)
}

var AddPublicationsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return AddPublications(cmd, args)
	},
}

func AddPublications(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)

	if err != nil {
		return fmt.Errorf("could not connect with the server: %w", err)
	}

	defer cnx.Close()

	stream, err := c.AddPublications(context.Background())
	if err != nil {
		return fmt.Errorf("could not create a grpc stream: %w", err)
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

	reader := bufio.NewReader(cmd.InOrStdin())
	lineNo := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read line from input: %w", err)
		}

		lineNo++

		p := &api.Publication{
			Payload: line,
		}

		req := &api.AddPublicationsRequest{Publication: p}
		if err := stream.Send(req); err != nil {
			return fmt.Errorf("could not send publication to the server: %w", err)
		}
	}

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
