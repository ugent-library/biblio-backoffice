package cli

import (
	"log/slog"
	_ "time/tzdata"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	_ "github.com/ugent-library/biblio-backoffice/snapstore"
	"go.uber.org/zap"
)

var (
	version   Version
	config    Config
	zapLogger *zap.SugaredLogger
	logger    *slog.Logger

	rootCmd = &cobra.Command{
		Use:   "biblio-backoffice",
		Short: "Biblio backoffice CLI",
	}
)

func init() {
	cobra.OnInitialize(initVersion, initConfig, initLogger)
	cobra.OnFinalize(func() {
		zapLogger.Sync()
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

// TODO completely switch to slog logger
func initLogger() {
	if config.Env == "local" {
		l, err := zap.NewDevelopment()
		cobra.CheckErr(err)
		zapLogger = l.Sugar()
	} else {
		l, err := zap.NewProduction()
		cobra.CheckErr(err)
		zapLogger = l.Sugar()
	}

}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}
