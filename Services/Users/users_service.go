package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_users_service/Config"
	"libery_users_service/databases/sqlite_users"
	"libery_users_service/handlers"
	"libery_users_service/repository"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(http_server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(handlers.ALIVE_ROUTE, handlers.AliveHandler(http_server))
	router.RegisterRoute(handlers.INITIAL_SETUP_ROUTE, handlers.IsInitialSetupHandler(http_server))
	router.RegisterRoute(handlers.USERS_ROUTE, handlers.UsersHandler(http_server))
	router.RegisterRoute(handlers.USER_AUTH_ROUTE, handlers.UserAuthHandler(http_server))
	router.RegisterRoute(handlers.USER_ROLES_ROUTE, handlers.UserRolesHandler(http_server))
}

func main() {
	echo.Echo(echo.GreenFG, "Booting up")

	app_config.VerifyConfig()
	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting %s", app_config.SERVICE_NAME))

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = app_config.SERVICE_NAME

	// ------ Repositories ------

	users_db, err := sqlite_users.NewUsersDB()
	if err != nil {
		echo.EchoFatal(err)
	}

	repository.SetUsersRepository(users_db)

	// ------ Servers ------

	users_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	users_service.OnAfterShutdown(func() {
		app_config.ClosePlatformCommunication()

		users_db.Close()
	})

	users_service.StartServer(BinderRoutes)
}
