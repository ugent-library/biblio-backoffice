package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/server"
)

type AddPublicationsCmd struct {
	RootCmd
}

func (c *AddPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *AddPublicationsCmd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.AddPublications(context.Background())
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
			log.Println(res.Message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	lineNo := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		lineNo++

		p := &models.Publication{}
		if err := json.Unmarshal(line, p); err != nil {
			log.Fatalf("Unable to decode publication at line %d : %v", lineNo, err)
		}

		req := &api.AddPublicationsRequest{Publication: server.PublicationToMessage(p)}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
