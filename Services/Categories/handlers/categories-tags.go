package handlers

import (
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/repository"
	"net/http"
	"strconv"

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
	var category_uuid string = request.URL.Query().Get("category_uuid")
	var tag_id_string string = request.URL.Query().Get("tag_id")
	var err error

	if category_uuid == "" || tag_id_string == "" {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go postCategoryTagsContentHandler: category_uuid<%s> or tag_id<%s> is empty", category_uuid, tag_id_string))
		response.WriteHeader(400)
		return
	}

	var tag_id int
	tag_id, err = strconv.Atoi(tag_id_string)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go postCategoryTagsContentHandler: tag_id<%s> is not an integer", tag_id_string))
		response.WriteHeader(400)
		return
	}

	category_content, err := repository.CategoriesRepo.GetCategoryMedias(request.Context(), category_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go postCategoryTagsContentHandler: error getting category medias: %s", err))
		response.WriteHeader(404)
		return
	}

	var media_uuids []string = make([]string, 0, len(category_content))

	for _, media := range category_content {
		media_uuids = append(media_uuids, media.Uuid)
	}

	tagged, err := communication.Metadata.TagEntities(tag_id, media_uuids, dungeon_models.ENTITY_TYPE_MEDIA)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go postCategoryTagsContentHandler: error tagging entities while calling metadata service: %s", err))
		response.WriteHeader(500)
		return
	}

	dungeon_helpers.WriteBooleanResponse(response, tagged)
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
	var category_uuid string = request.URL.Query().Get("category_uuid")
	var tag_id_string string = request.URL.Query().Get("tag_id")
	var err error

	if category_uuid == "" || tag_id_string == "" {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go deleteCategoryTagsContentHandler: category_uuid<%s> or tag_id<%s> is empty", category_uuid, tag_id_string))
		response.WriteHeader(400)
		return
	}

	var tag_id int

	tag_id, err = strconv.Atoi(tag_id_string)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go deleteCategoryTagsContentHandler: tag_id<%s> is not an integer", tag_id_string))
		response.WriteHeader(400)
		return
	}

	category_content, err := repository.CategoriesRepo.GetCategoryMedias(request.Context(), category_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go deleteCategoryTagsContentHandler: error getting category medias: %s", err))
		response.WriteHeader(404)
		return
	}

	var media_uuids []string = make([]string, 0, len(category_content))

	for _, media := range category_content {
		media_uuids = append(media_uuids, media.Uuid)
	}

	untagged, err := communication.Metadata.UntagEntities(tag_id, media_uuids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Services/Categories/handlers/categories-tags.go deleteCategoryTagsContentHandler: error untagging entities while calling metadata service: %s", err))
		response.WriteHeader(500)
		return
	}

	dungeon_helpers.WriteBooleanResponse(response, untagged)
}

func putCategoryTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
