package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	FileCmd.AddCommand(AddFileCmd)
}

var AddFileCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add file by path",
	Long: `
	Add a single file to the filestore.
	File provided is the file that will be added to the filestore.

	Prints id and local path to stdout:

		<id> <path>

	The id is the sha256 checksum of the file. You can use it to check if a
	file was succesfully added to the filestore:

		$ ./biblio-backoffice file add /path/to/file.txt > /path/to/id.txt
		$ sha256sum -c /path/to/id.txt
	`,
	Args: cobra.ExactArgs(1),
	RunE: AddFile,
}

func AddFile(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		stream, err := c.AddFile(context.Background())
		if err != nil {
			return fmt.Errorf("could not create a grpc stream: %w", err)
		}

		path := args[0]
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}
		defer f.Close()

		r := bufio.NewReader(f)

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

		cmd.Printf("%s %s", res.GetSha256(), path)

		return nil
	})
}
