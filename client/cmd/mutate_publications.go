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
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	PublicationCmd.AddCommand(MutatePublicationsCmd)
}

var MutatePublicationsCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate publications",
	Long: `
	Add one or more mutations from a JSONL (JSON Lines) formatted file via stdin.
	Each line represents a single mutation.

	Outputs either a success message with the publication ID or an error message.
	Each message contains the number pointing to the corresponding line in the input file:

		$ ./biblio-backoffice publication mutate < file.jsonl
		mutated publication [ID] at line [LINENO]
		# or
		failed to mutate publication [ID] at line [LINENO]: [MSG]

	Example input file:

	    {"id":"1234", "op":"add_keyword", "args":["dna", "double helix"]}
	    {"id":"2345", "op":"set_classification", "args":["A3"]}

	`,
	RunE: MutatePublications,
}

func MutatePublications(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		stream, err := c.MutatePublications(context.Background())
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

			req := &api.MutateRequest{}
			if err := protojson.Unmarshal(line, req); err != nil {
				return fmt.Errorf("could not decode mutation: %w", err)
			}

			if err := stream.Send(req); err != nil {
				return fmt.Errorf("could not send mutation to the server: %w", err)
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
