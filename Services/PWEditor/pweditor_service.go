package main

import (
	"context"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery-pw-editor-service/Config"
	"libery-pw-editor-service/handlers"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

func BinderRoutes(server libery_networking.Server, router *patriot_router.Router) {
	router.RegisterRoute(patriot_router.NewRoute("/alive", true), handlers.AliveHandler(server))
}

func main() {
	app_config.VerifyConfig()

	echo.Echo(echo.GreenFG, "Starting %s", app_config.SERVICE_NAME)

	var new_server_config *libery_networking.ServerConfig = new(libery_networking.ServerConfig)
	new_server_config.Port = app_config.SERVICE_PORT
	new_server_config.ServiceName = string(app_config.SERVICE_NAME)

	echo.EchoDebug(fmt.Sprintf("server config: %+v", new_server_config))

	pw_editor_service, err := libery_networking.NewBroker(context.Background(), new_server_config)
	if err != nil {
		echo.EchoFatal(err)
	}

	pw_editor_service.OnAfterShutdown(app_config.ClosePlatformCommunication)

	pw_editor_service.StartServer(BinderRoutes)
}
