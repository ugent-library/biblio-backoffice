package cmd

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(TransferPublicationsCmd)
}

var TransferPublicationsCmd = &cobra.Command{
	Use:   "transfer UID UID [PUBID]",
	Short: "Transfer publications between people",
	Long: `
	Transfer one or multiple publications between two persons.

	Each person id identified by an UUID. The first argument is the source, the second argument is the target person.
	Transferring a publication means replacing all matching instances of the source ID with the target's ID across all
	publicatoin fields (user, last_user & contributor fields).

	This operation transfers the current and all previous snapshots of a publication between persons.

	A publication ID can be passed as an optional third argument. If no publication ID is passed, the transfer will
	happen across all stored publications. If a publication ID is passed, the transfer command will be limited to that
	specific stored publication.

	The command outputs either a success message or an error message to stdout:

		$ ./biblio-client publication transfer UID UID
		p: ID: s: SNAPSHOT-ID ::: creator: UID -> UID
		p: ID: s: SNAPSHOT-ID ::: supervisor: UID -> UID
		p: ID: s: SNAPSHOT-ID ::: editor: UID -> UID

		$ ./biblio-client publication transfer UID UID
		Error: could not retrieve person UID: record not found

	If no matching instances of the source UID could be found, the transfer command won't produce any output.
	`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return TransferPublications(cmd, args)
	},
}

func TransferPublications(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		source := args[0]
		dest := args[1]

		pubid := ""
		if len(args) > 2 {
			pubid = args[2]
		}

		req := &api.TransferPublicationsRequest{
			Src:           source,
			Dest:          dest,
			Publicationid: pubid,
		}

		stream, err := c.TransferPublications(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		stream.CloseSend()

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// read done.
				break
			}

			// return gRPC level error
			if err != nil {
				if st, ok := status.FromError(err); ok {
					return errors.New(st.Message())
				}

				return err
			}

			// Application level error
			if ge := res.GetError(); ge != nil {
				sre := status.FromProto(ge)
				cmd.Printf("%s\n", sre.Message())
			}

			if rr := res.GetMessage(); rr != "" {
				cmd.Printf("%s\n", rr)
			}
		}

		return nil
	})
}
