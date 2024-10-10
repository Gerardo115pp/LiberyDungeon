package handlers

import (
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	"libery_users_service/workflows/common_workflows"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var initial_setup_route_path string = "/is-initial-setup"

var INITIAL_SETUP_ROUTE *patriot_router.Route = patriot_router.NewRoute(initial_setup_route_path, true)

func IsInitialSetupHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			response.WriteHeader(405)
			return
		}

		var is_initial_setup_done bool

		is_initial_setup_done, err := common_workflows.InitialSetupDone(request.Context())
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In Handles/IsInitialSetupHandler, while checking if initial setup is done, error: %s", err))
			response.WriteHeader(500)
			return
		}

		dungeon_helpers.WriteBooleanResponse(response, !is_initial_setup_done)
	}
}
