package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/server"
)

type ImportPublicationsCmd struct {
	RootCmd
}

func (c *ImportPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *ImportPublicationsCmd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.ImportPublications(context.Background())
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

		req := &api.ImportPublicationsRequest{Publication: server.PublicationToMessage(p)}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
