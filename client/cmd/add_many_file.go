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
	FileCmd.AddCommand(AddManyFileCmd)
}

var AddManyFileCmd = &cobra.Command{
	Use:   "add-many [file]",
	Short: "Add multiple files",
	Long: `
	Add many files to the filestore.
	Expected is txt file containing a file path per line.
	Reads from the stdin when file is not provided.
	For each path the new id followed by the old path is printed to the stdout:

		<id> <path>

	Can easily be checked as following:

		$ ./biblio-backoffice file add_many < /path/to/file_paths.txt > /path/to/ids.txt
		$ sha256sum -c /path/to/ids.txt
	`,
	RunE: AddManyFile,
}

func AddManyFile(cmd *cobra.Command, args []string) error {
	// file add < /path/to/files_paths.txt
	var fhIn = os.Stdin
	var fhInErr error

	// file add /path/to/files_paths.txt
	if len(args) > 0 {
		fhIn, fhInErr = os.Open(args[0])
		if fhInErr != nil {
			return fmt.Errorf("unable to open file %s: %w", args[0], fhInErr)
		}
	}

	var txErr error

	scanner := bufio.NewScanner(fhIn)
	for scanner.Scan() {
		path := scanner.Text()

		f, err := os.Open(path)
		if err != nil {
			cmd.Printf("cannot open file: %s\n", err.Error())
			continue
		}
		defer f.Close()

		txErr = cnx.Handle(config, func(c api.BiblioClient) error {
			stream, err := c.AddFile(context.Background())
			if err != nil {
				return fmt.Errorf("could not create a grpc stream: %w", err)
			}

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

			cmd.Printf("%s %s\n", res.GetSha256(), path)

			return nil
		})

		if txErr != nil {
			return txErr
		}
	}

	return nil
}
