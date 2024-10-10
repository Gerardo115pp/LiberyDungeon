package handlers

import (
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	"libery_JD_service/workflows/jobs"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func PlatformEventsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getPlatformEventsHandler(response, request)
		case http.MethodPost:
			postPlatformEventsHandler(response, request)
		case http.MethodPatch:
			patchPlatformEventsHandler(response, request)
		case http.MethodDelete:
			deletePlatformEventsHandler(response, request)
		case http.MethodPut:
			putPlatformEventsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getPlatformEventsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/platform-events/public/suscribe":
		getPlatformEventsPublicSubscriptionHandler(response, request)
	default:
		response.WriteHeader(404)
	}
}

func getPlatformEventsPublicSubscriptionHandler(response http.ResponseWriter, request *http.Request) {
	upgrader := jobs.TownCrier.GetUpgrader()

	connection, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getPlatformEventsPublicSubscriptionHandler while upgrading connection: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	jobs.TownCrier.RegisterListener(connection)
	return
}

func postPlatformEventsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchPlatformEventsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deletePlatformEventsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putPlatformEventsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
