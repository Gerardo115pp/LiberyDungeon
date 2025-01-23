package handlers

import (
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_medias_service/Config"
	service_helpers "libery_medias_service/helpers"
	service_workflows "libery_medias_service/workflows"
	service_common_workflows "libery_medias_service/workflows/common"
	service_workflows_errors "libery_medias_service/workflows/errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func ThumbnailsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getThumbnailsHandler(response, request)
		case http.MethodPost:
			postThumbnailsHandler(response, request)
		case http.MethodPatch:
			patchThumbnailsHandler(response, request)
		case http.MethodDelete:
			deleteThumbnailsHandler(response, request)
		case http.MethodPut:
			putThumbnailsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	echo.EchoDebug(fmt.Sprintf("resource: %s", resource))

	switch {
	case strings.HasPrefix(resource, "/thumbnails-fs/libery-trashcan"):
		getTrashcanThumbnailsHandler(response, request)
	default:
		getMediaThumbnailsHandler(response, request)
	}
}

func getMediaThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_uuid string = request.URL.Query().Get("cluster_uuid")
	if cluster_uuid == "" {
		dungeon_helpers.WriteRejection(response, 400, "Missing cluster_uuid")
		return
	}

	media_path, err := service_common_workflows.GetMediaFsPathFromRequest(request, request.URL.Path, cluster_uuid)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting media fs path from request: %s", err.Error()))
		response.WriteHeader(http.StatusNotFound)
		return
	}

	file_descriptor, err := service_helpers.GetFileDescriptor(media_path)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting media file descriptor: %s", err.Error()))
		response.WriteHeader(http.StatusNotFound)
		return
	}
	defer file_descriptor.Close()

	var str_width string = request.URL.Query().Get("width")
	var thumbnail_width int = app_config.THUMBNAIL_WIDTH

	if str_width != "" {
		requested_width, err := strconv.Atoi(str_width)
		if err != nil {
			http.Error(response, fmt.Sprintf("Invalid width parameter: %s", str_width), http.StatusBadRequest)
			return
		}

		if requested_width >= app_config.MIN_THUMBNAIL_WIDTH {
			thumbnail_width = requested_width
		}
	}

	thumbnail_response, labeled_error := service_workflows.GetFileThumbnail(file_descriptor, thumbnail_width)
	if labeled_error != nil {
		var status_code int = 500 // Internal Server Error
		if labeled_error.Label == service_workflows_errors.ErrMediaNotSupported {
			status_code = 501 // Not Implemented
		}
		labeled_error.AppendContext("In getThumbnailsHandler, in service_workflows.GetFileThumbnail call")
		fmt.Println(labeled_error)

		http.Error(response, "Failed to process thumbnail", status_code)
		return
	}

	response.Header().Set("Content-Type", thumbnail_response.MimeType)
	response.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", thumbnail_response.Filename))
	response.Header().Set("Content-Length", fmt.Sprintf("%d", thumbnail_response.MediaLength))
	response.Header().Set("Cache-Control", "private, max-age=604800") // 1 week

	echo.EchoDebug(fmt.Sprintf("thumbnail_data.Len(): %d", thumbnail_response.MediaLength))

	response.WriteHeader(http.StatusOK)
	response.Write(thumbnail_response.MediaStream.Bytes())
}

func getTrashcanThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	var media_name string = filepath.Base(request.URL.Path)
	var width_str string = request.URL.Query().Get("width")

	if media_name == "" {
		response.WriteHeader(400)
		return
	}

	var width int = app_config.THUMBNAIL_WIDTH

	if width_str != "" {
		width, _ = strconv.Atoi(width_str)

		if width < app_config.MIN_THUMBNAIL_WIDTH {
			width = app_config.MIN_THUMBNAIL_WIDTH
		}
	}

	file_path := filepath.Join(app_config.TRASH_MEDIA_PATH, media_name)

	file_descriptor, err := service_helpers.GetFileDescriptor(file_path)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting file descriptor: %s", err.Error()))
		response.WriteHeader(http.StatusNotFound)
		return
	}

	thumbnail_response, labeled_error := service_workflows.GetFileThumbnail(file_descriptor, width)
	if labeled_error != nil {
		var status_code int = 500 // Internal Server Error
		if labeled_error.Label == service_workflows_errors.ErrMediaNotSupported {
			status_code = 501 // Not Implemented
		}
		labeled_error.AppendContext("In getThumbnailsHandler, in service_workflows.GetFileThumbnail call")
		fmt.Println(labeled_error)

		http.Error(response, "Failed to process thumbnail", status_code)
		return
	}

	response.Header().Set("Content-Type", thumbnail_response.MimeType)
	response.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", thumbnail_response.Filename))
	response.Header().Set("Content-Length", fmt.Sprintf("%d", thumbnail_response.MediaLength))
	response.Header().Set("Cache-Control", "private, max-age=604800") // 1 week

	echo.EchoDebug(fmt.Sprintf("thumbnail_data.Len(): %d", thumbnail_response.MediaLength))

	response.WriteHeader(http.StatusOK)
	response.Write(thumbnail_response.MediaStream.Bytes())
}

func postThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putThumbnailsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
