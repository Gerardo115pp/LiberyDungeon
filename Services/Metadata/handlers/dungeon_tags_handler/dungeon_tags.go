package dungeon_tags_handler

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	"libery-metadata-service/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var dungeon_tags_path string = "/dungeon-tags(/.+)?"

var DUNGEON_TAGS_ROUTE *patriot_router.Route = patriot_router.NewRoute(dungeon_tags_path, false)

func DungeonTagsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var tag_handler_prefix string = "/dungeon-tags/tags"
		var tag_taxonomy_prefix string = "/dungeon-tags/taxonomies"
		var resource string = request.URL.Path

		if strings.HasPrefix(resource, tag_handler_prefix) {
			tagsHandler(response, request)
			return
		} else if strings.HasPrefix(resource, tag_taxonomy_prefix) {
			tagTaxonomiesHandler(response, request)
			return
		}

		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getDungeonTagsHandler
		case http.MethodPost:
			request_handler_func = postDungeonTagsHandler
		case http.MethodPatch:
			request_handler_func = patchDungeonTagsHandler
		case http.MethodDelete:
			request_handler_func = deleteDungeonTagsHandler
		case http.MethodPut:
			request_handler_func = putDungeonTagsHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func postDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {

	case "/dungeon-tags/tag-entity":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTag(postDungeonTagEntityHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
}

func postDungeonTagEntityHandler(response http.ResponseWriter, request *http.Request) {
	var tag_id_str string = request.URL.Query().Get("tag_id")
	var entity_identifier string = request.URL.Query().Get("entity")

	if tag_id_str == "" || entity_identifier == "" {
		echo.Echo(echo.RedFG, "In postDungeonTagEntityHandler, tag_id or entity is empty\n")
		response.WriteHeader(400)
		return
	}

	var tag_id int
	tag_id, err := strconv.Atoi(tag_id_str)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagEntityHandler, while converting tag_id to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	tagging_id, err := repository.DungeonTagsRepo.TagEntityCTX(request.Context(), tag_id, entity_identifier)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagEntityHandler, while tagging entity: %s\n", err))
		response.WriteHeader(500)
		return
	}

	dungeon_helpers.WriteSingleIntResponseWithStatus(response, int(tagging_id), 201)
}

func patchDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/tag-entity":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTag(deleteDungeonTagEntityHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
}

func deleteDungeonTagEntityHandler(response http.ResponseWriter, request *http.Request) {
	var entity string = request.URL.Query().Get("entity")
	var tag_id_str string = request.URL.Query().Get("tag_id")

	if entity == "" || tag_id_str == "" {
		echo.Echo(echo.RedFG, "In deleteDungeonTagEntityHandler, entity or tag_id are empty\n")
		response.WriteHeader(400)
		return
	}

	var tag_id int
	tag_id, err := strconv.Atoi(tag_id_str)

	err = repository.DungeonTagsRepo.RemoveTagFromEntity(tag_id, entity)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagEntityHandler, while deleting tag: %s\n", err))
		response.WriteHeader(500)
		return
	}

	err = repository.DungeonTagsRepo.RemoveTagFromEntityCTX(request.Context(), tag_id, entity)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagHandler, while deleting tag: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func putDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
