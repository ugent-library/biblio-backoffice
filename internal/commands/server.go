package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/services"
	"github.com/ugent-library/biblio-backend/services/orcidworker"
	"github.com/ugent-library/biblio-backend/services/webapp"
)

func init() {
	serverCmd.PersistentFlags().String("base-url", "", "base url")

	serverStartCmd.Flags().String("mode", defaultMode, "server mode (development, staging or production)")
	serverStartCmd.Flags().String("host", defaultHost, "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")
	serverStartCmd.Flags().String("csrf-name", "", "csrf cookie name")
	serverStartCmd.Flags().String("csrf-secret", "", "csrf cookie secret")

	// serverCmd.AddCommand(serverRoutesCmd)
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backend HTTP server",
}

// var serverRoutesCmd = &cobra.Command{
// 	Use:   "routes",
// 	Short: "print routes",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		router := buildRouter()
// 		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
// 			hostTemplate, err := route.GetHostTemplate()
// 			if err == nil {
// 				fmt.Println("HOST:", hostTemplate)
// 			}
// 			pathTemplate, err := route.GetPathTemplate()
// 			if err == nil {
// 				fmt.Println("ROUTE:", pathTemplate)
// 			}
// 			pathRegexp, err := route.GetPathRegexp()
// 			if err == nil {
// 				fmt.Println("Path regexp:", pathRegexp)
// 			}
// 			queriesTemplates, err := route.GetQueriesTemplates()
// 			if err == nil {
// 				fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
// 			}
// 			queriesRegexps, err := route.GetQueriesRegexp()
// 			if err == nil {
// 				fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
// 			}
// 			methods, err := route.GetMethods()
// 			if err == nil {
// 				fmt.Println("Methods:", strings.Join(methods, ","))
// 			}
// 			fmt.Println()
// 			return nil
// 		})
// 	},
// }

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		wa, err := webapp.New(e)
		if err != nil {
			log.Fatal(err)
		}
		ow, err := orcidworker.New(e)
		if err != nil {
			log.Fatal(err)
		}
		if err = services.Start(ow, wa); err != nil {
			log.Fatal(err)
		}
	},
}
