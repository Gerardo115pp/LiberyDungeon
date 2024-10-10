package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func SearchHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getSearchHandler(response, request)
		case http.MethodPost:
			postSearchHandler(response, request)
		case http.MethodPatch:
			patchSearchHandler(response, request)
		case http.MethodDelete:
			deleteSearchHandler(response, request)
		case http.MethodPut:
			putSearchHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getSearchHandler(response http.ResponseWriter, request *http.Request) {
	var query string = request.URL.Query().Get("query")
	var ignore string = request.URL.Query().Get("ignore")
	var cluster_id string = request.URL.Query().Get("cluster_id")

	if cluster_id == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing cluster_id parameter"))
		response.WriteHeader(400)
		return
	}

	if query == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing query parameter"))
		response.WriteHeader(400)
		return
	}

	all_categories, err := repository.CategoriesRepo.GetClusterCategories(request.Context(), cluster_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error getting all categories: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	// Exclude categories such that category.uuid == ignore
	if ignore != "" {
		filtered_categories := make([]dungeon_models.Category, 0)
		for _, category := range all_categories {
			if category.Parent != ignore {
				filtered_categories = append(filtered_categories, category)
			}
		}

		all_categories = filtered_categories
	}

	var matches []*dungeon_models.Category = make([]*dungeon_models.Category, 0)

	matches = workflows.GetUniqueMatches(query, all_categories, 0.85)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	err = json.NewEncoder(response).Encode(matches)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error encoding response: %s", err.Error()))
		return
	}
}

func postSearchHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func patchSearchHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteSearchHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putSearchHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
