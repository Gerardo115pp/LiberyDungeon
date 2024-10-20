package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	service_models "libery-metadata-service/models"
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
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/global-taxonomies":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonGlobalTaxonomiesHandler)
	case "/dungeon-tags/cluster-taxonomies":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonClusterTaxonomiesHandler)
	case "/dungeon-tags/tag":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonTagsTagHandler)
	}

	handler_func(response, request)
}

func getDungeonGlobalTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	taxonomies, err := repository.DungeonTagsRepo.GetGlobalTaxonomiesCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonGlobalTaxonomiesHandler, while getting global taxonomies: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(taxonomies)
}

func getDungeonClusterTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_uuid string = request.URL.Query().Get("cluster_uuid")

	if cluster_uuid == "" {
		echo.Echo(echo.RedFG, "In getDungeonClusterTaxonomiesHandler, cluster_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	taxonomies, err := repository.DungeonTagsRepo.GetClusterTaxonomiesCTX(request.Context(), cluster_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonClusterTaxonomiesHandler, while getting cluster taxonomies: %s\n", err))
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(taxonomies)
}

func getDungeonTagsTagHandler(response http.ResponseWriter, request *http.Request) {
	var taxonomy_uuid string = request.URL.Query().Get("taxonomy")
	var tag_id string = request.URL.Query().Get("id")
	var tag_name string = request.URL.Query().Get("name")

	if tag_id != "" {
		getDungeonTagByIdHandler(response, request)
		return
	} else if tag_name != "" && taxonomy_uuid != "" {
		getDungeonTagByNameHandler(response, request)
		return
	}

	echo.Echo(echo.RedFG, "In getDungeonTagsTagHandler, invalid query parameters\n")
	dungeon_helpers.WriteRejection(response, 400, "Invalid query parameters")
}

func getDungeonTagByIdHandler(response http.ResponseWriter, request *http.Request) {
	var tag_id_str string = request.URL.Query().Get("id")
	var tag_id int
	tag_id, err := strconv.Atoi(tag_id_str)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTagByIdHandler, while converting tag_id to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	tag, err := repository.DungeonTagsRepo.GetTagByIdCTX(request.Context(), tag_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTagByIdHandler, while getting tag by id: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(tag)
}

func getDungeonTagByNameHandler(response http.ResponseWriter, request *http.Request) {
	var taxonomy_uuid string = request.URL.Query().Get("taxonomy")
	var tag_name string = request.URL.Query().Get("name")

	tag, err := repository.DungeonTagsRepo.GetTagByNameCTX(request.Context(), tag_name, taxonomy_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTagByNameHandler, while getting tag by name: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(tag)
}

func postDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomy":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(postDungeonTagTaxonomyHandler)
	case "/dungeon-tags/tag":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(postDungeonTagHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
}

func postDungeonTagTaxonomyHandler(response http.ResponseWriter, request *http.Request) {
	var new_taxonomy *service_models.TagTaxonomy = new(service_models.TagTaxonomy)

	err := json.NewDecoder(request.Body).Decode(new_taxonomy)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagTaxonomyHandler, while decoding request body: %s\n", err))
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = repository.DungeonTagsRepo.CreateTaxonomyCTX(request.Context(), new_taxonomy)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagTaxonomyHandler, while creating taxonomy: %s\n", err))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(201)
}

func postDungeonTagHandler(response http.ResponseWriter, request *http.Request) {
	var new_tag *service_models.DungeonTag = new(service_models.DungeonTag)

	err := json.NewDecoder(request.Body).Decode(new_tag)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagHandler, while decoding request body: %s\n", err))
		response.WriteHeader(400)
		return
	}

	if new_tag.Taxonomy == "" || new_tag.Name == "" {
		echo.Echo(echo.RedFG, "In postDungeonTagHandler, taxonomy or name is empty\n")
		response.WriteHeader(400)
		return
	}

	new_tag.Name = strings.ToLower(new_tag.Name)

	new_tag.RecalculateNameTaxonomy()

	err = repository.DungeonTagsRepo.CreateTagCTX(request.Context(), new_tag)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagHandler, while creating tag: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(201)

	json.NewEncoder(response).Encode(new_tag)
}

func patchDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
