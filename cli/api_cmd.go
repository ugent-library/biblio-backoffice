package cli

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/server"
)

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(apiStartCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API commands",
}

var apiStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		users := server.Users{
			&server.User{
				Username: config.AdminUsername,
				Password: config.AdminPassword,
				Role:     "admin",
			},
			&server.User{
				Username: config.CuratorUsername,
				Password: config.CuratorPassword,
				Role:     "curator",
			},
		}
		srv := server.New(newServices(), users, zapLogger)
		addr := fmt.Sprintf("%s:%d", config.API.Host, config.API.Port)
		zapLogger.Infof("Listening at %s", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		return srv.Serve(listener)
	},
}
