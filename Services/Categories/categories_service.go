package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_categories_service/Config"
	"libery_categories_service/database"
	"libery_categories_service/handlers"
	"libery_categories_service/middleware"
	"libery_categories_service/repository"
	"libery_categories_service/server"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/categories-tree.*", false), middleware.CheckMainCategoryProxyID(handlers.CategoriesTreeHandler(server)))
	router.RegisterRoute(patriot_router.NewRoute("/categories(/.+)?$", false), handlers.CategoriesHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/clusters", false), handlers.CategoriesClustersHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/service-fs(/.+)?$", false), handlers.ServiceFSHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/search", true), handlers.SearchHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/trashcan(/.+)?$", false), handlers.TrashcanHandler(server))
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting %s", app_config.SERVICE_NAME))

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = app_config.SERVICE_NAME

	// ------ Repositories ------

	var categories_repo *database.CategoriesMysql

	categories_repo, err := database.NewCategoriesMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	categories_clusters_repo, err := database.NewCategoriesClustersMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	trash_repo, err := database.NewTrashcanDatabase(app_config.TRASH_STORAGE_PATH)
	if err != nil {
		echo.EchoFatal(err)
	}
	defer trash_repo.Save()

	medias_repo, err := database.NewMediasMysql()
	if err != nil {
		echo.EchoFatal(err)
	}

	repository.SetCategoriesImplementation(categories_repo)
	repository.SetCategoriesClustersImplementation(categories_clusters_repo)
	repository.SetTrashImplementation(trash_repo)
	repository.SetMediasImplementation(medias_repo)

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	// ------ Server ------

	categories_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	// ------ GRPC Server ------

	categories_grpc_server, err := server.NewCategoriesServer(app_config.CATEGORIES_GRPC_SERVER_PORT)
	if err != nil {
		echo.EchoFatal(err)
	}

	categories_service.SetGrpcServer(categories_grpc_server)

	categories_service.OnAfterShutdown(app_config.ClosePlatformCommunication)

	categories_service.StartServer(BinderRoutes)
}
