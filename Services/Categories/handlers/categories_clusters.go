package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/dungeonsec"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	service_models "libery_categories_service/models"
	"libery_categories_service/repository"
	service_fs_workflows "libery_categories_service/workflows/servicefs_workflows"
	"net/http"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CategoriesClustersHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getCategoriesClustersHandler(response, request)
		case http.MethodPost:
			postCategoriesClustersHandler(response, request)
		case http.MethodPatch:
			patchCategoriesClustersHandler(response, request)
		case http.MethodDelete:
			deleteCategoriesClustersHandler(response, request)
		case http.MethodPut:
			putCategoriesClustersHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	// if resource_path == "/clusters" {
	// 	getAllCategoriesClustersHandler(response, request)
	// } else if resource_path == "/clusters/sign-access" {
	// 	getSignAccessCategoriesClustersHandler(response, request)
	// }
	switch resource_path {
	case "/clusters":
		handler_func = getAllCategoriesClustersHandler
	case "/clusters/sign-access":
		handler_func = getSignAccessCategoriesClustersHandler
	}

	handler_func(response, request)
}

func getAllCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var categories_clusters []dungeon_models.CategoryCluster
	var err error

	var has_private_cluster_access bool = false
	var user_claims *dungeon_models.PlatformUserClaims

	user_claims, err = dungeon_middlewares.GetUserClaims(request, app_config.JWT_SECRET)
	if err == nil {
		has_private_cluster_access = dungeonsec.CanViewPrivateClusters(user_claims.UserGrants)
		echo.Echo(echo.GreenFG, fmt.Sprintf("User '%s' can view private clusters: %t", user_claims.UserUUID, has_private_cluster_access))
	}

	echo.Echo(echo.GreenBG, "GET", "CategoriesClustersHandler", "Get all categories clusters")

	categories_clusters, err = repository.CategoriesClustersRepo.GetClusters(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting categories clusters: %s", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var filtered_categories_clusters []dungeon_models.CategoryCluster = make([]dungeon_models.CategoryCluster, 0)
	var private_cluster_uuids []string = make([]string, 0)

	if !has_private_cluster_access {
		private_cluster_uuids, err = communication.Metadata.GetAllPrivateClusters()
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/categories_clusters.go getAllCategoriesClustersHandler :: Error getting all private clusters from metadata service: %s", err))
			http.Error(response, "Can't correctly form the response", 500)
			return
		}
	}

	for _, cluster := range categories_clusters {
		var is_opaque bool = true
		var is_private_cluster bool

		if !has_private_cluster_access {
			for _, private_cluster_uuid := range private_cluster_uuids {
				if cluster.Uuid == private_cluster_uuid {
					is_private_cluster = true
					break
				}
			}

			is_opaque = !is_private_cluster
		}

		if is_opaque {
			filtered_categories_clusters = append(filtered_categories_clusters, cluster)
		}
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	err = json.NewEncoder(response).Encode(filtered_categories_clusters)
	if err != nil {
		echo.Echo(echo.RedFG, "GET", "CategoriesClustersHandler", "Error encoding categories clusters")
	}
}

func getSignAccessCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_id string = request.URL.Query().Get("cluster")
	var cluster dungeon_models.CategoryCluster
	access_response := &dungeon_models.ClaimsStandardResponse{
		RedirectURL: "/",
		Granted:     false,
	}

	cluster, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), cluster_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting cluster: %s", err.Error()))
		response.WriteHeader(404) // Not found
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(access_response)
		return
	}

	var cluster_is_private bool

	cluster_is_private, err = communication.Metadata.CheckClusterPrivate(cluster.Uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error checking if cluster is private: %s", err.Error()))
		response.WriteHeader(500)
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(access_response)
		return
	}

	if cluster_is_private {
		user_claims, err := dungeon_middlewares.GetUserClaims(request, dungeon_secrets.GetDungeonJwtSecret())
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error getting user claims: %s", err.Error()))
			response.WriteHeader(500)
			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(access_response)
			return
		}

		var user_can_access bool = dungeonsec.CanViewPrivateClusters(user_claims.UserGrants)
		if !user_can_access {
			response.WriteHeader(403) // Forbidden
			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(access_response)
			return
		}
	}

	access_response.RedirectURL = fmt.Sprintf("%s/%s", app_config.MEDIAS_APP_DUNGEON_EXPLORER_ROUTE, cluster.RootCategory)
	access_response.Granted = true

	var access_expiration_time time.Time = time.Now().Add(time.Hour * 24)

	var access_token string
	access_token, err = dungeon_models.GenerateCategoriesClusterAccess(cluster, access_expiration_time, app_config.JWT_SECRET)

	var access_cookie http.Cookie = http.Cookie{
		Name:     app_config.CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME,
		Value:    access_token,
		Path:     "/",
		Expires:  access_expiration_time,
		HttpOnly: true,
	}

	http.SetCookie(response, &access_cookie)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(access_response)

	// TODO: check if the user has ACCESS to the cluster.
}

func postCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var new_cluster_data *dungeon_models.CategoryCluster = new(dungeon_models.CategoryCluster)

	err := json.NewDecoder(request.Body).Decode(new_cluster_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postCategoriesClustersHandler: there was an error decoding the request body because '%s'", err.Error()))
		http.Error(response, "Error decoding the request body", 400)
		return
	}

	var new_cluster *dungeon_models.CategoryCluster

	new_cluster, labeled_err := service_fs_workflows.CreateNewCategoryCluster(new_cluster_data)
	if labeled_err != nil {
		echo.EchoErr(labeled_err)
		http.Error(response, "Error creating new cluster", 500)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(201)

	err = json.NewEncoder(response).Encode(new_cluster)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In postCategoriesClustersHandler: there was an error encoding the response body because '%s'", err.Error()))
		http.Error(response, "Error encoding the response body", 500)
		return
	}
}

func patchCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var new_cluster_data dungeon_models.CategoryCluster
	response.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&new_cluster_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchCategoriesClustersHandler: there was an error decoding the request body because '%s'", err.Error()))
		http.Error(response, "Error decoding the request body", 400)
		return
	}

	old_cluster_data, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), new_cluster_data.Uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchCategoriesClustersHandler: there was an error getting the cluster because '%s'", err.Error()))
		http.Error(response, "Error getting the cluster, likely it does not exist", 404)
		return
	}

	err = old_cluster_data.UpdateClusterData(&new_cluster_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchCategoriesClustersHandler: there was an error updating the cluster because '%s'", err.Error()))
		http.Error(response, "Error updating the cluster", 422) // Unprocessable Entity
		return
	}

	err = repository.CategoriesClustersRepo.UpdateCluster(request.Context(), old_cluster_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In patchCategoriesClustersHandler: there was an error updating the cluster because '%s'", err.Error()))
		http.Error(response, "Error updating the cluster", 500)
		return
	}

	response.WriteHeader(200)
	json.NewEncoder(response).Encode(old_cluster_data)
}

func deleteCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	if resource_path == "/clusters/platform-data" {
		deletePlatformDataCategoriesClustersHandler(response, request)
	} else {
		response.WriteHeader(http.StatusNotFound)
	}

}

func deletePlatformDataCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_id string = request.URL.Query().Get("cluster")
	var labeled_err *dungeon_models.LabeledError

	echo.EchoDebug(fmt.Sprintf("DELETE /clusters/platform-data?cluster=%s", cluster_id))

	labeled_err = repository.CategoriesClustersRepo.DeleteCluster(request.Context(), cluster_id)
	if labeled_err != nil {
		echo.EchoErr(labeled_err)
		if labeled_err.Label == service_models.ErrDB_CouldNotFindCategoryCluster {
			http.Error(response, "Cluster not found", 404)
			return
		} else if labeled_err.Label == dungeon_models.ErrDB_FailedToCommitTX {
			echo.EchoWarn("Failed to commit transaction")
		}

		http.Error(response, "Error deleting cluster", 500)
	}

	response.WriteHeader(204)
}

func putCategoriesClustersHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
