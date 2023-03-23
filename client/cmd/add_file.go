package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	FileCmd.AddCommand(AddFileCmd)
}

// Set file buffer size
var fileBufSize = 524288

var AddFileCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return AddFile(cmd, args)
	},
}

func AddFile(cmd *cobra.Command, args []string) error {
	err := client.Transmit(config, func(c api.BiblioClient) error {
		stream, err := c.AddFile(context.Background())
		if err != nil {
			return fmt.Errorf("could not create a grpc stream: %w", err)
		}

		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, fileBufSize)

		for {
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("cannot read chunk to buffer: %w", err)
			}

			req := &api.AddFileRequest{Chunk: buf[:n]}

			if err = stream.Send(req); err != nil {
				return fmt.Errorf("cannot send chunk to server: %w", err)
			}
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}

			return fmt.Errorf("could not close a grpc stream: %w", err)
		}

		cmd.Printf(res.GetSha256())

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
