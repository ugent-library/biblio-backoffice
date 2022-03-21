package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
)

const (
	defaultPgConn        = "host=localhost dbname=biblio-backend sslmode=disable"
	defaultMode          = "production"
	defaultHost          = ""
	defaultPort          = 3000
	defaultSessionName   = "biblio-backend"
	defaultSessionMaxAge = 86400 * 30 // 30 days
	defaultCSRFName      = "biblio-backend.csrf-token"
)

var rootCmd = &cobra.Command{
	Use:   "biblio-backend [command]",
	Short: "The biblio-backend CLI",
	// flags override env vars
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				viper.Set(f.Name, f.Value.String())
			}
		})
		return nil
	},
}

func init() {
	viper.SetEnvPrefix("biblio-backend")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("pg-conn", defaultPgConn)
	viper.SetDefault("mode", defaultMode)
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("session-name", defaultSessionName)
	viper.SetDefault("session-max-age", defaultSessionMaxAge)
	viper.SetDefault("csrf-name", defaultCSRFName)
	viper.SetDefault("temporal-host-port", client.DefaultHostPort)

	rootCmd.PersistentFlags().String("librecat-url", "", "librecat rest api url")
	rootCmd.PersistentFlags().String("librecat-username", "", "librecat rest api username")
	rootCmd.PersistentFlags().String("librecat-password", "", "librecat rest api password")

	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "postgres connection string")

	rootCmd.PersistentFlags().String("temporal-host-port", client.DefaultHostPort, "temporal server host and port")

	rootCmd.PersistentFlags().String("orcid-client-id", "", "orcid client id")
	rootCmd.PersistentFlags().String("orcid-client-secret", "", "orcid client secret")
	rootCmd.PersistentFlags().Bool("orcid-sandbox", false, "use the orcid sandbox in development")

	rootCmd.PersistentFlags().String("oidc-url", "", "openid connect url")
	rootCmd.PersistentFlags().String("oidc-client-id", "", "openid connect client id")
	rootCmd.PersistentFlags().String("oidc-client-secret", "", "openid connect client secret")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
