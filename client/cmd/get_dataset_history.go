package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/server"
)

type GetDatasetHistoryCmd struct {
	RootCmd
}

func (c *GetDatasetHistoryCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-history [id]",
		Short: "Get dataset history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetDatasetHistoryCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.GetDatasetHistoryRequest{Id: args[0]}
	stream, err := c.Client.GetDatasetHistory(context.Background(), req)
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

		j, err := json.Marshal(server.MessageToDataset(res.Dataset))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	}
}
