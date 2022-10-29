package commands

// func init() {
// 	apiStartCmd.Flags().String("api-host", defaultAPIHost, "api server host")
// 	apiStartCmd.Flags().Int("api-port", defaultAPIPort, "api server port")
// 	apiStartCmd.Flags().String("api-username", "", "api server administrator username")
// 	apiStartCmd.Flags().String("api-password", "", "api server administrator password")
// 	apiStartCmd.Flags().Bool("api-tls-enabled", false, "api server enable TLS encryped connections")
// 	apiStartCmd.Flags().String("api-tls-servercert", "", "api server location of server certificate file")
// 	apiStartCmd.Flags().String("api-tls-serverkey", "", "api server location of server certificate key file")

// 	apiCmd.AddCommand(apiStartCmd)
// 	rootCmd.AddCommand(apiCmd)
// }

// var apiCmd = &cobra.Command{
// 	Use:   "api [command]",
// 	Short: "api commands",
// }

// var apiStartCmd = &cobra.Command{
// 	Use:   "start",
// 	Short: "start the api server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		logger := newLogger()

// 		srv := server.New(Services(), logger)
// 		addr := fmt.Sprintf("%s:%d", viper.GetString("api-host"), viper.GetInt("api-port"))
// 		log.Printf("Listening at %s", addr)
// 		listener, err := net.Listen("tcp", addr)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if err := srv.Serve(listener); err != nil {
// 			log.Fatal(err)
// 		}
// 	},
// }
