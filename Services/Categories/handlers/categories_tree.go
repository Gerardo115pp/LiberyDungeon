package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CategoriesTreeHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getCategoriesTreeHandler(response, request)
		case http.MethodPost:
			postCategoriesTreeHandler(response, request)
		case http.MethodPatch:
			patchCategoriesTreeHandler(response, request)
		case http.MethodDelete:
			deleteCategoriesTreeHandler(response, request)
		case http.MethodPut:
			putCategoriesTreeHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postCategoriesTreeHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchCategoriesTreeHandler(response http.ResponseWriter, request *http.Request) {
	var resources_requested string = request.URL.Path

	switch resources_requested {
	case "/categories-tree/move":
		handleMoveCategory(response, request)
	default:
		response.WriteHeader(http.StatusNotFound)
	}
}

func handleMoveCategory(response http.ResponseWriter, request *http.Request) {
	move_category_request := &struct {
		NewParentCategory string `json:"new_parent_category"`
		MovedCategory     string `json:"moved_category"`
	}{}

	err := json.NewDecoder(request.Body).Decode(move_category_request)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error decoding move category request: %s", err.Error()))
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if move_category_request.NewParentCategory == "" || move_category_request.MovedCategory == "" {
		echo.Echo(echo.RedFG, "receiver_category and moved_category are required")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = workflows.MoveCategory(move_category_request.MovedCategory, move_category_request.NewParentCategory, request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handleMoveCategory: Error moving category because '%s'", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	modified_category, err := repository.CategoriesRepo.GetCategory(request.Context(), move_category_request.MovedCategory)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handleMoveCategory: Error getting modified category because '%s'", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(modified_category)
}

func deleteCategoriesTreeHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putCategoriesTreeHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func getCategoriesTreeHandler(response http.ResponseWriter, request *http.Request) {
	var request_path string = request.URL.Path

	echo.Echo(echo.GreenFG, fmt.Sprintf("Request path: %s", request_path))

	if "/categories-tree" == request_path {
		getCategoriesLeaf(response, request)
		return
	} else if "/categories-tree/short" == request_path {
		getCategoriesShort(response, request)
		return
	} else {
		echo.Echo(echo.RedFG, fmt.Sprintf("Invalid path: %s", request_path))
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getCategoriesLeaf(response http.ResponseWriter, request *http.Request) {
	var category_id string = request.URL.Query().Get("category_id")
	var category_path string = request.URL.Query().Get("category_path")
	var category_cluster string = request.URL.Query().Get("category_cluster")

	echo.EchoDebug("category_id: " + category_id)

	if category_id == "" && (category_path == "" || category_cluster == "") {
		http.Error(response, "Invalid request, missing all posible category identifiers. either pass category_id or both category_path and category_cluster", 400)
		return
	} else if category_id == "" {
		getCategoryLeafByFullpath(response, request)
		return
	}

	var category_content, err = repository.CategoriesRepo.GetCategoryContent(request.Context(), category_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting category content: %s", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	category_content.SortContentSeries()

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(category_content)
}

// Returns a list of childs categories names and uuids, nothing else, no media or other info
func getCategoriesShort(response http.ResponseWriter, request *http.Request) {
	var category_id string = request.URL.Query().Get("category_id")

	echo.EchoDebug("processing short categories tree")

	if category_id == "" {
		echo.Echo(echo.OrangeFG, "category_id is required")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var category_content, err = repository.CategoriesRepo.GetCategoryChildsByID(request.Context(), category_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting category content: %s", err.Error()))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(category_content)
}

func getCategoryLeafByFullpath(response http.ResponseWriter, request *http.Request) {
	var category_path string = request.URL.Query().Get("category_path")
	var category_cluster string = request.URL.Query().Get("category_cluster")
	echo.EchoDebug(fmt.Sprintf("category_path: %s, category_cluster: %s", category_path, category_cluster))

	if category_path == "" || category_cluster == "" {
		echo.Echo(echo.OrangeFG, "category_path and category_cluster are required")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var category, err = repository.CategoriesRepo.GetCategoryContentByFullpath(request.Context(), category_path, category_cluster)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting category content with path '%s' and cluster '%s': %s", category_path, category_cluster, err.Error()))
		http.Error(response, fmt.Sprintf("Error getting category content with path '%s' and cluster '%s'", category_path, category_cluster), 404)
		return
	}

	category_content, err := repository.CategoriesRepo.GetCategoryContent(request.Context(), category.Uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting category content of category<%s>: %s", category.Uuid, err.Error()))
		http.Error(response, fmt.Sprintf("Error getting category content with path '%s' and cluster '%s'", category_path, category_cluster), 404)
		return
	}

	category_content.SortContentSeries()

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(category_content)
}
