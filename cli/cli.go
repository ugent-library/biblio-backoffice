package cli

import (
	_ "time/tzdata"

	"github.com/caarlos0/env/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/ugent-library/biblio-backoffice/internal/snapstore"
)

var (
	version Version
	config  Config
	logger  *zap.SugaredLogger

	rootCmd = &cobra.Command{
		Use:   "biblio-backoffice",
		Short: "Biblio backoffice CLI",
	}
)

func init() {
	cobra.OnInitialize(initVersion, initConfig, initLogger)
	cobra.OnFinalize(func() {
		logger.Sync()
	})
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
		l, err := zap.NewDevelopment()
		cobra.CheckErr(err)
		logger = l.Sugar()
	} else {
		l, err := zap.NewProduction()
		cobra.CheckErr(err)
		logger = l.Sugar()
	}
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}
