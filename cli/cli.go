package cli

import (
	"log"
	"strings"
	_ "time/tzdata"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/ugent-library/biblio-backoffice/internal/snapstore"
)

const (
	defaultTimezone            = "Europe/Brussels"
	defaultIndexRetention      = 2
	defaultAPIPort             = 30000
	defaultMode                = "production"
	defaultPort                = 3000
	defaultSessionName         = "biblio-backoffice"
	defaultSessionMaxAge       = 86400 * 30 // 30 days
	defaultCSRFName            = "biblio-backoffice.csrf-token"
	defaultHandleServerEnabled = false
	defaultMaxFileSize         = 2_000_000_000
)

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
var rootCmd = &cobra.Command{
	Use:   "biblio-backoffice [command]",
	Short: "The biblio-backoffice CLI",
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
	viper.SetEnvPrefix("biblio-backoffice")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("timezone", defaultTimezone)
	viper.SetDefault("index-retention", defaultIndexRetention)
	viper.SetDefault("s3-region", "us-east-1")
	viper.SetDefault("mode", defaultMode)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("session-name", defaultSessionName)
	viper.SetDefault("session-max-age", defaultSessionMaxAge)
	viper.SetDefault("csrf-name", defaultCSRFName)
	viper.SetDefault("max-file-size", defaultMaxFileSize)

	rootCmd.PersistentFlags().String("file-dir", "", "file store root directory")
	rootCmd.PersistentFlags().String("s3-endpoint", "", "S3 endpoint url")
	rootCmd.PersistentFlags().String("s3-region", "", "S3 region")
	rootCmd.PersistentFlags().String("s3-id", "", "S3 access key id")
	rootCmd.PersistentFlags().String("s3-secret", "", "S3 secret access key")
	rootCmd.PersistentFlags().String("s3-bucket", "", "S3 file bucket name")
	rootCmd.PersistentFlags().String("s3-temp-bucket", "", "S3 temp file bucket name")

	rootCmd.PersistentFlags().String("frontend-url", "", "biblio frontend url")
	rootCmd.PersistentFlags().String("frontend-username", "", "biblio frontend username")
	rootCmd.PersistentFlags().String("frontend-password", "", "biblio frontend password")

	rootCmd.PersistentFlags().String("pg-conn", "", "postgres connection string")
	rootCmd.PersistentFlags().String("es6-url", "", "elasticsearch 6.x url, separate multiple with comma")
	rootCmd.PersistentFlags().String("dataset-index", "", "dataset index name")
	rootCmd.PersistentFlags().String("publication-index", "", "publication index name")
	rootCmd.PersistentFlags().Int("index-retention", defaultIndexRetention, "number of old indexes to retain after index switch")

	rootCmd.PersistentFlags().String("frontend-es6-url", "", "frontend elasticsearch 6.x url, separate multiple with comma")

	rootCmd.PersistentFlags().String("orcid-client-id", "", "orcid client id")
	rootCmd.PersistentFlags().String("orcid-client-secret", "", "orcid client secret")
	rootCmd.PersistentFlags().Bool("orcid-sandbox", false, "use the orcid sandbox in development")

	rootCmd.PersistentFlags().String("oidc-url", "", "openid connect url")
	rootCmd.PersistentFlags().String("oidc-client-id", "", "openid connect client id")
	rootCmd.PersistentFlags().String("oidc-client-secret", "", "openid connect client secret")

	rootCmd.PersistentFlags().String("citeproc-url", "", "citeproc url")

	rootCmd.PersistentFlags().String("mongodb-url", "", "mongodb connection uri (for authority database)")

	// rootCmd.PersistentFlags().String("imagor-url", "", "imagor url")
	// rootCmd.PersistentFlags().String("imagor-secret", "", "imagor secret")
	rootCmd.PersistentFlags().Bool("hdl-srv-enabled", false, "enable updates to handle server (disabled by default)")
	rootCmd.PersistentFlags().String("hdl-srv-url", "", "handle server base url (without trailing slash)")
	rootCmd.PersistentFlags().String("hdl-srv-prefix", "", "handle server base prefix")
	rootCmd.PersistentFlags().String("hdl-srv-username", "", "handle server auth basic username")
	rootCmd.PersistentFlags().String("hdl-srv-password", "", "handle server auth basic password")

	rootCmd.PersistentFlags().Int("max-file-size", defaultMaxFileSize, "maximum file size")
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
