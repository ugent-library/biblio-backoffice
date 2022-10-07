package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
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
			log.Println(res.Messsage)
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

		pub := &api.Publication{}
		if err := c.Unmarshaller.Unmarshal(line, pub); err != nil {
			log.Fatal(err)
		}

		req := &api.AddPublicationsRequest{Publication: pub}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
