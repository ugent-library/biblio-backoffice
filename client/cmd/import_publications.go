package cmd

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(ImportPublicationsCmd)
}

var ImportPublicationsCmd = &cobra.Command{
	Use:   "import",
	Short: "Import publications",
	Run: func(cmd *cobra.Command, args []string) {
		ImportPublications(cmd, args)
	},
}

func ImportPublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	stream, err := c.ImportPublications(context.Background())
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

		p := &api.Publication{
			Payload: line,
		}

		req := &api.ImportPublicationsRequest{Publication: p}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
