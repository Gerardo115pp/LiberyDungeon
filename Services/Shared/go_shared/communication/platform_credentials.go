package communication

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"

	"google.golang.org/grpc/credentials"
)

var grpc_transport credentials.TransportCredentials
var http_transport *http.Transport

func getCertificateAuthorityData(path string) ([]byte, error) {
	ca_stat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if ca_stat.Mode()&os.ModeSymlink != 0 {
		path, err = os.Readlink(path)
	}

	ca_file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer ca_file.Close()

	certificate_authority_data, err := io.ReadAll(ca_file)
	if err != nil {
		return nil, err
	}

	return certificate_authority_data, nil
}

func setGrpcTransport(path string) (err error) {
	certificate_authority_data, err := getCertificateAuthorityData(path)

	cert_pool := x509.NewCertPool()

	if ok := cert_pool.AppendCertsFromPEM(certificate_authority_data); !ok {
		return err
	}

	new_credentials_config := &tls.Config{
		RootCAs: cert_pool,
	}

	grpc_transport = credentials.NewTLS(new_credentials_config)

	return nil
}

func setHttpTransport(path string) error {
	certificate_authority_data, err := getCertificateAuthorityData(path)
	if err != nil {
		return err
	}

	authority_pool := x509.NewCertPool()
	authority_pool.AppendCertsFromPEM(certificate_authority_data)

	transport_credentials := &tls.Config{
		RootCAs: authority_pool,
	}

	http_transport = &http.Transport{TLSClientConfig: transport_credentials}

	return nil
}
