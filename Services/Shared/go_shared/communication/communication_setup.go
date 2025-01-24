package communication

import (
	"libery-dungeon-libs/communication/service_clients"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"os"
)

var (
	CATEGORIES_SERVER string = os.Getenv("CATEGORIES_SERVER")
	MEDIAS_SERVER     string = os.Getenv("MEDIAS_SERVER")
	COLLECT_SERVER    string = os.Getenv("COLLECT_SERVER")
	DOWNLOADS_SERVER  string = os.Getenv("DOWNLOADS_SERVER")
	PWEDITOR_SERVER   string = os.Getenv("PWEDITOR_SERVER")
	JD_SERVER         string = os.Getenv("JD_SERVER")
	METADATA_SERVER   string = os.Getenv("METADATA_SERVER")
	GRPC_SERVER       string = os.Getenv("GRPC_SERVER")
	BASE_DOMAIN       string = os.Getenv("BASE_DOMAIN")
	SSL_CA_PATH       string = os.Getenv("SSL_CA_PATH")
	JWT_SECRET        string = os.Getenv("JWT_SECRET")
	DOMAIN_SECRET     string = os.Getenv("DOMAIN_SECRET")
)

// type ServiceCommunicationSetupParams struct {
// 	CaPath string `json:"ca_path"`
// }

// Verifies that all the endpoints are set and valid to the possible extent
// if it finds an error it will panic
func verifyEndpointConfig() {
	if JD_SERVER == "" {
		panic("JD_SERVER environment variable is required")
	}

	if CATEGORIES_SERVER == "" {
		panic("CATEGORIES_SERVER environment variable is required")
	}

	if MEDIAS_SERVER == "" {
		panic("MEDIAS_SERVER environment variable is required")
	}

	if PWEDITOR_SERVER == "" {
		panic("PWEDITOR_SERVER environment variable is required")
	}

	if COLLECT_SERVER == "" {
		panic("COLLECT_SERVER environment variable is required")
	}

	if DOWNLOADS_SERVER == "" {
		panic("DOWNLOADS_SERVER environment variable is required")
	}

	if METADATA_SERVER == "" {
		panic("METADATA_SERVER environment variable is required")
	}

	if GRPC_SERVER == "" {
		panic("GRPC_SERVER environment variable is required")
	}

	if BASE_DOMAIN == "" {
		panic("BASE_DOMAIN environment variable is required")
	}

	if SSL_CA_PATH == "" {
		panic("SSL_CA_PATH environment variable is required")
	}

	if JWT_SECRET == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if DOMAIN_SECRET == "" {
		panic("DOMAIN_SECRET environment variable is required")
	}
}

// Setup communication.
// If there is a configuration error, it will panic
func Setup() (err error) {
	verifyEndpointConfig() // if not valid, it will panic

	err = setGrpcTransport(SSL_CA_PATH)
	if err != nil {
		return
	}

	err = setHttpTransport(SSL_CA_PATH)
	if err != nil {
		return
	}

	var base_service_data service_clients.BaseServiceClient = service_clients.BaseServiceClient{
		GrpcTransport: grpc_transport,
		HttpTransport: http_transport,
		HttpAddress:   "",
		GrpcAddress:   GRPC_SERVER,
		BaseDomain:    BASE_DOMAIN,
	}

	// JD Communication

	base_service_data.HttpAddress = JD_SERVER

	JD = &service_clients.JD_Client{
		BaseServiceClient: base_service_data,
	}

	// Metadata Communication

	base_service_data.HttpAddress = METADATA_SERVER

	Metadata = &service_clients.MetadataServiceClient{
		BaseServiceClient: base_service_data,
	}

	// Medias Communication

	base_service_data.HttpAddress = MEDIAS_SERVER

	Medias = &service_clients.MediaServiceClient{
		BaseServiceClient: base_service_data,
	}

	// ----------------- Dungeonsec secrets -----------------

	dungeon_secrets.SetDungeonJwtSecret(JWT_SECRET)
	dungeon_secrets.SetDungeonDomainSecret(DOMAIN_SECRET)

	return
}

var JD *service_clients.JD_Client
var Metadata *service_clients.MetadataServiceClient
var Medias *service_clients.MediaServiceClient
