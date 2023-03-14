package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(PurgeDatasetCmd)
}

var PurgeDatasetCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge dataset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return PurgeDataset(cmd, args)
	},
}

func PurgeDataset(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("ContextDeadlineExceeded: true")
	}

	id := args[0]
	req := &api.PurgeDatasetRequest{Id: id}
	res, err := c.PurgeDataset(context.Background(), req)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errors.New(st.Message())
		}

		return err
	}

	if ge := res.GetError(); ge != nil {
		sre := status.FromProto(ge)
		cmd.Printf("%s", sre.Message())
	}

	if res.GetOk() {
		cmd.Printf("purged dataset %s", id)
	}

	return nil
}
