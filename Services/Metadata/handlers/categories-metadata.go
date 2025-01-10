package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication/service_requests/metadata_requests"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	"libery-metadata-service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var categories_metadata_path string = "/categories-metadata"

var CATEGORIES_METADATA_ROUTE *patriot_router.Route = patriot_router.NewRoute(fmt.Sprintf("%s(/.+)?", categories_metadata_path), false)

func CategoriesMetadataHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getCategoriesMetadataHandler
		case http.MethodPost:
			request_handler_func = postCategoriesMetadataHandler
		case http.MethodPatch:
			request_handler_func = patchCategoriesMetadataHandler
		case http.MethodDelete:
			request_handler_func = deleteCategoriesMetadataHandler
		case http.MethodPut:
			request_handler_func = putCategoriesMetadataHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getCategoriesMetadataHandler(response http.ResponseWriter, request *http.Request) {
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler
	var resource_path string = request.URL.Path

	switch resource_path {
	case fmt.Sprintf("%s/category-config", categories_metadata_path):
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getCategoryConfigHandler)
	}

	handler_func(response, request)
}

func getCategoryConfigHandler(response http.ResponseWriter, request *http.Request) {
	var category_uuid string = request.URL.Query().Get("category_uuid")

	if category_uuid == "" {
		echo.Echo(echo.RedFG, "In handlers/categories-metadata.go getCategoryConfigHandler: category_uuid is empty")
		response.WriteHeader(400)
		return
	}

	category_config := repository.CategoriesConfigRepo.GetCategoryConfig(category_uuid)

	if category_config == nil {
		echo.Echo(echo.RedFG, "In handlers/categories-metadata.go getCategoryConfigHandler: category config should never be nil. Programming error.")
		response.WriteHeader(500)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store")

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(category_config)
}

func postCategoriesMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchCategoriesMetadataHandler(response http.ResponseWriter, request *http.Request) {
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler
	var resource_path string = request.URL.Path

	switch resource_path {
	case fmt.Sprintf("%s/category-config/billboard-tags", categories_metadata_path):
		handler_func = dungeon_middlewares.CheckUserCan_ContentAlter(patchCategoryMetadataBillboardTagsHandler)
	case fmt.Sprintf("%s/category-config/billboard-medias", categories_metadata_path):
		handler_func = dungeon_middlewares.CheckUserCan_ContentAlter(patchCategoryMetadataBillboardMediasHandler)
	}

	handler_func(response, request)
}

func patchCategoryMetadataBillboardTagsHandler(response http.ResponseWriter, request *http.Request) {
	var request_body *metadata_requests.PatchCategoryBillboardTagsRequest = new(metadata_requests.PatchCategoryBillboardTagsRequest)

	err := json.NewDecoder(request.Body).Decode(request_body)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/categories-metadata.go patchCategoryMetadataBillboardTagsHandler: while decoding request body: %s", err))
		response.WriteHeader(400)
		return
	}

	if request_body.CategoryUUID == "" {
		echo.Echo(echo.RedFG, "In handlers/categories-metadata.go patchCategoryMetadataBillboardTagsHandler: category uuid is empty")
		response.WriteHeader(400)
		return
	}

	category_config := repository.CategoriesConfigRepo.GetCategoryConfig(request_body.CategoryUUID)

	category_config.BillboardDungeonTags = request_body.BillboardDungeonTags

	err = repository.CategoriesConfigRepo.UpdateCategoryConfig(category_config)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/categories-metadata.go patchCategoryMetadataBillboardTagsHandler: while updating category config: %s", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func patchCategoryMetadataBillboardMediasHandler(response http.ResponseWriter, request *http.Request) {
	var request_body *metadata_requests.PatchCategoryBillboardMediasRequest = new(metadata_requests.PatchCategoryBillboardMediasRequest)

	err := json.NewDecoder(request.Body).Decode(request_body)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/categories-metadata.go patchCategoryMetadataBillboardMediasHandler: while decoding request body: %s", err))
		response.WriteHeader(400)
		return
	}

	if request_body.CategoryUUID == "" {
		echo.Echo(echo.RedFG, "In handlers/categories-metadata.go patchCategoryMetadataBillboardMediasHandler: category uuid is empty")
		response.WriteHeader(400)
		return
	}

	category_config := repository.CategoriesConfigRepo.GetCategoryConfig(request_body.CategoryUUID)

	category_config.BillboardMediaUUIDs = request_body.BillboardMediaUUIDs

	err = repository.CategoriesConfigRepo.UpdateCategoryConfig(category_config)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/categories-metadata.go patchCategoryMetadataBillboardMediasHandler: while updating category config: %s", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func deleteCategoriesMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putCategoriesMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
