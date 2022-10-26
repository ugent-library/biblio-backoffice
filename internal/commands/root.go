package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/ugent-library/biblio-backend/internal/snapstore"
)

const (
	defaultPgConn               = "postgres://localhost:5432/biblio_backend?sslmode=disable"
	defaultEs6URL               = "http://localhost:9200"
	defaultDatasetIndex         = "biblio_backend_datasets"
	defaultPublicationIndex     = "biblio_backend_publications"
	defaultAPIHost              = ""
	defaultAPIPort              = 30000
	defaultAPIAdminUsername     = "admin"
	defaultAPIAdminPassword     = "admin"
	defaultMode                 = "production"
	defaultHost                 = ""
	defaultPort                 = 3000
	defaultSessionName          = "biblio-backend"
	defaultSessionMaxAge        = 86400 * 30 // 30 days
	defaultCSRFName             = "biblio-backend.csrf-token"
	defaultHandleServerURL      = "http://localhost:4000/handles"
	defaultHandleServerPrefix   = "1854"
	defaultHandleServerUsername = "handle"
	defaultHandleServerPassword = "handle"
)

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
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
	viper.SetDefault("es6-url", defaultEs6URL)
	viper.SetDefault("dataset-index", defaultDatasetIndex)
	viper.SetDefault("publication-index", defaultPublicationIndex)
	viper.SetDefault("api-host", defaultAPIHost)
	viper.SetDefault("api-port", defaultAPIPort)
	viper.SetDefault("api-username", defaultAPIAdminUsername)
	viper.SetDefault("api-password", defaultAPIAdminPassword)
	viper.SetDefault("mode", defaultMode)
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("session-name", defaultSessionName)
	viper.SetDefault("session-max-age", defaultSessionMaxAge)
	viper.SetDefault("csrf-name", defaultCSRFName)
	viper.SetDefault("hdl-srv-url", defaultHandleServerURL)
	viper.SetDefault("hdl-srv-prefix", defaultHandleServerPrefix)
	viper.SetDefault("hdl-srv-username", defaultHandleServerUsername)
	viper.SetDefault("hdl-srv-password", defaultHandleServerPassword)

	rootCmd.PersistentFlags().String("file-dir", "", "file store root directory")

	rootCmd.PersistentFlags().String("frontend-url", "", "biblio frontend url")
	rootCmd.PersistentFlags().String("frontend-username", "", "biblio frontend username")
	rootCmd.PersistentFlags().String("frontend-password", "", "biblio frontend password")

	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "postgres connection string")
	rootCmd.PersistentFlags().String("es6-url", defaultEs6URL, "elasticsearch 6.x url, separate multiple with comma")
	rootCmd.PersistentFlags().String("dataset-index", defaultDatasetIndex, "dataset index name")
	rootCmd.PersistentFlags().String("publication-index", defaultPublicationIndex, "publication index name")

	rootCmd.PersistentFlags().String("orcid-client-id", "", "orcid client id")
	rootCmd.PersistentFlags().String("orcid-client-secret", "", "orcid client secret")
	rootCmd.PersistentFlags().Bool("orcid-sandbox", false, "use the orcid sandbox in development")

	rootCmd.PersistentFlags().String("oidc-url", "", "openid connect url")
	rootCmd.PersistentFlags().String("oidc-client-id", "", "openid connect client id")
	rootCmd.PersistentFlags().String("oidc-client-secret", "", "openid connect client secret")

	rootCmd.PersistentFlags().String("citeproc-url", "", "citeproc url")

	// rootCmd.PersistentFlags().String("imagor-url", "", "imagor url")
	// rootCmd.PersistentFlags().String("imagor-secret", "", "imagor secret")
	rootCmd.PersistentFlags().String("hdl-srv-url", defaultHandleServerURL, "handle server base url (without trailing slash)")
	rootCmd.PersistentFlags().String("hdl-srv-prefix", defaultHandleServerPrefix, "handle server base prefix")
	rootCmd.PersistentFlags().String("hdl-srv-username", defaultHandleServerUsername, "handle server auth basic username")
	rootCmd.PersistentFlags().String("hdl-srv-password", defaultHandleServerPassword, "handle server auth basic password")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
