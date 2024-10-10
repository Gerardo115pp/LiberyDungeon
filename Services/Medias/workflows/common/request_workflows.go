package common_flows

import (
	"fmt"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	"libery_medias_service/repository"
	"net/http"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

/**
 * Attempts get the category cluster from the request first by checking if the cluster access cookie. If that fails, it will attempt to get the
 * cluster from the category id and the unsafe cluster path header. If the cluster access cookie is set. The category id can be an empty string
 */
func GetRequestCategoryCluster(request *http.Request, category_id string) (*dungeon_models.CategoryCluster, error) {
	var cluster_access_cookie *http.Cookie
	var category_cluster *dungeon_models.CategoryCluster
	var err error

	unsafe_cluster_path := request.Header.Get(app_config.IH_DOWNLOAD_RECIPIENT_PATH)

	cluster_access_cookie, err = request.Cookie(app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME)
	if err == nil {
		category_cluster, err = dungeon_models.GetCategoriesClusterFromToken(cluster_access_cookie.Value, app_config.JWT_SECRET)
		if err != nil {
			return nil, fmt.Errorf("Error getting cluster from token: %s.", err.Error())
		}
	} else if unsafe_cluster_path != "" && category_id != "" {
		echo.EchoDebug(fmt.Sprintf("Will attempt to get the cluster from the category. This will throw an error if the header '%s' is not set on the request and matches the cluster fs path.", app_config.IH_DOWNLOAD_RECIPIENT_PATH))

		category_cluster, err = repository.CategoriesClustersRepo.GetCategoryCluster(request.Context(), category_id)
		if err != nil {
			return nil, fmt.Errorf("Error getting category cluster: %s", err.Error())
		}
	}

	if category_cluster == nil {
		err = fmt.Errorf("Failed to determine category cluster")
	}

	return category_cluster, err
}

// Extracts the category cluster from the cookie app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME. If successful, it will construct a category cluster object from the
// jwt data in the cookie.
func GetCategoryClusterFromCookie(request *http.Request) (*dungeon_models.CategoryCluster, error) {
	cluster_access_cookie, err := request.Cookie(app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error getting cluster cookie: %s", err.Error())
	}

	category_cluster, err := dungeon_models.GetCategoriesClusterFromToken(cluster_access_cookie.Value, app_config.JWT_SECRET)
	if err != nil {
		return nil, fmt.Errorf("Error getting cluster from token: %s", err.Error())
	}

	return category_cluster, nil
}

/*
* Attempts to infer the category cluster from the request. It does the following:
  - 1. Tries to get the category cluster from the cookie app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME.
  - 2. If that fails, attempts to get it either from a category or media. So it checks if either the category_id or media_id were send as query parameters.
  - 3. If both category_id and media_id are empty, it will return an error. If category_id is empty but media_id is not, it will get the media and use its main category
    to get the category cluster. If category_id is not empty, it will get the category cluster from the category_id.
*/
func InferCategoryCluster(request *http.Request) (*dungeon_models.CategoryCluster, error) {
	var request_category_cluster *dungeon_models.CategoryCluster
	var err error

	request_category_cluster, err = GetCategoryClusterFromCookie(request)
	if err == nil && request_category_cluster != nil {
		return request_category_cluster, nil
	}

	var category_id string = request.URL.Query().Get("category_id")
	var media_id string = request.URL.Query().Get("media_id")

	if category_id == "" && media_id == "" {
		return nil, fmt.Errorf("No category_id or media_id provided")
	}

	if category_id == "" {
		media, err := repository.MediasRepo.GetMediaByID(request.Context(), media_id)
		if err != nil {
			return nil, fmt.Errorf("Error getting media: %s", err.Error())
		}

		category_id = media.MainCategory
	}

	request_category_cluster, err = repository.CategoriesClustersRepo.GetCategoryCluster(request.Context(), category_id)
	if err != nil {
		return nil, fmt.Errorf("Error getting category cluster: %s", err.Error())
	}

	return request_category_cluster, err
}

/*
*

		Returns the fs path of a media. Expects the media_path to be the full path only missing the cluster fs path.
	 	And attempts to get the category cluster from the request. If trust is set to false, it will attempt to get
		the category cluster from the cookie app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME. If trust is set to true,
		it will use InferCategoryCluster to get the category cluster.
*/
func GetMediaFsPathFromRequest(request *http.Request, media_path string, trust bool) (string, error) {
	var category_cluster *dungeon_models.CategoryCluster
	var err error

	if trust {
		category_cluster, err = InferCategoryCluster(request)
	} else {
		category_cluster, err = GetCategoryClusterFromCookie(request)
	}

	if err != nil {
		return "", err
	}

	fs_media_path := filepath.Join(category_cluster.FsPath, media_path)

	return fs_media_path, nil
}
