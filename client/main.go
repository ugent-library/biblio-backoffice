package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/client/cmd"
)

const (
	defaultHost = ""
	defaultPort = 30000
)

func main() {
	viper.SetEnvPrefix("biblio-backend")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	viper.SetDefault("insecure", false)
	// viper.SetDefault("api-tls-disabled", false)
	// viper.SetDefault("api-ca-cert", "")

	rootCmd.PersistentFlags().String("host", defaultHost, "api server host")
	rootCmd.PersistentFlags().Int("port", defaultPort, "api server port")
	rootCmd.PersistentFlags().String("username", "ddd", "api server user username")
	rootCmd.PersistentFlags().String("password", "", "api server user password")
	rootCmd.PersistentFlags().Bool("insecure", false, "disable api client TLS")
	// rootCmd.PersistentFlags().Bool("api-tls-disabled", false, "api client TLS enabled")
	// rootCmd.PersistentFlags().String("api-tls-ca-cert", "", "api client location of the CA certificate")

	fileCmd := (&cmd.FileCmd{}).Command()
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand((&cmd.GetFileCMd{}).Command())
	fileCmd.AddCommand((&cmd.AddFileCMd{}).Command())

	publicationCmd := (&cmd.PublicationCmd{}).Command()
	rootCmd.AddCommand(publicationCmd)
	publicationCmd.AddCommand((&cmd.GetPublicationCmd{}).Command())
	publicationCmd.AddCommand((&cmd.GetAllPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.SearchPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.UpdatePublicationCmd{}).Command())
	publicationCmd.AddCommand((&cmd.AddPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.ImportPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.GetPublicationHistoryCmd{}).Command())
	publicationCmd.AddCommand((&cmd.PurgePublicationCmd{}).Command())
	publicationCmd.AddCommand((&cmd.PurgeAllPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.ValidatePublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.PublicationRelateDatasetCmd{}).Command())

	datasetCmd := (&cmd.DatasetCmd{}).Command()
	rootCmd.AddCommand(datasetCmd)
	datasetCmd.AddCommand((&cmd.GetDatasetCmd{}).Command())
	datasetCmd.AddCommand((&cmd.GetAllDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.SearchDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.UpdateDatasetCmd{}).Command())
	datasetCmd.AddCommand((&cmd.AddDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.ImportDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.GetDatasetHistoryCmd{}).Command())
	datasetCmd.AddCommand((&cmd.PurgeDatasetCmd{}).Command())
	datasetCmd.AddCommand((&cmd.PurgeAllDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.ValidateDatasetsCmd{}).Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
var rootCmd = &cobra.Command{
	Use:   "biblio-client [command]",
	Short: "biblio client",
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
