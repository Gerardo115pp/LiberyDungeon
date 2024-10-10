package service_clients

import (
	"net/http"

	"google.golang.org/grpc/credentials"
)

type BaseServiceClient struct {
	GrpcTransport credentials.TransportCredentials
	BaseDomain    string
	HttpTransport *http.Transport
	HttpAddress   string
	GrpcAddress   string
}
