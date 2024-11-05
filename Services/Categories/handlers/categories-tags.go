package handlers

import (
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func categoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch request.Method {
	case http.MethodGet:
		request_handler_func = getCategoryTagsHandler
	case http.MethodPost:
		request_handler_func = postCategoryTagsHandler
	case http.MethodPatch:
		request_handler_func = patchCategoryTagsHandler
	case http.MethodDelete:
		request_handler_func = deleteCategoryTagsHandler
	case http.MethodPut:
		request_handler_func = putCategoryTagsHandler
	case http.MethodOptions:
		request_handler_func = dungeon_helpers.AllowAllHandler
	default:
		request_handler_func = dungeon_helpers.MethodNotAllowedHandler
	}

	request_handler_func(response, request)
}

func getCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func postCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/categories/tags/content":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTag(postCategoryTagsContentHandler)
	}

	handler_func(response, request)
}

func postCategoryTagsContentHandler(response http.ResponseWriter, request *http.Request) {
	// TODO: Implement postCategoryTagsContentHandler

	echo.Echo(echo.GreenFG, "Tagging category content")

	dungeon_helpers.WriteBooleanResponse(response, true)
}

func patchCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/categories/tags/content":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsUntag(deleteCategoryTagsContentHandler)
	default:
		echo.Echo(echo.RedFG, "Invalid resource: "+resource)
	}

	handler_func(response, request)
}

func deleteCategoryTagsContentHandler(response http.ResponseWriter, request *http.Request) {
	// TODO: Implement deleteCategoryTagsContentHandler

	echo.Echo(echo.GreenFG, "Untagging category content")

	dungeon_helpers.WriteBooleanResponse(response, true)
}

func putCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
