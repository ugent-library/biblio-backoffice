package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(UpdateDatasetCmd)
}

var UpdateDatasetCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dataset",
	RunE: func(cmd *cobra.Command, args []string) error {
		return UpdateDataset(cmd, args)
	},
}

func UpdateDataset(cmd *cobra.Command, args []string) error {
	err := client.Transmit(config, func(c api.BiblioClient) error {
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

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("ContextDeadlineExceeded: true")
	}

	return err
}
