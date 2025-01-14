package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery-metadata-service/Config"
	"libery-metadata-service/database/categories_metadata"
	cluster_metadata_database "libery-metadata-service/database/clusters_metadata"
	dungeon_tags_database "libery-metadata-service/database/dungeon_tags"
	"libery-metadata-service/database/video_moments"
	watch_point_database "libery-metadata-service/database/watch_points"
	"libery-metadata-service/handlers"
	"libery-metadata-service/handlers/dungeon_tags_handler"
	"libery-metadata-service/repository"
	"libery-metadata-service/server"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
	router.RegisterRoute(patriot_router.NewRoute("/watch-points", true), handlers.WatchPointsHandler(server))
	router.RegisterRoute(handlers.CLUSTER_METADATA_ROUTE, handlers.ClusterMetadataHandler(server))
	router.RegisterRoute(dungeon_tags_handler.DUNGEON_TAGS_ROUTE, dungeon_tags_handler.DungeonTagsHandler(server)) // This is the new standard for route registration and ownership. TODO: Adapt all routes(in all services) to this new standard
	router.RegisterRoute(handlers.CATEGORIES_METADATA_ROUTE, handlers.CategoriesMetadataHandler(server))
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, "Starting libery-metadata-service")

	// ----------------- Repositories -----------------

	var watch_point_impl repository.WatchPointRepository
	watch_point_impl = watch_point_database.NewWatchPointDatabase()

	repository.SetWatchPointRepository(watch_point_impl)

	var video_moments_impl repository.VideoMomentsRepository
	video_moments_impl = video_moments.NewVideoMomentsDB()

	repository.SetVideoMomentsRepo(video_moments_impl)

	var cluster_metadata_impl repository.ClusterMetadataRepository
	cluster_metadata_impl = cluster_metadata_database.NewClusterMetadataDB()

	repository.SetClusterMetadataRepository(cluster_metadata_impl)

	var categories_metadata_impl repository.CategoriesConfigurationRepository
	categories_metadata_impl = categories_metadata.NewCategoryConfigDB(app_config.OPERATION_DATA_PATH)

	repository.SetCategoriesConfigRepository(categories_metadata_impl)

	var dungeon_tags_impl repository.DungeonTagsRepository
	dungeon_tags_impl = dungeon_tags_database.NewDungeonTagsDB()

	repository.SetDungeonTagsRepository(dungeon_tags_impl)

	// ----------------- Services -----------------

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = app_config.SERVICE_NAME

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	metadata_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	grpc_server, err := server.NewMetadataGrpcServer(app_config.METADATA_GRPC_SERVER_PORT)
	if err != nil {
		echo.EchoFatal(err)
	}

	// RUN

	metadata_service.SetGrpcServer(grpc_server)

	metadata_service.OnAfterShutdown(app_config.ClosePlatformCommunication)

	metadata_service.StartServer(BinderRoutes)
}
