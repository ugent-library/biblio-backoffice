package cmd

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(ValidatePublicationsCmd)
}

var ValidatePublicationsCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ValidatePublications(cmd, args)
	},
}

func ValidatePublications(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	stream, err := c.ValidatePublications(context.Background())
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
			}

			// Application level error
			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.Printf("%s\n", sre.Message())
			}

			if rr := res.GetResults(); rr != nil {
				j, err := marshaller.Marshal(res)
				if err != nil {
					errorc <- err
				}
				cmd.Printf("%s\n", j)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		p := &api.Publication{
			Payload: line,
		}

		req := &api.ValidatePublicationsRequest{Publication: p}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
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
