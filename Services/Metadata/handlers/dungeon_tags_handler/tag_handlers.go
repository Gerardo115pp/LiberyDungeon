package dungeon_tags_handler

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication/service_requests/metadata_requests"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	service_models "libery-metadata-service/models"
	"libery-metadata-service/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func tagsHandler(response http.ResponseWriter, request *http.Request) {
	var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch request.Method {
	case http.MethodGet:
		request_handler_func = getTagHandler
	case http.MethodPost:
		request_handler_func = postTagHandler
	case http.MethodPatch:
		request_handler_func = patchTagHandler
	case http.MethodDelete:
		request_handler_func = deleteTagHandler
	case http.MethodPut:
		request_handler_func = putTagHandler
	case http.MethodOptions:
		request_handler_func = dungeon_helpers.AllowAllHandler
	default:
		request_handler_func = dungeon_helpers.MethodNotAllowedHandler
	}

	request_handler_func(response, request)
}

func getTagHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/tags":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonTagsTagHandler)
	case "/dungeon-tags/tags/by-ids":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonTagsByIDListHandler)
	case "/dungeon-tags/tags/entity":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getEntityTagsHandler)
	case "/dungeon-tags/tags/cluster":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonClusterTagsHandler)
	case "/dungeon-tags/tags/user-defined/cluster":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getDungeonClusterNonInternalTaxonomiesHandler)
	case "/dungeon-tags/tags/matching-entities":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getEntitiesWithTagsHandler)
	case "/dungeon-tags/tags/paginated/matching-entities":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getEntitiesWithTagsPaginatedHandler)
	}

	handler_func(response, request)
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

func getDungeonTagsByIDListHandler(response http.ResponseWriter, request *http.Request) {
	request_body, err := metadata_requests.ParseTagIDsListParams(request)
	if err != nil {
		dungeon_helpers.WriteRejection(response, 400, "Invalid query parameters")
		return
	}

	dungeon_tags, err := repository.DungeonTagsRepo.GetTagsByIDsCTX(request.Context(), request_body.TagList)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonTagsByIDListHandler, while getting tags by ids: %s\n", err))
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(dungeon_tags)
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
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(tag)
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

func getDungeonClusterNonInternalTaxonomiesHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_uuid string = request.URL.Query().Get("cluster_uuid")

	if cluster_uuid == "" {
		echo.Echo(echo.RedFG, "In getDungeonClusterNonInternalTaxonomiesHandler, cluster_uuid is empty\n")
		response.WriteHeader(400)
		return
	}

	taxonomy_tags, err := repository.DungeonTagsRepo.GetClusterTagsByInternalValueCTX(request.Context(), cluster_uuid, false)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getDungeonClusterNonInternalTaxonomiesHandler, while getting cluster taxonomies: %s\n", err))
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(taxonomy_tags)
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

	var entities_by_type map[string][]string = make(map[string][]string)

	for _, entity := range entities {
		entities_of_type, list_exists := entities_by_type[entity.EntityType]

		if !list_exists {
			entities_of_type = make([]string, 0)
		}

		entities_of_type = append(entities_of_type, entity.TaggedEntityUUID)
		entities_by_type[entity.EntityType] = entities_of_type
	}

	response.Header().Add("Content-Type", "application/json")

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(entities_by_type)
}

func getEntitiesWithTagsPaginatedHandler(response http.ResponseWriter, request *http.Request) {
	var page_num int
	var page_size int

	var query_param string = request.URL.Query().Get("page")

	if query_param == "" {
		echo.Echo(echo.RedFG, "In getEntitiesWithTagsPaginatedHandler, page query parameter is empty\n")
		response.WriteHeader(400)
		return
	}

	page_num, err := strconv.Atoi(query_param)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsPaginatedHandler, while converting page query parameter to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	query_param = request.URL.Query().Get("page_size")

	if query_param == "" {
		echo.Echo(echo.RedFG, "In getEntitiesWithTagsPaginatedHandler, page_size query parameter is empty\n")
		response.WriteHeader(400)
		return
	}

	page_size, err = strconv.Atoi(query_param)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsPaginatedHandler, while converting page_size query parameter to int: %s\n", err))
		response.WriteHeader(400)
		return
	}

	var tag_ids []int

	tag_ids, err = dungeon_helpers.ParseQueryParameterAsIntSlice(request, "tags")
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

	starting_index := page_size * (page_num - 1)
	max_index := page_size * page_num

	if len(entities) < starting_index {
		echo.Echo(echo.RedFG, "In getEntitiesWithTagsPaginatedHandler, more content than whats available\n")
		response.WriteHeader(404)
		return
	}

	var entities_by_type map[string][]string = make(map[string][]string)

	var entity_iterator int = starting_index

	for entity_iterator < max_index && entity_iterator < len(entities) {
		var entity service_models.DungeonTaggingCompact = entities[entity_iterator]
		entity_iterator++

		entities_of_type, list_exists := entities_by_type[entity.EntityType]

		if !list_exists {
			entities_of_type = make([]string, 0)
		}

		entities_of_type = append(entities_of_type, entity.TaggedEntityUUID)
		entities_by_type[entity.EntityType] = entities_of_type
	}

	var entites_count int = len(entities)

	dungeon_helpers.WritePaginatedResponse(response, entities_by_type, page_num, entites_count/page_size, entites_count)
	return
}

func postTagHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {

	case "/dungeon-tags/tags":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(postDungeonTagHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In postDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
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

func patchTagHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/tags/name":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(patchDungeonTagNameHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
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

func deleteTagHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/dungeon-tags/tags":
		handler_func = dungeon_middlewares.CheckUserCan_DungeonTagsCreate(deleteDungeonTagHandler)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In deleteDungeonTagsHandler, invalid resource: %s\n", resource))
	}

	handler_func(response, request)
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

func putTagHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
