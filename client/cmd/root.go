package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaller   = protojson.MarshalOptions{UseProtoNames: true}
	unmarshaller = protojson.UnmarshalOptions{}
)

const (
	defaultHost = ""
	defaultPort = 443
)

func init() {
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

func Execute() error {
	return rootCmd.Execute()
}
