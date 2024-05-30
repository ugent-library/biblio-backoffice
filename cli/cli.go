package cli

import (
	"log/slog"
	"os"
	_ "time/tzdata"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	_ "github.com/ugent-library/biblio-backoffice/snapstore"
)

var (
	version Version
	config  Config
	logger  *slog.Logger

	rootCmd = &cobra.Command{
		Use:   "biblio-backoffice",
		Short: "Biblio backoffice CLI",
	}
)

func init() {
	cobra.OnInitialize(initVersion, initConfig, initLogger)
}

func initVersion() {
	cobra.CheckErr(env.Parse(&version))
}

func initConfig() {
	cobra.CheckErr(env.ParseWithOptions(&config, env.Options{
		Prefix: "BIBLIO_BACKOFFICE_",
	}))
}

func initLogger() {
	if config.Env == "local" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}
