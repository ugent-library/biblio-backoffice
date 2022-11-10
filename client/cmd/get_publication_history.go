package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/server"
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

		j, err := json.Marshal(server.MessageToPublication(res.Publication))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	}
}