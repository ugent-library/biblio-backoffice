package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(ReindexPublicationCmd)
}

var ReindexPublicationCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all publications",
	Run: func(cmd *cobra.Command, args []string) {
		ReindexPublications(cmd, args)
	},
}

func ReindexPublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.ReindexPublicationsRequest{}
	stream, err := c.ReindexPublications(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt)

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("TEST TEST")
	// 		os.Exit(0)
	// 		return
	// 	case s := <-sigCh:
	// 		log.Printf("got signal %v, attempting graceful shutdown", s)
	// 		stream.CloseSend()
	// 		cancel()
	// 		cnx.Close()
	// 		os.Exit(0)
	// 	}
	// }()

	//go func() {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		fmt.Printf("%s\n", res.Message)
	}
	// }()
}
