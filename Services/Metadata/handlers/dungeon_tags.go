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
	case "/dungeon-tags/taxonomy/tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonTaxonomyTags)
	case "/dungeon-tags/entity-tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getEntityTagsHandler)
	case "/dungeon-tags/cluster-tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonClusterTagsHandler)
	case "/dungeon-tags/entities-with-tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getEntitiesWithTagsHandler)
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

func getEntityTagsHandler(response http.ResponseWriter, request *http.Request) {
	var entity_uuid string = request.URL.Query().Get("entity")
	var cluster_domain string = request.URL.Query().Get("cluster_domain")

	if entity_uuid == "" {
		echo.Echo(echo.RedFG, "In getEntityTagsHandler, entity_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	tags, err := repository.DungeonTagsRepo.GetEntityTaggingsCTX(request.Context(), entity_uuid, cluster_domain)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntityTagsHandler, while getting entity tags: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(tags)
}

func getDungeonClusterTagsHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_uuid string = request.URL.Query().Get("cluster_uuid")

	if cluster_uuid == "" {
		echo.Echo(echo.RedFG, "In getDungeonClusterTagsHandler, cluster_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	tags, err := repository.DungeonTagsRepo.GetClusterTagsCTX(request.Context(), cluster_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonClusterTagsHandler, while getting cluster tags: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(tags)
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

func postDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomy":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(postDungeonTagTaxonomyHandler)
	case "/dungeon-tags/tag":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(postDungeonTagHandler)
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
		echo.Echo(echo.RedFG, "In postDungeonTagEntityHandler, tag_id or entity_uuid is empty\n")
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
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomy/name":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(patchDungeonTagTaxonomyNameHandler)
	case "/dungeon-tags/tag/name":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(patchDungeonTagNameHandler)
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

func patchDungeonTagNameHandler(response http.ResponseWriter, request *http.Request) {
	var tag_id_str string = request.URL.Query().Get("id")
	var new_name string = request.URL.Query().Get("new_name")

	if tag_id_str == "" || new_name == "" {
		echo.Echo(echo.RedFG, "In patchDungeonTagNameHandler, tag_id or new_name is empty\n")
		response.WriteHeader(400)
		return
	}

	var tag_id int
	tag_id, err := strconv.Atoi(tag_id_str)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchDungeonTagNameHandler, while converting tag_id to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	err = repository.DungeonTagsRepo.UpdateTagNameCTX(request.Context(), tag_id, new_name)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchDungeonTagNameHandler, while updating tag name: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func deleteDungeonTagsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/taxonomy":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTaxonomyCreate(deleteDungeonTagTaxonomyHandler)
	case "/dungeon-tags/tag":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(deleteDungeonTagHandler)
	case "/dungeon-tags/tag-entity":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsTag(deleteDungeonTagEntityHandler)
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

func deleteDungeonTagHandler(response http.ResponseWriter, request *http.Request) {
	var tag_id_str string = request.URL.Query().Get("id")

	if tag_id_str == "" {
		echo.Echo(echo.RedFG, "In deleteDungeonTagHandler, tag_id is empty\n")
		response.WriteHeader(400)
		return
	}

	var tag_id int

	tag_id, err := strconv.Atoi(tag_id_str)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagHandler, while converting tag_id to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	err = repository.DungeonTagsRepo.DeleteTagCTX(request.Context(), tag_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagHandler, while deleting tag: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
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

func getEntitiesWithTagsHandler(response http.ResponseWriter, request *http.Request) {
	var tag_ids []int

	tag_ids, err := dungeon_helpers.ParseQueryParameterAsIntSlice(request, "tags")
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsHandler, while parsing tags query parameter: %s\n", err))
		response.WriteHeader(400)
	}

	entities, err := repository.DungeonTagsRepo.GetEntitiesWithTaggingsCTX(request.Context(), tag_ids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsHandler, while getting entities with taggings: %s\n", err))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(entities)
}
