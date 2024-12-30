package handlers

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	"libery-dungeon-libs/libs/platform_claims"
	"libery-dungeon-libs/libs/platform_claims/resource_access"
	app_config "libery_categories_service/Config"
	"libery_categories_service/repository"
	"net/http"
	"time"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var shared_content_path string = "/shared-content"

var SHARED_CONTENT_ROUTE *patriot_router.Route = patriot_router.NewRoute(fmt.Sprintf("%s(/.+)?", shared_content_path), false)

func SharedContentHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getSharedContentHandler
		case http.MethodPost:
			request_handler_func = postSharedContentHandler
		case http.MethodPatch:
			request_handler_func = patchSharedContentHandler
		case http.MethodDelete:
			request_handler_func = deleteSharedContentHandler
		case http.MethodPut:
			request_handler_func = putSharedContentHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getSharedContentHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path = request.URL.Path
	var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource_path {
	case fmt.Sprintf("%s/media-share-token", shared_content_path):
		request_handler_func = dungeon_middlewares.CheckUserCan_ShareContent(getShareMediaTokenHandler)
	}

	request_handler_func(response, request)
}

func getShareMediaTokenHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("media_uuid")

	if media_uuid == "" {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/shared_content.getShareMediaTokenHandler: media_uuid query parameter was empty."))
		dungeon_helpers.WriteRejection(response, 400, "missing 'media_uuid'")
		return
	}

	media_identity, err := repository.CategoriesRepo.GetMediaIdentity(request.Context(), media_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/shared_content.getShareMediaTokenHandler: while getting the media identity for '%s': %s", media_uuid, err.Error()))
		dungeon_helpers.WriteRejection(response, 404, "media not found")
		return
	}

	has_access_to_media_cluster, err := resource_access.UserClaimsCookieHasClusterAccess(media_identity.ClusterUUID, request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/shared_content.getShareMediaTokenHandler: while checking if user has access to media cluster: %s", err.Error()))
		dungeon_helpers.WriteRejection(response, 401, "unauthorized")
		return
	}

	if !has_access_to_media_cluster {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/shared_content.getShareMediaTokenHandler: user does not have access to media cluster"))
		dungeon_helpers.WriteRejection(response, 401, "unauthorized")
		return
	}

	platform_claims.WriteMediaShareTokenResponse(response, media_identity, time.Now().Add(time.Second*time.Duration(app_config.SHARED_MEDIA_EXPIRATION_SECS)))
}

func postSharedContentHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchSharedContentHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteSharedContentHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putSharedContentHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
