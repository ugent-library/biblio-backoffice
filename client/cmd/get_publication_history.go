package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type GetPublicationHistoryCmd struct {
	RootCmd
}

func (c *GetPublicationHistoryCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-history [id]",
		Short: "Get publication history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetPublicationHistoryCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.GetPublicationHistoryRequest{Id: args[0]}
	stream, err := c.Client.GetPublicationHistory(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		fmt.Printf("%s\n", res.Publication.Payload)
	}
}
