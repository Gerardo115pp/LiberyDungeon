package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	"libery_medias_service/repository"
	"libery_medias_service/workflows"
	common_flows "libery_medias_service/workflows/common"

	"net/http"
	"strconv"
	"strings"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var medias_resource_path string = "/medias(/.+)?"

var MEDIAS_ROUTE *patriot_router.Route = patriot_router.NewRoute(medias_resource_path, false)

func MediasHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getMediasHandler(response, request)
		case http.MethodPost:
			postMediasHandler(response, request)
		case http.MethodPatch:
			patchMediasHandler(response, request)
		case http.MethodDelete:
			deleteMediasHandler(response, request)
		case http.MethodPut:
			putMediasHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getMediasHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	handler_func := dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/medias/by-uuid":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getMediaByUUIDHandler)
	case "/medias/identity":
		handler_func = dungeon_middlewares.CheckUserCan_ViewContent(getMediaIdentityHandler)
	default:
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediasHandler: Resource not found: %s", resource))
	}

	handler_func(response, request)
}

func getMediaByUUIDHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("uuid")

	if media_uuid == "" {
		echo.Echo(echo.RedBG, "In MediasService.medias.getMediaByUUIDHandler: Missing media uuid query parameter")
		response.WriteHeader(400)
		return
	}

	media, err := repository.MediasRepo.GetMediaByID(request.Context(), media_uuid)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaByUUIDHandler: Error getting media by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=10600") // 3 hours

	response.WriteHeader(200)
	json.NewEncoder(response).Encode(media)
}

// Returns a MediaIdentity item which is an extension of the Media struct which includes all the necessary information from the category and cluster to request a media file.
func getMediaIdentityHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("uuid")

	access_cluster, err := common_flows.GetCategoryClusterFromCookie(request)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Error getting category cluster from cookie: %s", err.Error()))
		response.WriteHeader(403)
		return
	}

	media, err := repository.MediasRepo.GetMediaByID(request.Context(), media_uuid)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Error getting media by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	category, err := repository.CategoriesRepo.GetCategoryByID(request.Context(), media.MainCategory)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Error getting category by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), category.Cluster)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Error getting category cluster by ID: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	if category_cluster.Uuid != access_cluster.Uuid {
		echo.Echo(echo.RedBG, fmt.Sprintf("In MediasService.medias.getMediaIdentityHandler: Access denied to category cluster"))
		response.WriteHeader(403)
		return
	}

	media_identity := dungeon_models.CreateNewMediaIdentity(media, &category, &category_cluster)

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=10600") // 3 hours

	response.WriteHeader(200)

	json.NewEncoder(response).Encode(media_identity)
}

// Deprecated: This handler will be removed once other services that use it are updated.
func postMediasHandler(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.RedBG, "DEPRECATED: use POST '/upload-streams/stream-fragment' instead")
	var category_cluster *dungeon_models.CategoryCluster
	var downloaded_from_uuid int64 = 0
	var err error

	main_category_uuid := request.URL.Query().Get("main_category")
	if main_category_uuid == "" {
		http.Error(response, "Missing main_category query parameter", 400)
		echo.Echo(echo.RedBG, fmt.Sprintf("Missing main_category query parameter"))
		return
	}

	category_cluster, err = common_flows.GetRequestCategoryCluster(request, main_category_uuid)
	if err != nil {
		http.Error(response, "Error getting category cluster", 400)
		echo.Echo(echo.RedBG, fmt.Sprintf("In postMediasHandler: Error getting category cluster because '%s'", err.Error()))
		return
	}

	download_from := request.URL.Query().Get("download_from")

	if download_from != "" {
		downloaded_from_uuid, err = strconv.ParseInt(download_from, 10, 64)
		if err != nil {
			http.Error(response, "Invalid download_from query parameter", 400)
			echo.Echo(echo.RedBG, fmt.Sprintf("Invalid download_from query parameter: %s", err.Error()))
			return
		}
	}

	main_category, err := repository.CategoriesRepo.GetCategoryByID(request.Context(), main_category_uuid)
	if err != nil {
		http.Error(response, "Error getting main category", 404)
		echo.Echo(echo.RedBG, fmt.Sprintf("Error getting main category: %s", err.Error()))
		return
	}

	err = request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(response, "Error parsing multipart form", 400)
		echo.Echo(echo.RedBG, fmt.Sprintf("Error parsing multipart form: %s", err.Error()))
		return
	}

	for _, file_headers := range request.MultipartForm.File {
		for _, file_header := range file_headers {
			file, err := file_header.Open()
			if err != nil {
				http.Error(response, "Error opening file", 400)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error opening file: %s", err.Error()))
				return
			}
			defer file.Close()

			// Get mime type
			mime_type := file_header.Header.Get("Content-Type")

			is_video := strings.HasPrefix(mime_type, "video")

			media := dungeon_models.CreateNewMedia(file_header.Filename, main_category_uuid, is_video, downloaded_from_uuid)

			echo.Echo(echo.CyanFG, fmt.Sprintf("About to insert media: %s", media.Uuid))

			media_identity := dungeon_models.CreateNewMediaIdentity(media, &main_category, category_cluster)

			err = workflows.SaveMediaFile(media_identity, &file)
			if err != nil {
				http.Error(response, "Error saving media file", 500)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error saving media file: %s", err.Error()))
				return
			}

			err = repository.MediasRepo.InsertMedia(request.Context(), media)
			if err != nil {
				http.Error(response, "Error inserting media", 500)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error inserting media: %s", err.Error()))
				return
			}
		}
	}

	response.WriteHeader(201)
	return
}

func patchMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putMediasHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
