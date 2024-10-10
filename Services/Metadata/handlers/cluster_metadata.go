package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	"libery-metadata-service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var cluster_metadata_path string = "/metadata/clusters(/.+)?"

var CLUSTER_METADATA_ROUTE *patriot_router.Route = patriot_router.NewRoute(cluster_metadata_path, false)

func ClusterMetadataHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = dungeon_middlewares.CheckUserCanViewPrivateClusters(getClusterMetadataHandler)
		case http.MethodPost:
			request_handler_func = dungeon_middlewares.CheckUserCan_AlterPrivateClusters(postClusterMetadataHandler)
		case http.MethodPatch:
			request_handler_func = patchClusterMetadataHandler
		case http.MethodDelete:
			request_handler_func = deleteClusterMetadataHandler
		case http.MethodPut:
			request_handler_func = putClusterMetadataHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getClusterMetadataHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/metadata/clusters/private-clusters":
		getAllPrivateClustersHandler(response, request)
	case "/metadata/clusters/is-private":
		getIsPrivateHandler(response, request)
	default:
		dungeon_helpers.ResourceNotFoundHandler(response, request)
	}
}

func getAllPrivateClustersHandler(response http.ResponseWriter, request *http.Request) {
	var private_cluster_uuids []string = make([]string, 0)

	private_cluster_uuids = repository.ClusterMetadataRepo.GetPrivateClusters()

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(private_cluster_uuids)
}

func getIsPrivateHandler(response http.ResponseWriter, request *http.Request) {
	var query_cluster_uuid string = request.URL.Query().Get("cluster_uuid")

	if query_cluster_uuid == "" {
		echo.Echo(echo.RedBG, fmt.Sprintf("In Handlers/clusters_metadata.go getIsPrivateHandler :: While getting query_cluster_uuid query parameter<%s>", query_cluster_uuid))
		response.WriteHeader(400)
		return
	}

	var is_private bool = repository.ClusterMetadataRepo.IsPrivateCluster(query_cluster_uuid)

	dungeon_helpers.WriteBooleanResponse(response, is_private)
}

func postClusterMetadataHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/metadata/clusters/privacy":
		postClusterPrivacyHandler(response, request)
	default:
		dungeon_helpers.ResourceNotFoundHandler(response, request)
	}
}

func postClusterPrivacyHandler(response http.ResponseWriter, request *http.Request) {
	var cluster_uuid string = request.URL.Query().Get("cluster_uuid")
	var cluster_privacy string = request.URL.Query().Get("privacy")

	if cluster_uuid == "" {
		echo.Echo(echo.RedBG, fmt.Sprintf("In Handlers/clusters_metadata.go postClusterPrivacyHandler :: While getting cluster_uuid query parameter<%s>", cluster_uuid))
		response.WriteHeader(400)
		return
	}

	var is_private bool = false

	if cluster_privacy == "private" {
		is_private = true
	}

	if is_private {
		repository.ClusterMetadataRepo.AddPrivateCluster(cluster_uuid)
	} else {
		repository.ClusterMetadataRepo.RemovePrivateCluster(cluster_uuid)
	}

	response.WriteHeader(200)
}

func patchClusterMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteClusterMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putClusterMetadataHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
