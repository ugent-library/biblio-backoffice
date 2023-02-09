package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/client/cmd"
)

const (
	defaultHost = ""
	defaultPort = 443
)

func main() {
	viper.SetEnvPrefix("biblio-backoffice")
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

	fileCmd := cmd.FileCmd
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(cmd.GetFileCMd)
	fileCmd.AddCommand(cmd.AddFileCmd)

	publicationCmd := cmd.PublicationCmd
	rootCmd.AddCommand(publicationCmd)
	publicationCmd.AddCommand(cmd.GetPublicationCmd)
	publicationCmd.AddCommand(cmd.GetAllPublicationsCmd)
	publicationCmd.AddCommand(cmd.SearchPublicationsCmd)
	publicationCmd.AddCommand(cmd.UpdatePublicationCmd)
	publicationCmd.AddCommand((cmd.AddPublicationsCmd))
	publicationCmd.AddCommand(cmd.ImportPublicationsCmd)
	publicationCmd.AddCommand(cmd.GetPublicationHistoryCmd)
	publicationCmd.AddCommand(cmd.PurgePublicationCmd)
	publicationCmd.AddCommand(cmd.PurgeAllPublicationsCmd)
	publicationCmd.AddCommand(cmd.ValidatePublicationsCmd)
	publicationCmd.AddCommand(cmd.PublicationRelateDatasetCmd)

	datasetCmd := cmd.DatasetCmd
	rootCmd.AddCommand(datasetCmd)
	datasetCmd.AddCommand(cmd.GetDatasetCmd)
	datasetCmd.AddCommand(cmd.GetAllDatasetsCmd)
	datasetCmd.AddCommand(cmd.SearchDatasetsCmd)
	datasetCmd.AddCommand(cmd.UpdateDatasetCmd)
	datasetCmd.AddCommand(cmd.AddDatasetsCmd)
	datasetCmd.AddCommand(cmd.ImportDatasetsCmd)
	datasetCmd.AddCommand(cmd.GetDatasetHistoryCmd)
	datasetCmd.AddCommand(cmd.PurgeDatasetCmd)
	datasetCmd.AddCommand(cmd.PurgeAllDatasetsCmd)
	datasetCmd.AddCommand(cmd.ValidateDatasetsCmd)

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
