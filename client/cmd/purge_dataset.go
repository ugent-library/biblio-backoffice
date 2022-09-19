package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type PurgeDatasetCmd struct {
	RootCmd
}

func (c *PurgeDatasetCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge [id]",
		Short: "Purge dataset",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *PurgeDatasetCmd) Run(args []string) {
	id := args[0]
	req := &api.PurgeDatasetRequest{Id: id}
	if _, err := c.Client.PurgeDataset(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
