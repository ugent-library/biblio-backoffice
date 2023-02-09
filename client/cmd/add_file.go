package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

// Set file buffer size
var fileBufSize = 524288

var AddFileCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		AddFile(cmd, args)
	},
}

func AddFile(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	stream, err := c.AddFile(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, fileBufSize)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &api.AddFileRequest{Chunk: buf[:n]}

		if err = stream.Send(req); err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.WriteString(res.Sha256 + "\n")
}
