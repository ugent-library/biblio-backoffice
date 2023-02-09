package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

var ValidatePublicationsCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate publications",
	Run: func(cmd *cobra.Command, args []string) {
		ValidatePublications(cmd, args)
	},
}

func ValidatePublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	stream, err := c.ValidatePublications(context.Background())
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

			j, err := marshaller.Marshal(res)
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
