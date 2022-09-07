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
	defaultAPIHost = ""
	defaultAPIPort = 30000
)

func main() {
	viper.SetEnvPrefix("biblio-backend")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("api-host", defaultAPIHost)
	viper.SetDefault("api-port", defaultAPIPort)
	viper.SetDefault("api-username", "")
	viper.SetDefault("api-password", "")

	rootCmd.PersistentFlags().String("api-host", defaultAPIHost, "api server host")
	rootCmd.PersistentFlags().Int("api-port", defaultAPIPort, "api server port")
	rootCmd.PersistentFlags().String("api-username", "ddd", "api server user username")
	rootCmd.PersistentFlags().String("api-password", "", "api server user password")

	fileCmd := (&cmd.FileCmd{}).Command()
	fileCmd.AddCommand((&cmd.GetFileCMd{}).Command())
	fileCmd.AddCommand((&cmd.AddFileCMd{}).Command())

	datasetCmd := (&cmd.DatasetCmd{}).Command()
	rootCmd.AddCommand(datasetCmd)
	datasetCmd.AddCommand((&cmd.GetDatasetCmd{}).Command())
	datasetCmd.AddCommand((&cmd.GetAllDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.SearchDatasetsCmd{}).Command())
	datasetCmd.AddCommand((&cmd.UpdateDatasetCmd{}).Command())
	datasetCmd.AddCommand((&cmd.AddDatasetsCmd{}).Command())

	publicationCmd := (&cmd.PublicationCmd{}).Command()
	rootCmd.AddCommand(publicationCmd)
	publicationCmd.AddCommand((&cmd.GetPublicationCmd{}).Command())
	publicationCmd.AddCommand((&cmd.GetAllPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.SearchPublicationsCmd{}).Command())
	publicationCmd.AddCommand((&cmd.UpdatePublicationCmd{}).Command())
	publicationCmd.AddCommand((&cmd.AddPublicationsCmd{}).Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
var rootCmd = &cobra.Command{
	Use:   "api [command]",
	Short: "biblio api client",
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
