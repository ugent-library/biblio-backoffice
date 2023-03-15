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
	rootCmd.AddCommand(PublicationRelateDatasetCmd)
}

var PublicationRelateDatasetCmd = &cobra.Command{
	Use:   "relate-dataset [id] [dataset-id]",
	Short: "Add related dataset to publication",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return PublicationRelateDataset(cmd, args)
	},
}

func PublicationRelateDataset(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.RelateRequest{
		One: &api.RelateRequest_PublicationOne{PublicationOne: args[0]},
		Two: &api.RelateRequest_DatasetTwo{DatasetTwo: args[1]},
	}
	res, err := c.Relate(context.Background(), req)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errors.New(st.Message())
		}
	}

	if ge := res.GetError(); ge != nil {
		sre := status.FromProto(ge)
		cmd.Printf("%s", sre.Message())
	} else {
		cmd.Printf("%s", res.GetMessage())
	}

	return nil
}
