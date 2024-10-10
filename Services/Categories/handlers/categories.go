package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CategoriesHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getCategoriesHandler(response, request)
		case http.MethodPost:
			postCategoriesHandler(response, request)
		case http.MethodPatch:
			patchCategoriesHandler(response, request)
		case http.MethodDelete:
			deleteCategoriesHandler(response, request)
		case http.MethodPut:
			putCategoriesHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getCategoriesHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	if resource_path == "/categories/name-available" {
		getCategoriesNameAvailableHandler(response, request)
		return
	} else if resource_path == "/categories/data" {
		getCategoriesDataHandler(response, request)
	} else if resource_path == "/categories/by-fullpath" {
		getCategoriesByFullpathHandler(response, request)
	} else {
		response.WriteHeader(404)
		echo.Echo(echo.RedBG, fmt.Sprintf("Resource not found: %s", resource_path))
	}

	return
}

func getCategoriesDataHandler(response http.ResponseWriter, request *http.Request) {
	var category_id string = request.URL.Query().Get("category_id")

	if category_id == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_id parameter"))
		response.WriteHeader(400)
		return
	}

	category, err := repository.CategoriesRepo.GetCategory(request.Context(), category_id)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting category: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(category)
}

func getCategoriesByFullpathHandler(response http.ResponseWriter, request *http.Request) {
	var category_fullpath string = request.URL.Query().Get("category_path")
	var category_cluster string = request.URL.Query().Get("category_cluster")

	if category_fullpath == "" || category_cluster == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_path or category_cluster parameter"))
		response.WriteHeader(400)
		return
	}

	category, err := repository.CategoriesRepo.GetCategoryContentByFullpath(request.Context(), category_fullpath, category_cluster)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting category by fullpath<%s> and cluster<%s>: %s", category_fullpath, category_cluster, err.Error()))
		response.WriteHeader(404)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(category)
}

func getCategoriesNameAvailableHandler(response http.ResponseWriter, request *http.Request) {
	var category_name string = request.URL.Query().Get("category_name")
	var parent_id string = request.URL.Query().Get("parent_id")

	if category_name == "" || parent_id == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_name or parent_id parameter"))
		response.WriteHeader(400)
		return
	}
	name_available, err := workflows.IsCategoryNameAvailable(category_name, parent_id, false)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error checking if category name is available: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	var status_code int = 200

	if !name_available {
		status_code = 409 // Conflict
	}

	response.WriteHeader(status_code)
}

func postCategoriesHandler(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.GreenFG, fmt.Sprintf("POST /categories, creating new category"))
	new_category_request := &struct {
		Name    string `json:"name"`
		Parent  string `json:"parent"`
		Cluster string `json:"cluster"`
	}{}

	err := json.NewDecoder(request.Body).Decode(new_category_request)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error decoding new category request: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	if new_category_request.Name == "" || new_category_request.Parent == "" || new_category_request.Cluster == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing name, fullpath or parent parameter"))
		response.WriteHeader(400)
		return
	}

	var new_category dungeon_models.Category

	new_category, err = workflows.CreateNewCategory(request.Context(), new_category_request.Name, new_category_request.Parent, new_category_request.Cluster)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error creating new category: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(201)
	response.Header().Add("Content-Type", "application/json")
	json.NewEncoder(response).Encode(new_category)
	return
}

func patchCategoriesHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path

	if resource_path == "/categories/rename" {
		patchCategoryRenameHandler(response, request)
	} else {
		patchCategoryContentHandler(response, request)
	}

	return
}

func patchCategoryRenameHandler(response http.ResponseWriter, request *http.Request) {
	var category_rename_request = &struct {
		CategoryID string `json:"category_id"`
		NewName    string `json:"new_name"`
	}{}

	err := json.NewDecoder(request.Body).Decode(category_rename_request)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error decoding category rename request: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	if category_rename_request.CategoryID == "" || category_rename_request.NewName == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_id or new_name parameter"))
		response.WriteHeader(400)
		return
	}

	err = workflows.RenameCategory(category_rename_request.CategoryID, category_rename_request.NewName)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error renaming category: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	modified_category, err := repository.CategoriesRepo.GetCategory(request.Context(), category_rename_request.CategoryID)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting modified category: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(200)
	json.NewEncoder(response).Encode(modified_category)
}

func patchCategoryContentHandler(response http.ResponseWriter, request *http.Request) {
	var rejected_medias []dungeon_models.Media
	var moved_medias map[string][]dungeon_models.Media

	category_changes := &struct {
		RejectedMedias []dungeon_models.Media            `json:"rejected_medias"`
		MovedMedias    map[string][]dungeon_models.Media `json:"moved_medias"`
	}{}

	var category_id string = request.URL.Query().Get("category_id")
	if category_id == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_id parameter"))
		response.WriteHeader(400)
		return
	}

	err := json.NewDecoder(request.Body).Decode(category_changes)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error decoding category changes: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	rejected_medias = category_changes.RejectedMedias
	moved_medias = category_changes.MovedMedias

	current_category, err := repository.CategoriesRepo.GetCategory(request.Context(), category_id)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting category: %s", err.Error()))
		response.WriteHeader(404)
		return
	}

	// Get the category cluster
	medias_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), current_category.Cluster)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error getting category cluster: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	if len(rejected_medias) > 0 {
		err = workflows.ProcessRejectedMedias(rejected_medias, current_category, &medias_cluster)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error processing rejected medias: %s", err.Error()))
			response.WriteHeader(500)
			return
		}
	}

	if len(moved_medias) > 0 {
		err = workflows.ProcessMovedMedias(moved_medias, current_category, &medias_cluster)
		if err != nil {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Error processing moved medias: %s", err.Error()))
			response.WriteHeader(500)
			return
		}
	}

	// Emit platform event

	category_changes_event := communication.NewClusterFSChangeEvent(app_config.JWT_SECRET, medias_cluster.Uuid, len(category_changes.RejectedMedias), 0, len(category_changes.MovedMedias))

	err = category_changes_event.Emit()
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error emitting category changes event: %s", err.Error()))
	}

	response.WriteHeader(200)
}

func deleteCategoriesHandler(response http.ResponseWriter, request *http.Request) {
	category_delete_request := &struct {
		CategoryID string `json:"category_id"`
		Force      bool   `json:"force"`
	}{}

	// If force is true then delete the category and all its childs and medias
	// if not, then we only delete the category if it has no childs	or medias

	err := json.NewDecoder(request.Body).Decode(category_delete_request)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error decoding category delete request: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	if category_delete_request.CategoryID == "" {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Missing category_id parameter"))
		response.WriteHeader(400)
		return
	}

	// Check if the category is empty
	empty, err := repository.CategoriesRepo.IsCategoryEmpty(request.Context(), category_delete_request.CategoryID)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error checking if category is empty: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	if !empty && !category_delete_request.Force {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Category is not empty"))
		response.WriteHeader(409) // Conflict
	}

	deleted_category, err := workflows.DeleteCategory(category_delete_request.CategoryID, category_delete_request.Force)
	if err != nil {
		echo.Echo(echo.YellowFG, fmt.Sprintf("Error deleting category: %s", err.Error()))
		response.WriteHeader(500)
		return
	}

	var status_code int = 205 // 205: Reset Content

	if !deleted_category {
		status_code = 304 // 304: Not Modified
	}

	response.WriteHeader(status_code)
	return
}

func putCategoriesHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
