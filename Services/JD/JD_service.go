package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_JD_service/Config"
	"libery_JD_service/databases"
	"libery_JD_service/handlers"
	"libery_JD_service/repository"
	"libery_JD_service/server"
	"libery_JD_service/workflows/jobs"
	"libery_JD_service/workflows/workers"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/platform-events(/.+)?$", false), handlers.PlatformEventsHandler(server))
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, "Starting JD_service")

	// ------ REPOSITORIES ------

	platform_services_datastore := databases.NewPlatformServicesDatastore()
	repository.SetPlatformFeaturesRepo(platform_services_datastore)

	platform_events_store := databases.NewPlatformEventsStore()
	repository.SetPlatformEventsRepo(platform_events_store)

	// ------ WORKERS ------

	town_crier := workers.NewEventCrier()
	jobs.SetTownCrierImplementation(town_crier)
	defer town_crier.Close()

	// ------ HTTP SERVER ------

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = app_config.SERVICE_NAME

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	JD_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	// ------ GRPC SERVER ------

	grpc_server, err := server.NewJDGrpcServer(app_config.JD_GRPC_SERVER_PORT)
	if err != nil {
		echo.EchoFatal(err)
	}

	// ------ RUN ------

	JD_service.SetGrpcServer(grpc_server)

	JD_service.StartServer(BinderRoutes)
}
