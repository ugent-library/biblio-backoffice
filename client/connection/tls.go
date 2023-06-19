package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials(config Config) (credentials.TransportCredentials, error) {
	// Use system CA certificates
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	// Override if a custom CA cert is provided
	if config.Cacert != "" {
		pemServerCA, err := os.ReadFile(config.Cacert)
		if err != nil {
			return nil, err
		}

		certPool = x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			return nil, fmt.Errorf("failed to add server CA's certificate")
		}
	}

	// Create the credentials and return it
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
