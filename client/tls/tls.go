package tls

// func LoadTLSCredentials() (credentials.TransportCredentials, error) {
// 	// Use system CA certificates
// 	certPool, err := x509.SystemCertPool()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Override if a custom CA cert is provided
// 	if cacert := viper.GetString("api-tls-ca-cert"); cacert != "" {
// 		pemServerCA, err := os.ReadFile(cacert)
// 		if err != nil {
// 			return nil, err
// 		}

// 		certPool = x509.NewCertPool()
// 		if !certPool.AppendCertsFromPEM(pemServerCA) {
// 			return nil, fmt.Errorf("failed to add server CA's certificate")
// 		}
// 	}

// 	// Create the credentials and return it
// 	config := &tls.Config{
// 		RootCAs: certPool,
// 	}

// 	return credentials.NewTLS(config), nil
// }
