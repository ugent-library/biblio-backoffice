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
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
)

func init() {
	DatasetCmd.AddCommand(ValidateDatasetsCmd)
}

var ValidateDatasetsCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate datasets",
	RunE:  ValidateDatasets,
}

func ValidateDatasets(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		stream, err := c.ValidateDatasets(context.Background())
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

			d := &api.Dataset{
				Payload: line,
			}

			req := &api.ValidateDatasetsRequest{Dataset: d}
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
			}
		}

		stream.CloseSend()
		<-waitc

		return nil
	})
}
