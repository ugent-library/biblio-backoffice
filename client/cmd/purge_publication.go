package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(PurgePublicationCmd)
}

var PurgePublicationCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge publication",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return PurgePublication(cmd, args)
	},
}

func PurgePublication(cmd *cobra.Command, args []string) error {
	err := cnx.Handle(config, func(c api.BiblioClient) error {
		id := args[0]
		req := &api.PurgePublicationRequest{Id: id}
		res, err := c.PurgePublication(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}

			return err
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		}

		if res.GetOk() {
			cmd.Printf("purged publication %s", id)
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("ContextDeadlineExceeded: true")
	}

	return err
}
