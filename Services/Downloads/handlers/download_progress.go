package handlers

import (
	"libery-dungeon-libs/libs/libery_networking"
	"libery_downloads_service/workflows"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func DownloadProgressHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getDownloadProgressHandler(response, request)
		case http.MethodPost:
			postDownloadProgressHandler(response, request)
		case http.MethodPatch:
			patchDownloadProgressHandler(response, request)
		case http.MethodDelete:
			deleteDownloadProgressHandler(response, request)
		case http.MethodPut:
			putDownloadProgressHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getDownloadProgressHandler(response http.ResponseWriter, request *http.Request) {
	var download_uuid string = request.URL.Query().Get("download_uuid")
	if download_uuid == "" {
		echo.EchoWarn("While getting download progress: No download_uuid provided")
		response.WriteHeader(400)
		return
	}

	err := workflows.RegisterProgressListener(download_uuid, response, request)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(500)
		return
	}
}

func postDownloadProgressHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchDownloadProgressHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteDownloadProgressHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putDownloadProgressHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
