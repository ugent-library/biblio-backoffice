package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type PurgePublicationCmd struct {
	RootCmd
}

func (c *PurgePublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge [id]",
		Short: "Purge publication",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *PurgePublicationCmd) Run(args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id := args[0]
	req := &api.PurgePublicationRequest{Id: id}
	if _, err := c.Client.PurgePublication(ctx, req); err != nil {
		log.Fatal(err)
	}
}
