package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaller   = protojson.MarshalOptions{UseProtoNames: true}
	unmarshaller = protojson.UnmarshalOptions{}
	configFile   string
	config       cnx.Config
)

const (
	defaultHost = ""
	defaultPort = 443
)

func init() {
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	viper.SetDefault("insecure", false)
	viper.SetDefault("cacert", "")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		cobra.CheckErr(viper.ReadInConfig())
	}

	viper.SetEnvPrefix("biblio-backoffice")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cobra.CheckErr(viper.Unmarshal(&config))
}

var rootCmd = &cobra.Command{
	Use: "biblio-client",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Silence the usage text if an error occurs
		cmd.SilenceUsage = true
		return nil
	},
}

func Execute() error {
	// Set the output to os.Stdout. If not set, cmd.Println would write to Stderr
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}
