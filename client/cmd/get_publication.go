package cmd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(GetPublicationCmd)
}

var GetPublicationCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetPublication(cmd, args)
	},
}

func GetPublication(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	id := args[0]
	req := &api.GetPublicationRequest{Id: id}
	res, err := c.GetPublication(ctx, req)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errors.New(st.Message())
		}
	}

	if ge := res.GetError(); ge != nil {
		sre := status.FromProto(ge)
		cmd.Printf("%s", sre.Message())
	} else {
		cmd.Printf("%s", res.GetPublication().GetPayload())
	}

	return nil
}
