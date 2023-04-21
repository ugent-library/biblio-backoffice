package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	DatasetCmd.AddCommand(PurgeAllDatasetsCmd)
}

var PurgeAllDatasetsCmd = &cobra.Command{
	Use:   "purge-all",
	Short: "Purge all datasets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return PurgeAllDatasets(cmd, args)
	},
}

func init() {
	PurgeAllDatasetsCmd.Flags().BoolP("yes", "y", false, "are you sure?")
}

func PurgeAllDatasets(cmd *cobra.Command, args []string) error {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		cmd.Printf("no confirmation flag set. you need to set the --yes flag")
		return nil
	}

	return cnx.Handle(config, func(c api.BiblioClient) error {
		req := &api.PurgeAllDatasetsRequest{
			Confirm: true,
		}
		res, err := c.PurgeAllDatasets(context.Background(), req)

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return errors.New(st.Message())
			}
		}

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s", sre.Message())
		}

		if res.GetOk() {
			cmd.Printf("purged all datasets")
		}

		return nil
	})
}
