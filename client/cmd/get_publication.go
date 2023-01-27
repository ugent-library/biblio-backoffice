package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/server"
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

	j, err := json.Marshal(server.MessageToPublication(res.Publication))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
