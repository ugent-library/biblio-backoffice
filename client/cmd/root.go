package cmd

import (
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaller   = protojson.MarshalOptions{UseProtoNames: true}
	unmarshaller = protojson.UnmarshalOptions{}
	configFile   string
	config       client.Config
)

const (
	defaultHost = ""
	defaultPort = 443
)

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	viper.SetDefault("insecure", false)
	viper.SetDefault("cacert", "")

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.biblioclient.toml)")
	rootCmd.PersistentFlags().String("host", defaultHost, "api server host")
	rootCmd.PersistentFlags().Int("port", defaultPort, "api server port")
	rootCmd.PersistentFlags().String("username", "ddd", "api server user username")
	rootCmd.PersistentFlags().String("password", "", "api server user password")
	rootCmd.PersistentFlags().Bool("insecure", false, "disable TLS encryption")
	rootCmd.PersistentFlags().String("cacert", "", "path to the CA certificate (if using self-signed certificates)")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.)
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".biblioclient.{yaml,toml}".
		viper.AddConfigPath(home)
		viper.SetConfigName(".biblioclient")
	}

	viper.SetEnvPrefix("biblio-backoffice")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
var rootCmd = &cobra.Command{
	Use:   "biblio-client [command]",
	Short: "biblio client",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Silence the usage text if an error occurs
		cmd.SilenceUsage = true

		// flags override env vars
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				viper.BindPFlag(f.Name, f)
			}
		})

		cobra.CheckErr(viper.ReadInConfig())
		cobra.CheckErr(viper.Unmarshal(&config))

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
