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
	DatasetCmd.AddCommand(ImportDatasetsCmd)
}

var ImportDatasetsCmd = &cobra.Command{
	Use:   "import",
	Short: "Import datasets",
	Run: func(cmd *cobra.Command, args []string) {
		ImportDatasets(cmd, args)
	},
}

func ImportDatasets(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	stream, err := c.ImportDatasets(context.Background())
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

		req := &api.ImportDatasetsRequest{Dataset: d}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
}
