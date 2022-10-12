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

type PublicationHistoryCmd struct {
	RootCmd
}

func (c *PublicationHistoryCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [id]",
		Short: "Publication history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *PublicationHistoryCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.PublicationHistoryRequest{Id: args[0]}
	stream, err := c.Client.PublicationHistory(context.Background(), req)
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
