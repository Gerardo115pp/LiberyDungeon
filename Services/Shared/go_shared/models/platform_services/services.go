package platform_services

import "github.com/Gerardo115pp/patriots_lib/echo"

type PlatformServiceName string

const (
	JD_SERVICE         PlatformServiceName = "libery_jd_service"
	CATEGORIES_SERVICE PlatformServiceName = "libery_categories_service"
	METADATA_SERVICE   PlatformServiceName = "libery_metadata_service"
	MEDIAS_SERVICE     PlatformServiceName = "libery_medias_service"
	PWEDITOR_SERVICE   PlatformServiceName = "libery-pw-editor-service"
	COLLECT_SERVICE    PlatformServiceName = "libery_collect_service"
	DOWNLOADS_SERVICE  PlatformServiceName = "libery_downloads_service"
	USERS_SERVICE      PlatformServiceName = "libery_users_service"
)

var ExistingServices = [...]PlatformServiceName{
	JD_SERVICE,
	CATEGORIES_SERVICE,
	METADATA_SERVICE,
	MEDIAS_SERVICE,
	PWEDITOR_SERVICE,
	COLLECT_SERVICE,
	DOWNLOADS_SERVICE,
	USERS_SERVICE,
}

type GrpcServer interface {
	Connect() error
}

func SetGrpcServers(server GrpcServer) {
	echo.Echo(echo.PinkBG, "Starting GRPC Servers")
	go server.Connect()
}

func ServiceNameValid(service_name string) bool {
	var service_platform_name PlatformServiceName = PlatformServiceName(service_name)
	var name_valid bool = false

	for _, service := range ExistingServices {
		if service == service_platform_name {
			name_valid = true
			break
		}
	}

	return name_valid
}
