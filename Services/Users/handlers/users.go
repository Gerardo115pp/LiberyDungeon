package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_users_service/Config"
	"libery_users_service/middlewares"
	service_models "libery_users_service/models"
	"libery_users_service/repository"
	"libery_users_service/workflows"
	"libery_users_service/workflows/common_workflows"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var users_route_path string = "^/users(/.+)?"

var USERS_ROUTE *patriot_router.Route = patriot_router.NewRoute(users_route_path, false)

func UsersHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getUsersHandler
		case http.MethodPost:
			request_handler_func = middlewares.CheckUserCanCreateUsers(postUsersHandler)
		case http.MethodPatch:
			request_handler_func = patchUsersHandler
		case http.MethodDelete:
			request_handler_func = deleteUsersHandler
		case http.MethodPut:
			request_handler_func = putUsersHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getUsersHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/users/identity":
		handler_func = getUserIdentityHandler
	case "/users/read-all":
		handler_func = dungeon_middlewares.CheckUserCan_ReadUsers(getAllUsersHandler)
	case "/users/roles":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(getAllRolesOfUserHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

func getAllRolesOfUserHandler(response http.ResponseWriter, request *http.Request) {
	var username string = request.URL.Query().Get("username")

	if username == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.getAllRolesOfUserHandler, missing username in query")
		response.WriteHeader(400)
		return
	}

	var user *service_models.User

	user, err := repository.UsersRepo.GetUserByUsernameCTX(request.Context(), username)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.getAllRolesOfUserHandler, while getting user<%s>, error: %s", username, err))
		response.WriteHeader(404)
		return
	}

	var user_roles []service_models.RoleTaxonomy

	user_roles, err = repository.UsersRepo.GetUserRolesCTX(request.Context(), user)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.getAllRolesOfUserHandler, while getting roles of user<%s>, error: %s", username, err))
		response.WriteHeader(500)
		return
	}

	var user_role_labels []string = make([]string, 0)

	for _, role := range user_roles {
		user_role_labels = append(user_role_labels, role.RoleLabel)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(user_role_labels)
}

func getAllUsersHandler(response http.ResponseWriter, request *http.Request) {
	var all_users_entires []service_models.UserEntry = make([]service_models.UserEntry, 0)

	all_users, err := repository.UsersRepo.GetAllUsersCTX(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler, while getting all users, error: %s", err))
		response.WriteHeader(500)
		return
	}

	for _, user := range all_users {
		all_users_entires = append(all_users_entires, user.GetAsEntry())
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(all_users_entires)
}

func getUserIdentityHandler(response http.ResponseWriter, request *http.Request) {
	var user_claims *dungeon_models.PlatformUserClaims

	user_claims, err := dungeon_middlewares.GetUserClaims(request, app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler, while getting user claims, error: %s", err))
		response.WriteHeader(401)
		return
	}

	var user_identity *service_models.UserIdentity = &service_models.UserIdentity{
		UUID:          user_claims.UserUUID,
		Username:      user_claims.UserName,
		RoleHierarchy: user_claims.UserHighestHierarchy,
		Grants:        user_claims.UserGrants,
	}

	response.Header().Set("Content-Type", "application/json")

	json.NewEncoder(response).Encode(user_identity)
}

func postUsersHandler(response http.ResponseWriter, request *http.Request) {
	var user_creation_params service_models.CreateUserParams

	var initial_setup_done bool

	initial_setup_done, err := common_workflows.InitialSetupDone(request.Context())
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handles/IsInitialSetupHandler, while checking if initial setup is done, error: %s", err))
		response.WriteHeader(500)
		return
	}

	err = json.NewDecoder(request.Body).Decode(&user_creation_params)
	if err != nil {
		response.WriteHeader(400)
		return
	}

	new_user, err := service_models.CreateNewUser(user_creation_params)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler, while creating new user, error: %s", err))
		response.WriteHeader(500)
		return
	}

	err = repository.UsersRepo.AddUserCTX(request.Context(), new_user)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler, while adding new user, error: %s", err))
		response.WriteHeader(500)
		return
	}

	var user_default_role string = service_models.VISITOR_ROLE

	if !initial_setup_done {
		echo.Echo(echo.YellowBG, "Initial setup is done, setting default role to SUPER_ADMIN")
		user_default_role = service_models.SUPER_ADMIN_ROLE
	}

	err = repository.UsersRepo.AddUserToRoleCTX(request.Context(), new_user, user_default_role)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler, while adding default role to new user, error: %s", err))
		response.WriteHeader(500)
		return
	}

	dungeon_helpers.WriteSingleStringResponseWithStatus(response, new_user.UUID, 201)
}

func patchUsersHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/users/role":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(patchUserRoleHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

func patchUserRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")

	if role_label == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.patchUserRoleHandler, missing role_label in query")
		response.WriteHeader(400)
		return
	}

	_, err := repository.UsersRepo.GetRoleCTX(request.Context(), role_label)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.patchUserRoleHandler, while getting role<%s>, error: %s", role_label, err))
		response.WriteHeader(404)
		return
	}

	var username string = request.URL.Query().Get("username")

	if username == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.patchUserRoleHandler, missing username in query")
		response.WriteHeader(400)
		return
	}

	var user *service_models.User

	user, err = repository.UsersRepo.GetUserByUsernameCTX(request.Context(), username)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.patchUserRoleHandler, while getting user<%s>, error: %s", username, err))
		response.WriteHeader(404)
		return
	}

	err = repository.UsersRepo.AddUserToRoleCTX(request.Context(), user, role_label)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.patchUserRoleHandler, while adding role<%s> to user<%s>, error: %s", role_label, user.UUID, err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func deleteUsersHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/users/user":
		handler_func = dungeon_middlewares.CheckUserCan_DeleteUsers(deleteUserByUUIDHandler)
	case "/users/role":
		handler_func = dungeon_middlewares.CheckUserCan_Grant(deleteUserRoleHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

// We do not ever delete users if they have a super admin role.
func deleteUserByUUIDHandler(response http.ResponseWriter, request *http.Request) {
	var user_uuid string = request.URL.Query().Get("user_uuid")

	if user_uuid == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.deleteUserByUUIDHandler, missing user_uuid in query")
		response.WriteHeader(400)
		return
	}

	var user *service_models.User

	user, err := repository.UsersRepo.GetUserByUuidCTX(request.Context(), user_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, while getting user<%s>, error: %s", user_uuid, err))
		response.WriteHeader(404)
		return
	}

	var user_roles []service_models.RoleTaxonomy

	user_roles, err = repository.UsersRepo.GetUserRolesCTX(request.Context(), user)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, while getting roles of user<%s>, error: %s", user_uuid, err))
		response.WriteHeader(500)
		return
	}

	highest_role_hierarchy := workflows.GetHighestRoleHierarchy(user_roles)

	if highest_role_hierarchy == 0 {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, user<%s> has a super admin role, refusing to delete directly", user_uuid))
		response.WriteHeader(422)
		return
	}

	// Refuse to have a user identity destroy itself.

	user_claims, err := dungeon_middlewares.GetUserClaims(request, app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, while getting user claims, error: %s", err))
		response.WriteHeader(500) // This is likely a server error, an unauthenticated user should not be able to reach this point.
	}

	if user_claims.UserUUID == user_uuid {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, user<%s> is trying to delete itself", user_uuid))
		response.WriteHeader(422)
		return
	}

	err = repository.UsersRepo.DeleteUserByUuidCTX(request.Context(), user_uuid)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserByUUIDHandler, while deleting user<%s>, error: %s", user_uuid, err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func deleteUserRoleHandler(response http.ResponseWriter, request *http.Request) {
	var role_label string = request.URL.Query().Get("role_label")

	if role_label == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.deleteUserRoleHandler, missing role_label in query")
		response.WriteHeader(400)
		return
	}

	var username string = request.URL.Query().Get("username")

	if username == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.deleteUserRoleHandler, missing username in query")
		response.WriteHeader(400)
		return
	}

	var user *service_models.User

	user, err := repository.UsersRepo.GetUserByUsernameCTX(request.Context(), username)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserRoleHandler, while getting user<%s>, error: %s", username, err))
		response.WriteHeader(404)
		return
	}

	err = repository.UsersRepo.DeleteUserFromRoleCTX(request.Context(), user, role_label)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.deleteUserRoleHandler, while deleting role<%s> from user<%s>, error: %s", role_label, user.UUID, err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func putUsersHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc

	switch resource {
	case "/users/user/secret":
		handler_func = dungeon_middlewares.CheckUserCan_ModifyUsers(putUserSecretHandler)
	case "/users/user/username":
		handler_func = dungeon_middlewares.CheckUserCan_ModifyUsers(putUserUsernameHandler)
	default:
		handler_func = dungeon_helpers.ResourceNotFoundHandler
	}

	handler_func(response, request)
}

// Changes the secret of a user. The repository method accepts a user data and only requires the uuid to match a user entry on the database.
// So this theoretically could change a user's username as well, this handles should only change the secret so if the username does not match, we should
// return 404 error immediately.
func putUserSecretHandler(response http.ResponseWriter, request *http.Request) {
	var new_user_data service_models.User

	err := json.NewDecoder(request.Body).Decode(&new_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserSecretHandler, while decoding new user data, error: %s", err))
		response.WriteHeader(400)
		return
	}

	if new_user_data.UUID == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.putUserSecretHandler, missing user_uuid in new user data")
		response.WriteHeader(404)
		return
	}

	var current_user_data *service_models.User

	current_user_data, err = repository.UsersRepo.GetUserByUuidCTX(request.Context(), new_user_data.UUID)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserSecretHandler, while getting user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(404)
		return
	}

	if current_user_data.Username != new_user_data.Username {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserSecretHandler, user<%s> does not match new user data", new_user_data.UUID))
		response.WriteHeader(421)
		return
	}

	err = current_user_data.UpdateFromOther(&new_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserSecretHandler, while updating user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(500)
		return
	}

	err = repository.UsersRepo.UpdateUserCTX(request.Context(), current_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserSecretHandler, while updating user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}

// Changes the user's username. As the putUserSecretHandler, this should only change the username, so if the secret does not match
// the request shall be refused.
func putUserUsernameHandler(response http.ResponseWriter, request *http.Request) {
	var new_user_data service_models.User

	err := json.NewDecoder(request.Body).Decode(&new_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserUsernameHandler, while decoding new user data, error: %s", err))
		response.WriteHeader(400)
		return
	}

	if new_user_data.UUID == "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.putUserUsernameHandler, missing user_uuid in new user data")
		response.WriteHeader(404)
		return
	}

	// if the secret is empty, user.UpdateFromOther will leave the current secret as is.
	if new_user_data.SecretHash != "" {
		echo.Echo(echo.RedFG, "In Handlers/UsersHandler.putUserUsernameHandler, secret hash is not empty, this handler only changes the username")
		response.WriteHeader(421)
		return
	}

	var current_user_data *service_models.User

	current_user_data, err = repository.UsersRepo.GetUserByUuidCTX(request.Context(), new_user_data.UUID)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserUsernameHandler, while getting user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(404)
		return
	}

	err = current_user_data.UpdateFromOther(&new_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserUsernameHandler, while updating user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(500)
		return
	}

	err = repository.UsersRepo.UpdateUserCTX(request.Context(), current_user_data)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/UsersHandler.putUserUsernameHandler, while updating user<%s>, error: %s", new_user_data.UUID, err))
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(204)
}
