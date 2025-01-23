package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	medias_http_requests "libery-dungeon-libs/communication/service_requests/medias_requests"
	"libery-dungeon-libs/dungeonsec/access_sec"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	"libery_medias_service/repository"
	"libery_medias_service/workflows"

	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var medias_resource_path string = "/medias"

var MEDIAS_ROUTE *patriot_router.Route = patriot_router.NewRoute(fmt.Sprintf("%s(/.+)?$", medias_resource_path), false)

func MediasHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getMediasHandler(response, request)
		case http.MethodPost:
			postMediasHandler(response, request)
		case http.MethodPatch:
			patchMediasHandler(response, request)
		case http.MethodDelete:
			deleteMediasHandler(response, request)
		case http.MethodPut:
			putMediasHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getMediasHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	handler_func := dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/medias/by-uuid":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getMediaByUUIDHandler)
	case "/medias/identity":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getMediaIdentityHandler)
	default:
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediasHandler: Resource not found: %s", resource))
	}

	handler_func(response, request)
}

func getMediaByUUIDHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("uuid")

	if media_uuid == "" {
		echo.Echo(echo.RedBG, "In MediasService.medias.getMediaByUUIDHandler: Missing media uuid query parameter")
		response.WriteHeader(400)
		return
	}

	media, err := repository.MediasRepo.GetMediaByID(request.Context(), media_uuid)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaByUUIDHandler: Error getting media by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=10600") // 3 hours

	response.WriteHeader(200)
	json.NewEncoder(response).Encode(media)
}

// Returns a MediaIdentity item which is an extension of the Media struct which includes all the necessary information from the category and cluster to request a media file.
func getMediaIdentityHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("uuid")

	media_identity, err := repository.MediasRepo.GetMediaIdentity(request.Context(), media_uuid)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Error getting media by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	has_cluster_access := access_sec.RequestHasClusterAccess(media_identity.ClusterUUID, request)
	if !has_cluster_access {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Request does not have access to cluster '%s'", media_identity.ClusterUUID))
		response.WriteHeader(403)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=10600") // 3 hours

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(media_identity)
}

// Deprecated: This endpoint is now removed.
func postMediasHandler(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.RedBG, "DEPRECATED: use POST '/upload-streams/stream-fragment' instead")
	dungeon_helpers.ResourceNotFoundHandler(response, request)
}

func patchMediasHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path
	var resource_handler http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource_path {
	case fmt.Sprintf("%s/rename", medias_resource_path):
		resource_handler = dungeon_middlewares.CheckUserCan_ContentAlter(patchRenameMediaHandler)
	}

	resource_handler(response, request)
}

func patchRenameMediaHandler(response http.ResponseWriter, request *http.Request) {
	var rename_request *medias_http_requests.RenameMediaRequest = new(medias_http_requests.RenameMediaRequest)

	err := json.NewDecoder(request.Body).Decode(rename_request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/medias.pathRenameMediaHandler: Error decoding request body: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	if rename_request.NewName == "" {
		echo.Echo(echo.RedFG, "In handlers/medias.pathRenameMediaHandler: Missing new name")
		response.WriteHeader(400)
		return
	}

	var target_media *dungeon_models.MediaIdentity

	var request_context context.Context = request.Context()

	target_media, err = repository.MediasRepo.GetMediaIdentity(request_context, rename_request.MediaUUID)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/medias.pathRenameMediaHandler: Error getting media by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	if target_media.Media.Name == rename_request.NewName {
		echo.Echo(echo.YellowFG, "In handlers/medias.pathRenameMediaHandler: New name is the same as the current name")
		response.WriteHeader(200)
		return
	}

	err = workflows.RenameMedia(request_context, *target_media, rename_request.NewName)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/medias.pathRenameMediaHandler: Error renaming media: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func deleteMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
