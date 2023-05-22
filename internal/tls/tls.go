package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/vipulvpatil/candidate-tracker-go/internal/config"
	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials(cfg *config.Config) (credentials.TransportCredentials, error) {
	caCert, err := loadStringFromBase64EncodedEnvVar(cfg.CaCertBase64)
	if err != nil {
		return nil, err
	}

	serverCert, err := loadStringFromBase64EncodedEnvVar(cfg.ServerCertBase64)
	if err != nil {
		return nil, err
	}

	serverKey, err := loadStringFromBase64EncodedEnvVar(cfg.ServerKeyBase64)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA's certificate")
	}

	// Load server's certificate and private key
	serverCertAndKey, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCertAndKey},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func loadStringFromBase64EncodedEnvVar(envVarValue string) ([]byte, error) {
	cert, err := base64.StdEncoding.DecodeString(envVarValue)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
