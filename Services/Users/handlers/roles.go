package handlers

import (
	"encoding/json"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	service_models "libery_users_service/models"
	"libery_users_service/repository"
	"net/http"
	"slices"
	"strconv"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var user_roles_route_path string = "^/roles(/.+)?"

var USER_ROLES_ROUTE *patriot_router.Route = patriot_router.NewRoute(user_roles_route_path, false)

func UserRolesHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getUserRolesHandler
		case http.MethodPost:
			request_handler_func = postUserRolesHandler
		case http.MethodPatch:
			request_handler_func = patchUserRolesHandler
		case http.MethodDelete:
			request_handler_func = deleteUserRolesHandler
		case http.MethodPut:
			request_handler_func = putUserRolesHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getUserRolesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/roles/read-all":
		handler_func = dungeon_middlewares.CheckUserCan_ReadUsers(getAllRolesHandler)
	case "/roles/grant/read-all":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(getAllGrantsHandler)
	case "/roles/below-hierarchy":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(getRoleBelowHierarchyHandler)
	case "/roles/role":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(getUserRoleHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

func getUserRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")

	if role_label == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getUserRoleHandler, missing role_label in query")
		response.WriteHeader(400)
		return
	}

	var role service_models.RoleTaxonomy
	var err error

	role, err = repository.UsersRepo.GetRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getUserRoleHandler, while getting role<%s>, error: %s", role_label, err)
		response.WriteHeader(404)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(role)
}

func getAllGrantsHandler(response http.ResponseWriter, request *http.Request) {
	var all_grants []string = make([]string, 0)

	all_grants, err := repository.UsersRepo.GetAllGrantsCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getAllGrantsHandler, while getting all grants, error: %s", err)
		response.WriteHeader(500)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(all_grants)
}

func getAllRolesHandler(response http.ResponseWriter, request *http.Request) {
	var all_roles []service_models.RoleTaxonomy
	var all_roles_labels []string = make([]string, 0)

	all_roles, err := repository.UsersRepo.GetAllRolesCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getAllRolesHandler, while getting all roles, error: %s", err)
		response.WriteHeader(500)
		return
	}

	for _, role := range all_roles {
		all_roles_labels = append(all_roles_labels, role.RoleLabel)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(all_roles_labels)
}

func getRoleBelowHierarchyHandler(response http.ResponseWriter, request *http.Request) {
	var roles_below_hierarchy []service_models.RoleTaxonomy = make([]service_models.RoleTaxonomy, 0)
	var role_hierarchy_string string = request.URL.Query().Get("hierarchy")
	var role_hierarchy int
	var err error

	if role_hierarchy_string == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getRoleBelowHierarchyHandler, missing hierarchy in query")
		response.WriteHeader(400)
		return
	}

	role_hierarchy, err = strconv.Atoi(role_hierarchy_string)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getRoleBelowHierarchyHandler, while converting hierarchy<%s> to int, error: %s", role_hierarchy_string, err)
		response.WriteHeader(400)
		return
	}

	roles_below_hierarchy, err = repository.UsersRepo.FindRolesBelowHierarchyCTX(request.Context(), role_hierarchy)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.getRoleBelowHierarchyHandler, while finding roles below hierarchy<%d>, error: %s", role_hierarchy, err)
		response.WriteHeader(500)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(roles_below_hierarchy)
}

func postUserRolesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/roles/grant":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(postCreateNewGrantHandler)
	case "/roles/add-grant":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(postAddGrantToRoleHandler)
	case "/roles/role":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(postCreateNewRoleHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler

	}

	handler_func(response, request)
}

func postCreateNewGrantHandler(response http.ResponseWriter, request *http.Request) {
	var new_grant string = request.URL.Query().Get("new_grant")

	if new_grant == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postCreateNewGrantHandler, missing new_grant in query")
		response.WriteHeader(400)
		return
	}

	err := repository.UsersRepo.AddGrantCTX(request.Context(), new_grant)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postCreateNewGrantHandler, while adding new grant, error: %s", err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(201)
}

func postAddGrantToRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")
	var grant string = request.URL.Query().Get("grant")

	if role_label == "" || grant == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postAddGrantToRoleHandler, missing role_label or grant in query")
		response.WriteHeader(400)
		return
	}

	_, err := repository.UsersRepo.GetRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postAddGrantToRoleHandler, while getting role<%s>, error: %s", role_label, err)
		response.WriteHeader(404)
		return
	}

	all_grants, err := repository.UsersRepo.GetAllGrantsCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postAddGrantToRoleHandler, while getting all grants, error: %s", err)
		response.WriteHeader(500)
		return
	}

	var grant_exists bool = slices.Contains(all_grants, grant)
	if !grant_exists {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postAddGrantToRoleHandler, grant<%s> does not exist", grant)
		response.WriteHeader(400)
		return
	}

	err = repository.UsersRepo.AddGrantToRoleCTX(request.Context(), role_label, grant, true)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postAddGrantToRoleHandler, while adding grant<%s> to role<%s>, error: %s", grant, role_label, err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(201)
}

func postCreateNewRoleHandler(response http.ResponseWriter, request *http.Request) {
	var new_role service_models.RoleTaxonomy

	err := json.NewDecoder(request.Body).Decode(&new_role)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postCreateNewRoleHandler, while decoding new role, error: %s", err)
		response.WriteHeader(400)
		return
	}

	err = repository.UsersRepo.CreateNewRoleCTX(request.Context(), new_role)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.postCreateNewRoleHandler, while creating new role, error: %s", err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(201)
}

func patchUserRolesHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteUserRolesHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/roles/grant":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(deleteGrantHandler)
	case "/roles/role":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(deleteRoleHandler)
	case "/roles/remove-grant":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(deleteGrantFromRoleHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

func deleteGrantHandler(response http.ResponseWriter, request *http.Request) {
	var grant string = request.URL.Query().Get("grant")

	if grant == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantHandler, missing grant in query")
		response.WriteHeader(400)
		return
	}

	err := repository.UsersRepo.DeleteGrantCTX(request.Context(), grant)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantHandler, while deleting grant<%s>, error: %s", grant, err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func deleteRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")

	if role_label == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteRoleHandler, missing role_label in query")
		response.WriteHeader(400)
		return
	}

	role, err := repository.UsersRepo.GetRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteRoleHandler, while getting role<%s>, error: %s", role_label, err)
		response.WriteHeader(404)
		return
	}

	if role.RoleHierarchy == 0 {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteRoleHandler, role<%s> is a root role", role_label)
		response.WriteHeader(422)
		return
	}

	err = repository.UsersRepo.DeleteRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteRoleHandler, while deleting role<%s>, error: %s", role_label, err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func deleteGrantFromRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")
	var grant string = request.URL.Query().Get("grant")

	if role_label == "" || grant == "" {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantFromRoleHandler, missing role_label or grant in query")
		response.WriteHeader(400)
		return
	}

	_, err := repository.UsersRepo.GetRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantFromRoleHandler, while getting role<%s>, error: %s", role_label, err)
		response.WriteHeader(400) // obscure what roles actually exist
		return
	}

	all_grants, err := repository.UsersRepo.GetAllGrantsCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantFromRoleHandler, while getting all grants, error: %s", err)
		response.WriteHeader(500)
		return
	}

	var grant_exists bool = slices.Contains(all_grants, grant)
	if !grant_exists {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantFromRoleHandler, grant<%s> does not exist", grant)
		response.WriteHeader(400)
		return
	}

	err = repository.UsersRepo.DeleteGrantFromRoleCTX(request.Context(), role_label, grant)
	if err != nil {
		echo.Echo(echo.RedFG, "In Handlers/UserRolesHandler.deleteGrantFromRoleHandler, while deleting grant<%s> from role<%s>, error: %s", grant, role_label, err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

func putUserRolesHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
