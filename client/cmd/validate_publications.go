package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type ValidatePublicationsCmd struct {
	RootCmd
}

func (c *ValidatePublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *ValidatePublicationsCmd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.ValidatePublications(context.Background())
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

		p := &api.Publication{
			Payload: line,
		}

		req := &api.ValidatePublicationsRequest{Publication: p}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
