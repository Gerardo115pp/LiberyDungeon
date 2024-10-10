package handlers

// a cluster can be created in:
// - a new directory on /. will not create a new cluster, on / if a directory with the same name exists.
// - any directory under app_config.SERVICE_CLUSTERS_ROOT
// - any new directory or subdirectory of /home
// To add an existing directory as a cluster, the directory must be under app_config.SERVICE_CLUSTERS_ROOT.
// Things to watch out for:
// - Directory variables that will result in a new cluster `fs_root` property shall not include symbols like `~` or `..`

import (
	"encoding/json"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	app_config "libery_categories_service/Config"
	"libery_categories_service/handlers/request_parameters"
	service_models "libery_categories_service/models"
	service_workflows "libery_categories_service/workflows"
	common_workflows "libery_categories_service/workflows/common"
	service_fs_workflows "libery_categories_service/workflows/servicefs_workflows"
	"libery_categories_service/workflows/servicefs_workflows/fs_sync"
	"net/http"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func ServiceFSHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getServiceFSHandler(response, request)
		case http.MethodPost:
			postServiceFSHandler(response, request)
		case http.MethodPatch:
			patchServiceFSHandler(response, request)
		case http.MethodDelete:
			deleteServiceFSHandler(response, request)
		case http.MethodPut:
			putServiceFSHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	switch resource_path {
	case "/service-fs/new-cluster-options":
		getNewClusterOptionsServiceFSHandler(response, request)
	case "/service-fs/validate-path":
		getNewClusterPathValidationHandler(response, request)
	case "/service-fs/cluster-root-default":
		getClusterRootDefaultHandler(response, request)
	default:
		http.Error(response, "Resource not found", 404)
	}
	return
}

// The idea of this handler is to return a list of directories that can be used to create a new cluster within the service.
func getNewClusterOptionsServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	var subdirectory_to_scan string = request.URL.Query().Get("subdirectory")
	var new_cluster_options []service_models.DirectoryOption
	var directory_to_scan string = app_config.SERVICE_CLUSTERS_ROOT

	if subdirectory_to_scan != "" {
		if dungeon_helpers.IsChildPath(app_config.SERVICE_CLUSTERS_ROOT, subdirectory_to_scan) {
			directory_to_scan = subdirectory_to_scan
		} else {
			http.Error(response, "Invalid subdirectory path", 400)
			return
		}
	}

	directory_to_scan, err := filepath.EvalSymlinks(directory_to_scan)
	if err != nil {
		echo.EchoErr(fmt.Errorf("Error evaluating symlinks: %s", err.Error()))
		http.Error(response, "Could not get directory options", 500)
		return
	}

	directory_to_scan, err = filepath.Abs(directory_to_scan)
	if err != nil {
		echo.EchoErr(fmt.Errorf("Error getting absolute path: %s", err.Error()))
		http.Error(response, "Error getting directory options", 500)
		return
	}
	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Scanning directory: %s", directory_to_scan))

	new_cluster_options, labeled_err := service_fs_workflows.GetDirectoryOptionsFromPath(directory_to_scan)
	if labeled_err != nil {
		var status_code int = 500
		if labeled_err.Label == service_workflows.ErrNoSuchFileOrDirectory {
			status_code = 404
		} else if labeled_err.Label == service_workflows.ErrForbiddenDirectoryScan {
			status_code = 403
		}

		echo.EchoErr(labeled_err)
		http.Error(response, "could not get directory options", status_code)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	response.WriteHeader(200)

	err = json.NewEncoder(response).Encode(new_cluster_options)
	if err != nil {
		http.Error(response, "Error encoding new cluster options", 500)
	}
}

/**
 * Returns a boolean response indicating if the path is valid for creating a new cluster.
 */
func getNewClusterPathValidationHandler(response http.ResponseWriter, request *http.Request) {
	var path_to_validate string = request.URL.Query().Get("unsafe_path")

	if path_to_validate == "" {
		echo.Echo(echo.RedFG, "Empty paths are never valid")
		dungeon_helpers.WriteBooleanResponse(response, false)
		return
	}

	var is_path_valid bool = true
	var rejection_reason string = ""

	labeled_err := service_fs_workflows.VerifyValidSubdirectoryScanPath(path_to_validate, true)
	if labeled_err != nil {
		echo.EchoErr(labeled_err)
		fmt.Printf(string(echo.RESET))

		is_path_valid = false

		switch labeled_err.Label {
		case service_workflows.ErrForbiddenDirectoryScan:
			rejection_reason = fmt.Sprintf("Creating a cluster in the directory '%s' is forbidden by the security policy", path_to_validate)
		case service_workflows.ErrNoSuchDirectory:
			rejection_reason = fmt.Sprintf("The directory '%s' does not exist or is not a directory", path_to_validate)
		case service_workflows.ErrPathValidation_PathIsAClusterAncestor:
			conflicting_cluster_name, lerr := labeled_err.GetStringVariable("cluster_name")
			if lerr != nil {
				echo.EchoErr(lerr)
				rejection_reason = fmt.Sprintf("The directory '%s' contains another cluster", path_to_validate)
			} else {
				rejection_reason = fmt.Sprintf("The directory '%s' contains the cluster '%s'", path_to_validate, conflicting_cluster_name)
			}
		case service_workflows.ErrPathValidation_PathIsACluster:
			rejection_reason = fmt.Sprintf("The directory '%s' is already registered as a cluster", path_to_validate)
		case service_workflows.ErrPathValidation_PathIsAClusterChild:
			conflicting_cluster_name, lerr := labeled_err.GetStringVariable("cluster_name")
			if lerr != nil {
				echo.EchoErr(lerr)
				rejection_reason = fmt.Sprintf("The directory '%s' is a child of another cluster", path_to_validate)
			} else {
				rejection_reason = fmt.Sprintf("The directory '%s' is a child of the cluster '%s'", path_to_validate, conflicting_cluster_name)
			}
		}
	}

	dungeon_helpers.WriteReasonedBooleanResponse(response, is_path_valid, rejection_reason)
}

/**
 * Returns the value of the app_config.SERVICE_CLUSTERS_ROOT variable.
 */
func getClusterRootDefaultHandler(response http.ResponseWriter, request *http.Request) {
	dungeon_helpers.WriteSingleStringResponse(response, app_config.SERVICE_CLUSTERS_ROOT)
}

func postServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putServiceFSHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/service-fs/cluster-sync":
		handler_func = putClusterSyncHandler
	}

	handler_func(response, request)
}

func putClusterSyncHandler(response http.ResponseWriter, request *http.Request) {
	var params *request_parameters.SyncClusterPathRequest

	params, err := request_parameters.NewSyncClusterPathFromRequest(request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error decoding request: %s", err.Error()))
		http.Error(response, "Error decoding request", 400)
		return
	}

	category_identity, lerr := common_workflows.GetCategoryIdentityFromUUID(params.SyncCategoryUUID)
	if lerr != nil {
		echo.EchoErr(lerr)
		http.Error(response, "Problematic identity requested", 404)
		return
	}

	if category_identity.ClusterUUID != params.ClusterUUID {
		echo.Echo(echo.RedFG, fmt.Sprintf("The cluster '%s' does not match the cluster of the category '%s'", params.ClusterUUID, params.SyncCategoryUUID))
		http.Error(response, "Cluster mismatch", 400)
		return
	}

	lerr = fs_sync.SyncCategoryBranch(category_identity)
	if lerr != nil {
		echo.EchoErr(lerr)
		http.Error(response, "Error syncing category branch", 500)
		return
	}

	response.WriteHeader(200)
	return
}
