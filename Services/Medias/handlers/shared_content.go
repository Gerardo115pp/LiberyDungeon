package handlers

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	common_flows "libery_medias_service/workflows/common"
	"net/http"

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
	case fmt.Sprintf("%s/shared-media", shared_content_path):
		request_handler_func = getSharedMediaContentHandler
	}

	request_handler_func(response, request)
}

func getSharedMediaContentHandler(response http.ResponseWriter, request *http.Request) {
	var share_token_string string = request.URL.Query().Get("share_token")

	if share_token_string == "" {
		echo.Echo(echo.RedFG, "In handlers/shared_content.getSharedMediaContentHandler: 'share_token' was empty.")
		dungeon_helpers.WriteRejection(response, 400, "missing share_token")
		return
	}

	media_claims, err := dungeon_models.ParseMediaShareToken(share_token_string, dungeon_secrets.GetDungeonJwtSecret())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/shared_content.getSharedMediaContentHandler: Error parsing share token: %s", err.Error()))
		dungeon_helpers.WriteRejection(response, 400, "invalid share_token")
		return
	}

	media_identity := media_claims.MediaIdentity

	common_flows.WriteMediaFileResponse(*media_identity, response, request)
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
