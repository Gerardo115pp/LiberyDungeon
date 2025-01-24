package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_medias_service/Config"
	"libery_medias_service/database"
	"libery_medias_service/handlers"
	"libery_medias_service/middleware"
	"libery_medias_service/repository"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/medias-fs.*", false), handlers.MediasFSHandler(server))
	router.RegisterRoute(handlers.MEDIAS_ROUTE, handlers.MediasHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/random-medias-fs.*", false), handlers.RandomMediasFsHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/thumbnails-fs.*", false), middleware.ParseResourcePath(handlers.ThumbnailsHandler(server)))
	router.RegisterRoute(patriot_router.NewRoute("/upload-streams.*", false), handlers.UploadStreamsHandler(server))
	router.RegisterRoute(handlers.SHARED_CONTENT_ROUTE, handlers.SharedContentHandler(server))
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, "Starting medias_service")

	// ----------------- Service Configuration -----------------

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = string(app_config.SERVICE_NAME)

	// ----------------- Repositories -----------------

	var medias_repo *database.MediasMysql
	var categories_repo *database.CategoriesMysql

	medias_repo, err := database.NewMediasMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	categories_repo, err = database.NewCategoriesMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	clusters_repo, err := database.NewCategoriesClustersMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	repository.SetMediasImplementation(medias_repo)
	repository.SetCategoriesImplementation(categories_repo)
	repository.SetCategoriesClustersImplementation(clusters_repo)

	// ----------------- Services -----------------

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	media_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	media_service.OnAfterShutdown(app_config.ClosePlatformCommunication)

	media_service.StartServer(BinderRoutes)
}
