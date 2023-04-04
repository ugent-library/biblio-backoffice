package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(UpdateDatasetCmd)
}

var UpdateDatasetCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dataset",
	Long: `
	Update one or multiple datasets.

	This command reads a JSONL formatted file from stdin and streams it to the store.

	It will output either a success message or an error message per record:

		$ ./biblio-backoffice dataset update < datasets.jsonl
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return UpdateDataset(cmd, args)
	},
}

func UpdateDataset(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		reader := bufio.NewReader(cmd.InOrStdin())
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("could not read from stdin: %v", err)
		}

		p := &api.Dataset{
			Payload: line,
		}

		req := &api.UpdateDatasetRequest{Dataset: p}
		res, err := c.UpdateDataset(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		}

		return nil
	})
}
