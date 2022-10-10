package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/server"
)

type ValidateDatasetsCmd struct {
	RootCmd
}

func (c *ValidateDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate datasets",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *ValidateDatasetsCmd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.ValidateDatasets(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatal(err)
			}

			j, err := c.Marshaller.Marshal(res)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", j)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		d := &models.Dataset{}
		if err := json.Unmarshal(line, d); err != nil {
			log.Fatal(err)
		}

		req := &api.ValidateDatasetsRequest{Dataset: server.DatasetToMessage(d)}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
