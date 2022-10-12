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

type DatasetHistoryCmd struct {
	RootCmd
}

func (c *DatasetHistoryCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [id]",
		Short: "Dataset history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *DatasetHistoryCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.DatasetHistoryRequest{Id: args[0]}
	stream, err := c.Client.DatasetHistory(context.Background(), req)
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
