package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("base-url", defaultBaseURL)

	rootCmd.PersistentFlags().String("librecat-url", "", "librecat rest api url")
	rootCmd.PersistentFlags().String("librecat-username", "", "librecat rest api username")
	rootCmd.PersistentFlags().String("librecat-password", "", "librecat rest api password")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
