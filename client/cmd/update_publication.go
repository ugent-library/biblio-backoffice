package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(UpdatePublicationCmd)
}

var UpdatePublicationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update publication",
	Long: `
	Update one or multiple publications.

	This command reads a JSONL formatted file from stdin and streams it to the store.

	It will output either a success message or an error message per record:

		$ ./biblio-backoffice publication update < publications.jsonl
	`,
	RunE: UpdatePublication,
}

func UpdatePublication(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		reader := bufio.NewReader(cmd.InOrStdin())
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("could not read from stdin: %v", err)
		}

		p := &api.Publication{
			Payload: line,
		}

		req := &api.UpdatePublicationRequest{Publication: p}
		res, err := c.UpdatePublication(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s\n", sre.Message())
		}

		return nil
	})
}
