package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type GetPublicationCmd struct {
	RootCmd
}

func (c *GetPublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get publication by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.SetOutput(cmd.OutOrStdout())

			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetPublicationCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id := args[0]
	req := &api.GetPublicationRequest{Id: id}
	res, err := c.Client.GetPublication(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Printf("%s\n", res.Publication.Payload)
}
