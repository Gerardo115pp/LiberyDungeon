package dungeon_tags_handler

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	service_models "libery-metadata-service/models"
	"libery-metadata-service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func tagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch request.Method {
	case http.MethodGet:
		request_handler_func = getTagTaxonomiesHandler
	case http.MethodPost:
		request_handler_func = postTagTaxonomiesHandler
	case http.MethodPatch:
		request_handler_func = patchTagTaxonomiesHandler
	case http.MethodDelete:
		request_handler_func = deleteTagTaxonomiesHandler
	case http.MethodPut:
		request_handler_func = putTagTaxonomiesHandler
	case http.MethodOptions:
		request_handler_func = dungeon_helpers.AllowAllHandler
	default:
		request_handler_func = dungeon_helpers.MethodNotAllowedHandler
	}

	request_handler_func(response, request)
}

func getTagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomies/global":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonGlobalTaxonomiesHandler)
	case "/dungeon-tags/taxonomies/cluster":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonClusterTaxonomiesHandler)
	case "/dungeon-tags/taxonomies/tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonTaxonomyTags)
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

func getDungeonTaxonomyTags(response http.ResponseWriter, request *http.Request) {
	var taxonomy_uuid string = request.URL.Query().Get("taxonomy")

	if taxonomy_uuid == "" {
		echo.Echo(echo.RedFG, "In getDungeonTaxonomyTags, taxonomy_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	tag_taxonomy, err := repository.DungeonTagsRepo.GetTagTaxonomyCTX(request.Context(), taxonomy_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTaxonomyTags, while getting taxonomy: %s\n", err))
		response.WriteHeader(404)
		return
	}

	tags, err := repository.DungeonTagsRepo.GetTaxonomyTagsCTX(request.Context(), taxonomy_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTaxonomyTags, while getting taxonomy tags: %s\n", err))
		response.WriteHeader(500)
		return
	}

	taxonomy_tags := service_models.TaxonomyTags{
		Taxonomy: &tag_taxonomy,
		Tags:     tags,
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(taxonomy_tags)
}

func postTagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomies":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(postDungeonTagTaxonomyHandler)
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

func patchTagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomies/name":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(patchDungeonTagTaxonomyNameHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
}

func patchDungeonTagTaxonomyNameHandler(response http.ResponseWriter, request *http.Request) {
	var taxonomy_uuid string = request.URL.Query().Get("uuid")
	var new_name string = request.URL.Query().Get("new_name")

	if taxonomy_uuid == "" || new_name == "" {
		echo.Echo(echo.RedFG, "In patchDungeonTagTaxonomyNameHandler, taxonomy_uuid or new_name is empty\n")
		response.WriteHeader(400)
		return
	}

	err := repository.DungeonTagsRepo.UpdateTaxonomyNameCTX(request.Context(), taxonomy_uuid, new_name)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchDungeonTagTaxonomyNameHandler, while updating taxonomy name: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func deleteTagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomies":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(deleteDungeonTagTaxonomyHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
}

func deleteDungeonTagTaxonomyHandler(response http.ResponseWriter, request *http.Request) {
	var taxonomy_uuid string = request.URL.Query().Get("uuid")

	if taxonomy_uuid == "" {
		echo.Echo(echo.RedFG, "In deleteDungeonTagTaxonomyHandler, taxonomy_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	err := repository.DungeonTagsRepo.DeleteTaxonomyCTX(request.Context(), taxonomy_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagTaxonomyHandler, while deleting taxonomy: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func putTagTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
