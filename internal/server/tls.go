package server

// func loadTLSCredentials() (credentials.TransportCredentials, error) {
// 	certFile := viper.GetString("api-tls-servercert")
// 	keyFile := viper.GetString("api-tls-serverkey")

// 	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create the credentials and return it
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{serverCert},
// 		ClientAuth:   tls.NoClientCert,
// 	}

// 	return credentials.NewTLS(config), nil
// }
