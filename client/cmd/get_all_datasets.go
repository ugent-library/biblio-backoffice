package cmd

import (
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
	DatasetCmd.AddCommand(GetAllDatasetsCmd)
}

var GetAllDatasetsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all datasets",
	Long: `
	Retrieve all stored datasets as a stream of JSONL formatted records.
	The stream will be outputted to stdout.

		$ ./biblio-backoffice dataset get-all > datasets.jsonl
	`,
	RunE: GetAllDatasets,
}

func GetAllDatasets(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.GetAllDatasetsRequest{}
		stream, err := c.GetAllDatasets(context.Background(), req)
		if err != nil {
			return fmt.Errorf("error while reading stream: %v", err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				if st, ok := status.FromError(err); ok {
					return errors.New(st.Message())
				}

				return err
			}

			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.Printf("%s\n", sre.Message())
			}

			if rr := res.GetDataset(); rr != nil {
				cmd.Printf("%s\n", rr.GetPayload())
			}
		}

		return nil
	})
}
