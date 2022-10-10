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

type GetAllDatasetsCmd struct {
	RootCmd
}

func (c *GetAllDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get all datasets",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetAllDatasetsCmd) Run(cmd *cobra.Command, args []string) {

	req := &api.GetAllDatasetsRequest{}
	stream, err := c.Client.GetAllDatasets(context.Background(), req)
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
