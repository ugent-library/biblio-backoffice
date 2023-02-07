package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type AddDatasetsCmd struct {
	RootCmd
}

func (c *AddDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add datasets",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *AddDatasetsCmd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.AddDatasets(context.Background())
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

		d := &api.Dataset{
			Payload: line,
		}

		req := &api.AddDatasetsRequest{Dataset: d}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
