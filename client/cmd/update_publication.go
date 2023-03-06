package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(UpdatePublicationCmd)
}

var UpdatePublicationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update publication",
	RunE: func(cmd *cobra.Command, args []string) error {
		return UpdatePublication(cmd, args)
	},
}

func UpdatePublication(cmd *cobra.Command, args []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("ContextDeadlineExceeded: true")
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("could not read from stdin: %v", err)
	}

	p := &api.Publication{
		Payload: line,
	}

	req := &api.UpdatePublicationRequest{Publication: p}
	res, err := c.UpdatePublication(ctx, req)

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
}
