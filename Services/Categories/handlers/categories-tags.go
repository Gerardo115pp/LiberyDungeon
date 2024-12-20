package handlers

import (
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/repository"
	"math"
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
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler
	echo.Echo(echo.GreenFG, "In getCategoryTagsHandler\n")

	switch resource {
	case "/categories/tags/content-tagged":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getTaggedCategoryContent)
	}

	handler_func(response, request)
}

// Returns a paginated response with the medias tagged with the provided tag list.
func getTaggedCategoryContent(response http.ResponseWriter, request *http.Request) {
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

	var medias []dungeon_models.MediaIdentity = make([]dungeon_models.MediaIdentity, 0)

	tagged_content, err := communication.Metadata.GetEntitiesWithTaggings(tag_ids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsPaginatedHandler, while getting entities with taggings: %s\n", err))
		response.WriteHeader(500)
		return
	}

	tagged_medias_uuids := tagged_content[dungeon_models.ENTITY_TYPE_MEDIA]

	medias, err = repository.CategoriesRepo.GetMediaIdentityList(request.Context(), tagged_medias_uuids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsPaginatedHandler, while getting media identities: %s\n", err))
		response.WriteHeader(500)
		return
	}

	tagged_categories_uuids := tagged_content[dungeon_models.ENTITY_TYPE_CATEGORY]

	for _, category_uuid := range tagged_categories_uuids {
		category_media_identities, err := repository.CategoriesRepo.GetCategoryMediaIdentities(request.Context(), category_uuid)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getEntitiesWithTagsPaginatedHandler, while getting category media identities: %s\n", err))
			response.WriteHeader(500)
			return
		}

		medias = append(medias, category_media_identities...)
	}

	var medias_length int = len(medias)

	var medias_in_page int = page_size

	var starting_index int = (page_num - 1) * page_size
	var ending_index int = starting_index + page_size
	if ending_index > medias_length {
		ending_index = medias_length - 1
		medias_in_page = medias_length - starting_index
	}

	var total_pages int = int(math.Ceil(float64(len(medias)) / float64(page_size)))

	requested_media_identities := make([]dungeon_models.MediaIdentity, medias_in_page)

	if starting_index > len(medias) {
		dungeon_helpers.WritePaginatedResponseList(response, make([]dungeon_models.MediaIdentity, 0), page_num, total_pages, len(medias))
		return
	}

	for starting_index <= ending_index && starting_index < len(medias) {
		requested_media_identities[starting_index%page_size] = medias[starting_index]
		starting_index++
	}

	dungeon_helpers.WritePaginatedResponseList(response, requested_media_identities, page_num, total_pages, len(medias))
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
