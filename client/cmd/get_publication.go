package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type GetPublicationCmd struct {
	RootCmd
}

func (c *GetPublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get publication by id",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

	j, err := c.Marshaller.Marshal(res.Publication)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
