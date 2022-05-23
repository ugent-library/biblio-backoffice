package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
)

func init() {
	fileCmd.AddCommand(fileAddCmd)
	fileCmd.AddCommand((fileAddManyCmd))
	rootCmd.AddCommand(fileCmd)
}

func newFileStore() (*filestore.Store, error) {
	fs, fsErr := filestore.New(viper.GetString("file-dir"))
	if fsErr != nil {
		return nil, fmt.Errorf("unable to initialize filestore: %s", fsErr)
	}
	return fs, nil
}

func addFile(fs *filestore.Store, path string) (string, error) {
	fh, fhErr := os.Open(path)
	defer fh.Close()
	if fhErr != nil {
		return "", fmt.Errorf("unable to %s for reading: %v", path, fhErr)
	}
	id, addErr := fs.Add(fh)
	if addErr != nil {
		return "", fmt.Errorf("unable to add file %s: %v", path, addErr)
	}
	return id, nil
}

var fileCmd = &cobra.Command{
	Use:   "file [command]",
	Short: "File commands",
}

var fileAddCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add file by path",
	Long: `
	Adds one file to the filestore.
	File provided is the file added to the filestore.

	Writes id and path to the stdout:

		<id> <path>

	Can easily be checked as following:
		$ ./biblio-backend file add /path/to/file.txt > /path/to/id.txt
		$ sha256sum -c /path/to/id.txt
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		fs, fsErr := newFileStore()
		if fsErr != nil {
			fmt.Fprintf(os.Stderr, "unable to add file %s: %s\n", path, fsErr.Error())
			os.Exit(1)
		}
		id, addErr := addFile(fs, path)
		if addErr != nil {
			fmt.Fprintf(os.Stderr, "unable to add file %s: %s\n", path, addErr.Error())
			os.Exit(1)
		}
		fmt.Printf("%s %s\n", id, path)
	},
}

var fileAddManyCmd = &cobra.Command{
	Use:   "add_many [file]",
	Short: "Add multiple files",
	Long: `
	Add many files to the filestore.
	Expected is txt file containing a file path per line.
	Reads from the stdin when file is not provided. 
	For each path the new id followed by the old path is printed to the stdout:
		<id> <path>

	Can easily be checked as following:
		$ ./biblio-backend file add_many < /path/to/file_paths.txt > /path/to/ids.txt
		$ sha256sum -c /path/to/ids.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// file add < /path/to/files_paths.txt
		var fhIn *os.File = os.Stdin
		var fhInErr error

		// file add /path/to/files_paths.txt
		if len(args) > 0 {
			fhIn, fhInErr = os.Open(args[0])
			if fhInErr != nil {
				fmt.Fprintf(os.Stderr, "Unable to open file %s: %s", args[0], fhInErr.Error())
				os.Exit(1)
			}
		}

		fs, fsErr := newFileStore()
		if fsErr != nil {
			fmt.Fprint(os.Stderr, fsErr.Error())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(fhIn)
		for scanner.Scan() {
			path := scanner.Text()
			id, addErr := addFile(fs, path)
			if addErr != nil {
				fmt.Fprintf(os.Stderr, "unable add file %s : %s\n", path, addErr.Error())
				continue
			}
			// <file-id> <old-path>
			fmt.Printf("%s %s\n", id, path)
		}
	},
}
