package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication/service_requests/categories_requests"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var medias_path string = "/medias"

var MEDIAS_ROUTE *patriot_router.Route = patriot_router.NewRoute(fmt.Sprintf("%s(/.+)?", medias_path), false)

func MediasHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getMediasHandler
		case http.MethodPost:
			request_handler_func = postMediasHandler
		case http.MethodPatch:
			request_handler_func = patchMediasHandler
		case http.MethodDelete:
			request_handler_func = deleteMediasHandler
		case http.MethodPut:
			request_handler_func = putMediasHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getMediasHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource_path {
	case fmt.Sprintf("%s/media-in-list", medias_path):
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getMediaInListHandler)
	}

	handler_func(response, request)
}

func getMediaInListHandler(response http.ResponseWriter, request *http.Request) {
	var request_body *categories_requests.GetMediaListRequest

	request_body, err := categories_requests.ParseMediaListParams(request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/medias.getMediaInListHandler: while parsing request body: \n\n%s", err))
		response.WriteHeader(400)
		return
	}

	medias, err := repository.CategoriesRepo.GetMediaIdentityList(request.Context(), request_body.MediaUUIDs)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/medias.getMediaInListHandler: while getting medias: \n\n%s", err))
		response.WriteHeader(400)
		return
	}

	var cluster_consistent_medias []dungeon_models.MediaIdentity = make([]dungeon_models.MediaIdentity, 0)
	for _, media := range medias {
		if media.ClusterUUID != request_body.ClusterUUID {
			continue
		}

		cluster_consistent_medias = append(cluster_consistent_medias, media)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(cluster_consistent_medias)
}

func postMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
