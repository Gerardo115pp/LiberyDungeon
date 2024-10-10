package handlers

import (
	"libery-dungeon-libs/libs/libery_networking"
	"libery_downloads_service/workflows/jobs"
	"net/http"
)

func DownloadsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getDownloadsHandler(response, request)
		case http.MethodPost:
			postDownloadsHandler(response, request)
		case http.MethodPatch:
			patchDownloadsHandler(response, request)
		case http.MethodDelete:
			deleteDownloadsHandler(response, request)
		case http.MethodPut:
			putDownloadsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getDownloadsHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	if resource_path == "/downloads/current-download" {
		getCurrentDownloadHandler(response, request)
		return
	} else {
		response.WriteHeader(404)
		return
	}
}

func getCurrentDownloadHandler(response http.ResponseWriter, request *http.Request) {
	var download_uuid string = jobs.DownloadWorker.GetCurrentDownloadUUID()

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	response.Write([]byte(`{"download_uuid": "` + download_uuid + `"}`))

	return
}

func postDownloadsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchDownloadsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteDownloadsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putDownloadsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
