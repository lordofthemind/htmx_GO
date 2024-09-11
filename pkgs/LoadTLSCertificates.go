package pkgs

import (
	"crypto/tls"
	"fmt"
)

// LoadTLSCertificate loads the TLS certificate and private key from the specified files.
// It returns a tls.Certificate object and an error if loading the certificate fails.
//
// Parameters:
// - certFile: The path to the TLS certificate file (usually .crt or .pem).
// - keyFile: The path to the private key file corresponding to the TLS certificate (usually .key).
//
// Returns:
// - tls.Certificate: A certificate object containing the loaded certificate and private key.
// - error: An error if the certificate or key could not be loaded or are invalid.
func LoadTLSCertificate(certFile, keyFile string) (tls.Certificate, error) {
	// Attempt to load the certificate and key from the provided files
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		// Return a formatted error message instead of logging and exiting
		return tls.Certificate{}, fmt.Errorf("failed to load TLS certificate from certFile %s and keyFile %s: %w", certFile, keyFile, err)
	}

	return cert, nil
}
