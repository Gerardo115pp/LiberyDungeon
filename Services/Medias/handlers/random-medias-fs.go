package handlers

import (
	"fmt"
	"io"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	service_helpers "libery_medias_service/helpers"
	"libery_medias_service/repository"
	"net/http"
	"path"
	"strconv"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/gotd/contrib/http_range"
)

func RandomMediasFsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getRandomMediasFsHandler(response, request)
		case http.MethodPost:
			postRandomMediasFsHandler(response, request)
		case http.MethodPatch:
			patchRandomMediasFsHandler(response, request)
		case http.MethodDelete:
			deleteRandomMediasFsHandler(response, request)
		case http.MethodPut:
			putRandomMediasFsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getRandomMediasFsHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_id string = request.URL.Query().Get("cluster_id")
	var category_id string = request.URL.Query().Get("category_id")
	var cache_seconds string = request.URL.Query().Get("cache_seconds")
	var cluster dungeon_models.CategoryCluster
	var media *dungeon_models.Media
	var category *dungeon_models.Category
	var err error

	if cluster_id == "" {
		cluster_cookie, err := request.Cookie(app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting cluster access cookie because '%s'", err.Error()))
			response.WriteHeader(404)
			return
		}

		token_cluster, err := dungeon_models.GetCategoriesClusterFromToken(cluster_cookie.Value, app_config.JWT_SECRET)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting cluster from token because '%s'", err.Error()))
			response.WriteHeader(404)
			return
		}

		cluster = *token_cluster
	} else {
		cluster, err = repository.CategoriesClustersRepo.GetClusterByID(request.Context(), cluster_id)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting cluster by id because '%s'", err.Error()))
			response.WriteHeader(404)
			return
		}

	}

	media, category, err = repository.MediasRepo.GetRandomMedia(request.Context(), cluster.Uuid, category_id, true)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting random media because '%s'", err.Error()))
		response.WriteHeader(404)
		return
	}

	file_descriptor, err := service_helpers.GetFileDescriptor(path.Join(cluster.FsPath, category.Fullpath, media.Name))
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting media file descriptor because '%s'", err.Error()))
		response.WriteHeader(404)
		return
	}

	defer file_descriptor.Close()
	file_header := make([]byte, 512)
	file_descriptor.Read(file_header)

	content_type := http.DetectContentType(file_header)

	file_stat, err := file_descriptor.Stat()
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error getting media file stat because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	var use_cache bool = false

	if cache_seconds != "" {
		_, err := strconv.Atoi(cache_seconds)

		use_cache = err == nil
	}

	file_size_64 := file_stat.Size()
	safe_filename := service_helpers.RemoveSpecialChars(media.Name)

	response.Header().Set("Content-Type", content_type)
	response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", safe_filename))
	response.Header().Set("Accept-Ranges", "bytes")

	if !use_cache {
		response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	} else {
		response.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", cache_seconds))
	}

	ranges, err := http_range.ParseRange(request.Header.Get("Range"), file_size_64)
	if err != nil || len(ranges) == 0 {
		file_descriptor.Seek(0, 0)
		var file_data []byte

		file_data, err := io.ReadAll(file_descriptor)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error reading media file because '%s'", err.Error()))
			response.WriteHeader(500)
			return
		}

		response.Header().Set("Content-Length", fmt.Sprintf("%d", file_size_64))

		response.WriteHeader(200)

		_, err = response.Write(file_data)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error writing media file because '%s'", err.Error()))
			response.WriteHeader(500)
			return
		}

		return
	}

	for _, rng := range ranges {
		start := rng.Start
		end := start + rng.Length

		response.Header().Set("Content-Ranges", fmt.Sprintf("bytes %d-%d/%d", start, end-1, file_size_64))
		response.Header().Set("Content-Length", strconv.FormatInt(rng.Length, 10))
		response.WriteHeader(206)

		file_descriptor.Seek(start, 0)
		_, err = io.CopyN(response, file_descriptor, rng.Length)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In getRandomMediasFsHandler: Error copying media file range because '%s'", err.Error()))
			response.WriteHeader(500)
			return
		}
	}

	// fmt.Println("Media: ", media)
	// fmt.Println("Category: ", category)

	// response.Header().Set("Content-Type", "application/json")
	// response.WriteHeader(200)
	// json.NewEncoder(response).Encode(media)
}
func postRandomMediasFsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchRandomMediasFsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteRandomMediasFsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putRandomMediasFsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
