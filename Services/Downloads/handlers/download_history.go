package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	"libery_downloads_service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func DownloadHistoryHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getDownloadHistoryHandler(response, request)
		case http.MethodPost:
			postDownloadHistoryHandler(response, request)
		case http.MethodPatch:
			patchDownloadHistoryHandler(response, request)
		case http.MethodDelete:
			deleteDownloadHistoryHandler(response, request)
		case http.MethodPut:
			putDownloadHistoryHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getDownloadHistoryHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	if resource_path == "/download-history/download" {
		getDownloadHandler(response, request)
	} else {
		echo.Echo(echo.RedBG, fmt.Sprintf("Resource not found: %s", resource_path))
		response.WriteHeader(404)
	}

	return
}

func getDownloadHandler(response http.ResponseWriter, request *http.Request) {
	var download_uuid string = request.URL.Query().Get("download_uuid")

	if download_uuid == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing download_uuid parameter"))
		response.WriteHeader(400)
		return
	}

	download_exists, err := repository.Downloads.DownloadExists(download_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error checking download existence: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	response_body := &struct {
		Exists        bool `json:"exists"`
		DownloadCount int  `json:"download_count"`
	}{
		Exists:        false,
		DownloadCount: 0,
	}

	if download_exists {

		registered_download, err := repository.Downloads.GetDownload(download_uuid)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error getting download: %s", err.Error()))
			response.WriteHeader(500)
			return
		}

		response_body.Exists = true
		response_body.DownloadCount = len(registered_download.DownloadFiles)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	err = json.NewEncoder(response).Encode(response_body)
}

func postDownloadHistoryHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchDownloadHistoryHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteDownloadHistoryHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putDownloadHistoryHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
