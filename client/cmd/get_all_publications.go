package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(GetAllPublicationsCmd)
}

var GetAllPublicationsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetAllPublications(cmd, args)
	},
}

func GetAllPublications(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.GetAllPublicationsRequest{}
		stream, err := c.GetAllPublications(context.Background(), req)
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

			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.Printf("%s\n", sre.Message())
			}

			if rr := res.GetPublication(); rr != nil {
				cmd.Printf("%s\n", rr.GetPayload())
			}
		}

		return nil
	})
}
