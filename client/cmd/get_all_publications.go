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

type GetAllPublicationsCmd struct {
	RootCmd
}

func (c *GetAllPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get all publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetAllPublicationsCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.GetAllPublicationsRequest{}
	stream, err := c.Client.GetAllPublications(context.Background(), req)
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
