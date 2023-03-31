package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	FileCmd.AddCommand(GetFileCmd)
}

var GetFileCmd = &cobra.Command{
	Use:   "get [sha256]",
	Short: "Get file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetFile(cmd, args)
	},
}

func GetFile(cmd *cobra.Command, args []string) error {
	err := cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.GetFileRequest{Sha256: args[0]}
		stream, err := c.GetFile(context.Background(), req)
		if err != nil {
			return fmt.Errorf("error while reading stream: %v", err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				if st, ok := status.FromError(err); ok {
					return errors.New(st.Message())
				}

				return err
			}

			if _, err := os.Stdout.Write(res.Chunk); err != nil {
				return fmt.Errorf("error writing to stdout: %w", err)
			}
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
