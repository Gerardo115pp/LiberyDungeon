package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	service_models "libery-metadata-service/models"
	"libery-metadata-service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func WatchPointsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getWatchPointsHandler(response, request)
		case http.MethodPost:
			postWatchPointsHandler(response, request)
		case http.MethodPatch:
			patchWatchPointsHandler(response, request)
		case http.MethodDelete:
			deleteWatchPointsHandler(response, request)
		case http.MethodPut:
			putWatchPointsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getWatchPointsHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("media_uuid")
	var err error

	if media_uuid == "" {
		echo.Echo(echo.RedFG, "media_uuid is required")
		response.WriteHeader(400)
		return
	}

	watch_point, lerr := repository.WatchPointRepo.GetWatchPointByMediaID(request.Context(), media_uuid)
	if lerr != nil {
		echo.EchoErr(lerr)
		status_code := 500

		if lerr.Label == service_models.ErrWatchPointNotFound {
			status_code = 404
		}

		response.WriteHeader(status_code)
		return
	}

	switch request.Header.Get("Accept") {
	case "application/json":
		response.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(response).Encode(watch_point)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error encoding response: %s", err.Error()))
		}
	default:
		response.Header().Set("Content-Type", "application/octet-stream")
		err = binary.Write(response, binary.LittleEndian, watch_point.StartTime)
		if err != nil {
			echo.EchoErr(err)
		}
	}
}

func postWatchPointsHandler(response http.ResponseWriter, request *http.Request) {
	watch_point_post_request := &struct {
		Media_uuid string `json:"media_uuid"`
		Start_time uint32 `json:"start_time"`
	}{}

	err := json.NewDecoder(request.Body).Decode(watch_point_post_request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error decoding request: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	lerr := repository.WatchPointRepo.InsertWatchPoint(request.Context(), watch_point_post_request.Media_uuid, watch_point_post_request.Start_time)
	if lerr != nil {
		echo.EchoErr(lerr)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(201)
}

func patchWatchPointsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteWatchPointsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putWatchPointsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
