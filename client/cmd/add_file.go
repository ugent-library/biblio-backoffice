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

type AddFileCMd struct {
	fileBufSize int
	RootCmd
}

func (c *AddFileCMd) Command() *cobra.Command {
	// Set file buffer size
	c.fileBufSize = 524288

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *AddFileCMd) Run(cmd *cobra.Command, args []string) {
	stream, err := c.Client.AddFile(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, c.fileBufSize)

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
