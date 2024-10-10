package handlers

import (
	"libery-dungeon-libs/libs/libery_networking"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
)

var alive_route_path string = "/alive"

var ALIVE_ROUTE *patriot_router.Route = patriot_router.NewRoute(alive_route_path, true)

func AliveHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		return
	}
}
