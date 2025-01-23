package common_flows

import (
	"fmt"
	"io"
	"libery-dungeon-libs/dungeonsec/access_sec"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	service_helpers "libery_medias_service/helpers"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/gotd/contrib/http_range"
)

/* ----------------------------- Cluster cookie ----------------------------- */

/**
 * Returns the fs path of a media. Expects the media_path to be the full path only missing the cluster fs path.
 * And attempts to get the category cluster from the request. If trust is set to false, it will attempt to get
 * the category cluster from the cookie app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME. If trust is set to true,
 * it will use InferCategoryCluster to get the category cluster.
 */
func GetMediaFsPathFromRequest(request *http.Request, media_path, cluster_uuid string) (string, error) {
	var category_cluster *dungeon_models.CategoryCluster
	var err error

	category_cluster, err = access_sec.GetSignedClusterOnRequest(cluster_uuid, request)
	if err != nil {
		return "", err
	}

	fs_media_path := filepath.Join(category_cluster.FsPath, media_path)

	return fs_media_path, nil
}

/* -------------------------------------------------------------------------- */

/**
 * Writes a media file to a response writer object.
 */
func WriteMediaFileResponse(media_identity dungeon_models.MediaIdentity, response http.ResponseWriter, request *http.Request) {
	var media_path string = media_identity.FsPath()

	file_descriptor, err := service_helpers.GetFileDescriptor(media_path)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("In workflows/common/request_workflows.WriteMEdiaFileResponse: Error getting media file descriptor for '%s': %s", media_path, err.Error()))
		dungeon_helpers.WriteRejection(response, 404, "file doesn't exists")
		return
	}
	defer file_descriptor.Close()

	file_header := make([]byte, 512)
	file_descriptor.Read(file_header)

	content_type := http.DetectContentType(file_header)

	file_stat, err := file_descriptor.Stat()
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("In workflows/common/request_workflows.WriteMEdiaFileResponse: Error getting media file stat: %s", err.Error()))
		dungeon_helpers.WriteRejection(response, 500, "")
		return
	}

	file_size_64 := file_stat.Size()
	var safe_filename string = service_helpers.RemoveSpecialChars(media_identity.Media.Name)

	response.Header().Set("Content-Type", content_type)
	response.Header().Set("Content-Disposition", "attachment; filename="+safe_filename)
	response.Header().Set("Cache-Control", "private, max-age=10800, must-revalidate")
	response.Header().Set("Accept-Ranges", "bytes")

	ranges, err := http_range.ParseRange(request.Header.Get("Range"), file_size_64)
	if err != nil || len(ranges) == 0 {
		file_descriptor.Seek(0, 0)
		var file_data []byte

		file_data, err := io.ReadAll(file_descriptor)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("In workflows/common/request_workflows.WriteMEdiaFileResponse: Error reading media file: %s", err.Error()))
			dungeon_helpers.WriteRejection(response, 500, "")
			return
		}

		response.Header().Set("Content-Length", strconv.FormatInt(file_size_64, 10))

		response.WriteHeader(200)

		response.Write(file_data)

		return
	}

	for _, r := range ranges {
		start := r.Start
		end := start + r.Length

		response.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end-1, file_size_64))
		response.Header().Set("Content-Length", strconv.FormatInt(r.Length, 10))
		response.WriteHeader(206)

		file_descriptor.Seek(start, 0)
		_, err := io.CopyN(response, file_descriptor, r.Length)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("In workflows/common/request_workflows.WriteMEdiaFileResponse: Error copying media file: %s", err.Error()))
		}
	}
}
