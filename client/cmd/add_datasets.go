package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(AddDatasetsCmd)
}

var AddDatasetsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add datasets",
	Long: `
	Add one or more datasets from a JSONL (JSON Lines) formatted file via stdin.
	Each line represents a single dataset.

	Outputs either a success message with the dataset ID or an error message.
	Each message contains the number pointing to the corresponding line in the input file:

		$ ./biblio-backoffice dataset add < file.jsonl
		stored and indexed dataset [ID] at line [LINENO]
		failed to validate dataset [ID] at line [LINENO]: [MSG]
	`,
	RunE: AddDatasets,
}

func AddDatasets(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		stream, err := c.AddDatasets(context.Background())
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

			d := &api.Dataset{
				Payload: line,
			}

			req := &api.AddDatasetsRequest{Dataset: d}
			if err := stream.Send(req); err != nil {
				return fmt.Errorf("could not send dataset to the server: %w", err)
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
	})
}
