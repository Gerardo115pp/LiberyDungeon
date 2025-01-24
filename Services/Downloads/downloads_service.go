package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_downloads_service/Config"
	"libery_downloads_service/database"
	"libery_downloads_service/handlers"
	"libery_downloads_service/repository"
	"libery_downloads_service/server"
	"libery_downloads_service/workflows/jobs"
	"libery_downloads_service/workflows/workers"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/download-history(/.+)?$", false), handlers.DownloadHistoryHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/ws/download-progress", true), handlers.DownloadProgressHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/downloads(/.+)?$", false), handlers.DownloadsHandler(server))
}

func SetGrpcServers(server libery_networking.GrpcServer) {
	echo.Echo(echo.PinkBG, "Starting GRPC Servers")
	go server.Connect()
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, "Starting downloads_service")

	// ------ REPOSITORIES ------

	categories_service_repo := database.NewCategoriesServiceConn()
	repository.SetCategoriesImplementation(categories_service_repo)

	downloads_repo, err := database.NewDownloadDB()
	if err != nil {
		echo.EchoFatal(err)
	}

	repository.SetDownloadsRepository(downloads_repo)

	// ------ WORKERS ------

	downloads_worker := workers.NewAsyncDownloader()
	jobs.SetDownloadWorkerImplementation(downloads_worker)

	// ------ HTTP SERVER ------

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = string(app_config.SERVICE_NAME)

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	microservice, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	// ------ GRPC SERVERS ------

	downloads_grpc_server, err := server.NewDownloadsServer(app_config.DOWNLOADS_GRPC_SERVER_PORT)
	if err != nil {
		echo.EchoFatal(err)
	}

	microservice.SetGrpcServer(downloads_grpc_server)

	microservice.OnAfterShutdown(app_config.ClosePlatformCommunication)

	microservice.StartServer(BinderRoutes)
}
