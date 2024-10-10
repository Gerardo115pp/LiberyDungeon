package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	"libery_medias_service/helpers"
	service_helpers "libery_medias_service/helpers"
	"libery_medias_service/workflows"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/gotd/contrib/http_range"
)

func MediasFSHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getMediasFSHandler(response, request)
		case http.MethodPost:
			postMediasFSHandler(response, request)
		case http.MethodPatch:
			patchMediasFSHandler(response, request)
		case http.MethodDelete:
			deleteMediasFSHandler(response, request)
		case http.MethodPut:
			putMediasFSHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getMediasFSHandler(response http.ResponseWriter, request *http.Request) {
	var fs_path string = ""
	var cluster_cookie *http.Cookie

	cluster_cookie, err := request.Cookie(app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME)
	if err == nil {
		cluster, err := dungeon_models.GetCategoriesClusterFromToken(cluster_cookie.Value, app_config.JWT_SECRET)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting cluster from token: %s", err.Error()))
			response.WriteHeader(400)
			return
		}

		fs_path = cluster.FsPath
	}

	media_path := request.URL.Path
	var use_mobile_version bool = false

	media_path = helpers.RemoveRoutePrefix(media_path, "medias-fs")
	if media_path == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Unprocessable media path: %s", request.URL.Path))
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(media_path, "/mobile") {
		media_path = strings.Replace(media_path, "/mobile", "", 1)
		use_mobile_version = true
	}

	if fs_path != "" {
		media_path = path.Join(fs_path, media_path)
	}

	file_descriptor, err := service_helpers.GetFileDescriptor(media_path)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting media file descriptor for '%s': %s", media_path, err.Error()))
		response.WriteHeader(http.StatusNotFound)
		return
	}

	defer file_descriptor.Close()
	file_header := make([]byte, 512)
	file_descriptor.Read(file_header)

	content_type := http.DetectContentType(file_header)

	is_video := strings.Contains(content_type, "video")

	file_stat, err := file_descriptor.Stat()
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting media file stat: %s", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	file_size_64 := file_stat.Size()
	var filename string

	segments := strings.Split(media_path, "/")

	if len(segments) == 0 {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Unprocessable media path: %s", request.URL.Path))
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	filename = segments[len(segments)-1]

	safe_filename := helpers.RemoveSpecialChars(filename)

	response.Header().Set("Content-Type", content_type)
	response.Header().Set("Content-Disposition", "attachment; filename="+safe_filename)
	// response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	response.Header().Set("Cache-Control", "private, max-age=10800, must-revalidate")
	response.Header().Set("Accept-Ranges", "bytes")

	ranges, err := http_range.ParseRange(request.Header.Get("Range"), file_size_64)
	if err != nil || len(ranges) == 0 {
		file_descriptor.Seek(0, 0)
		var file_data []byte

		file_data, err := ioutil.ReadAll(file_descriptor)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error reading media file: %s", err.Error()))
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Handle mobile resizing
		if use_mobile_version && !is_video {
			resized_image, err := workflows.ResizeImage(file_data, media_path, app_config.MOBILE_MAX_WIDTH)
			if err == nil {
				file_size_64 = int64(resized_image.Len())
				file_data = resized_image.Bytes()
			}
		}

		response.Header().Set("Content-Length", strconv.FormatInt(file_size_64, 10))

		response.WriteHeader(200)

		_, err = response.Write(file_data)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error writing media file to response: %s", err.Error()))
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	for _, r := range ranges {
		start := r.Start
		end := start + r.Length

		response.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end-1, file_size_64))
		response.Header().Set("Content-Length", strconv.FormatInt(r.Length, 10))
		response.WriteHeader(206) // Partial Content

		file_descriptor.Seek(start, 0)
		_, err = io.CopyN(response, file_descriptor, r.Length)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error copying media file to response: %s", err.Error()))
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

func postMediasFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchMediasFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteMediasFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putMediasFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
