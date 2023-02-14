package cmd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	DatasetCmd.AddCommand(PurgeDatasetCmd)
}

var PurgeDatasetCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge dataset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		PurgeDataset(cmd, args)
	},
}

func PurgeDataset(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	id := args[0]
	req := &api.PurgeDatasetRequest{Id: id}
	if _, err := c.PurgeDataset(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
