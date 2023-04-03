package cmd

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
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

type importFile struct {
	File   string `json:"file,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
	Size   int64  `json:"size,omitempty"`
}

func init() {
	FileCmd.AddCommand(ImportFileCmd)
}

var ImportFileCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import multiple files",
	Long: `
	Add many files to the filestore.
	Expected is json file containing "file" and "sha256":

		{ "file": "/path/to/file.txt", "sha256": "<sha256>" }

	Reads from the stdin when file is not provided.
	For each path the new id followed by the old path is printed to the stdout:
		<id> <path>

	Can easily be checked as following:
		$ ./biblio-backoffice file import < /path/to/file_paths.json > /path/to/ids.txt
		$ sha256sum -c /path/to/ids.txt

	If the filestore already contains an identical checksum, the import will be skipped.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ImportFile(cmd, args)
	},
}

func ImportFile(cmd *cobra.Command, args []string) error {
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

	decoder := json.NewDecoder(bufio.NewReader(fhIn))
	lineNo := 0

	for {
		var txErr error

		importFile := importFile{}
		if err := decoder.Decode(&importFile); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("unable to decode line %d", lineNo)
		}

		if errs := importFile.Validate(); len(errs) > 0 {
			for _, e := range errs {
				cmd.Printf("validation error for line %d: %s\n", lineNo, e.Error())
			}
			lineNo++
			continue
		}

		f, err := os.Open(importFile.File)
		if err != nil {
			return fmt.Errorf("cannot open file: %s", err.Error())
		}
		defer f.Close()

		// checksum check before sending file. Make sure it's the same file
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return fmt.Errorf("could not calculate sha256 signature: %s", err.Error())
		}

		sha256res := fmt.Sprintf("%x", h.Sum(nil))

		if sha256res != importFile.Sha256 {
			cmd.Printf("sha256 do not match: expect %s, got %s for file %s", importFile.Sha256, sha256res, importFile.File)
			lineNo++
			continue
		}

		var fileExists bool

		txErr = cnx.Handle(config, func(c api.BiblioClient) error {
			req := &api.ExistsFileRequest{Sha256: importFile.Sha256}
			res, err := c.ExistsFile(context.Background(), req)

			if err != nil {
				if st, ok := status.FromError(err); ok {
					return errors.New(st.Message())
				}
			}

			fileExists = res.Exists

			return nil
		})

		if errors.Is(txErr, context.DeadlineExceeded) {
			log.Fatal("ContextDeadlineExceeded: true")
		}

		// skip files that are already in the store
		if fileExists {
			// <file-id> <old-path>
			cmd.Printf("%s %s\n", importFile.Sha256, importFile.File)
			lineNo++
			continue
		}

		txErr = cnx.Handle(config, func(c api.BiblioClient) error {
			stream, err := c.AddFile(context.Background())
			if err != nil {
				return fmt.Errorf("could not create a grpc stream: %w", err)
			}

			f, err := os.Open(importFile.File)
			if err != nil {
				return fmt.Errorf("cannot open file: %s", err.Error())
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

			cmd.Printf("%s %s\n", res.GetSha256(), importFile.File)

			return nil
		})

		if txErr != nil {
			if errors.Is(txErr, context.DeadlineExceeded) {
				log.Fatal("ContextDeadlineExceeded: true")
			}

			cmd.Printf("%s", txErr.Error())
			lineNo++
			continue
		}

		lineNo++
	}

	return nil
}

func (f *importFile) Validate() (errs []error) {
	// file: required
	if f.File == "" {
		errs = append(errs, errors.New("attribute 'file' is required"))
	}
	// size: required
	if f.Size <= 0 {
		errs = append(errs, errors.New("attribute 'size' should be 1 or higher"))
	} else if f.File != "" {
		fhStat, fhStatErr := os.Stat(f.File)
		if fhStatErr != nil {
			errs = append(errs, fmt.Errorf("file \"%s\" not found", f.File))
		} else if fhStat.Size() != f.Size {
			errs = append(
				errs,
				fmt.Errorf(
					"file \"%s\" has different file size (%d <=> %d)",
					f.File,
					f.Size,
					fhStat.Size(),
				),
			)
		}
	}
	// sha256: required
	if len(f.Sha256) != 64 {
		errs = append(errs, errors.New("attribute 'sha256' is required and must be 64 characters long"))
	}
	return
}
