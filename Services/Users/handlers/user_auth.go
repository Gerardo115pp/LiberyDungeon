package handlers

import (
	"encoding/json"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_users_service/Config"
	service_models "libery_users_service/models"
	"libery_users_service/repository"
	"libery_users_service/workflows"
	"net/http"
	"time"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var user_auth_route_path string = "/user-auth(/.+)?"

var USER_AUTH_ROUTE *patriot_router.Route = patriot_router.NewRoute(user_auth_route_path, false)

func UserAuthHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getUserAuthHandler(response, request)
		case http.MethodPost:
			postUserAuthHandler(response, request)
		case http.MethodPatch:
			patchUserAuthHandler(response, request)
		case http.MethodDelete:
			deleteUserAuthHandler(response, request)
		case http.MethodPut:
			putUserAuthHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getUserAuthHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/user-auth":
		echo.Echo(echo.SkyBlueFG, "Requesting user login")
		getLoginUserHandler(response, request)
	case "/user-auth/verify":
		echo.Echo(echo.SkyBlueFG, "Requesting user auth token validation")
		getVerifyUserHandler(response, request)
	case "/user-auth/logout":
		echo.Echo(echo.SkyBlueFG, "Requesting user logout")
		getLogoutUserHandler(response, request)
	default:
		response.WriteHeader(404)
	}
}

func getLoginUserHandler(response http.ResponseWriter, request *http.Request) {
	var err error
	var auth_username string = request.URL.Query().Get("username")
	var auth_secret string = request.URL.Query().Get("secret")

	if auth_username == "" || auth_secret == "" {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getUserAuthHandler, missing username<%s> or secret<%s> in query", auth_username, auth_secret))
		response.WriteHeader(400)
		return
	}

	access_response := &struct {
		Granted  bool                         `json:"granted"`
		UserData *service_models.UserIdentity `json:"user_data"`
	}{
		Granted:  false,
		UserData: nil,
	}

	user_credentials, err := repository.UsersRepo.GetUserByUsernameCTX(request.Context(), auth_username)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getUserAuthHandler, while getting user<%s> by username, error: %s. Assuming not found", auth_username, err))
		response.WriteHeader(404)
		return
	}

	err = user_credentials.CompareSecret(auth_secret)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getUserAuthHandler, while comparing user<%s> secret, error: %s", user_credentials.Username, err))
		response.WriteHeader(200)
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(access_response)
		return
	}

	user_roles, err := repository.UsersRepo.GetUserRolesCTX(request.Context(), user_credentials)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getUserAuthHandler, while getting user roles for user<%s>, error: %s", user_credentials.Username, err))
		response.WriteHeader(500)
		return
	}

	var user_highest_role_hierarchy int = workflows.GetHighestRoleHierarchy(user_roles)

	var user_grants []string = workflows.CompileUserGrants(user_roles)

	claims_expiration_time := time.Now().Add(app_config.USER_CLAIMS_EXPIRATION_HOURS)

	user_token, err := dungeon_models.GeneratePlatformUserClaims(user_credentials.UUID, user_credentials.Username, user_highest_role_hierarchy, user_grants, claims_expiration_time, app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getUserAuthHandler, while generating user token for user<%s>, error: %s", user_credentials.Username, err))
		response.WriteHeader(500)
		return
	}

	var user_claims_cookie http.Cookie = http.Cookie{
		Name:     app_config.USER_CLAIMS_COOKIE_NAME,
		Value:    user_token,
		Path:     "/",
		Expires:  claims_expiration_time,
		HttpOnly: true,
	}

	http.SetCookie(response, &user_claims_cookie)

	access_response.Granted = true

	access_response.UserData = service_models.CreateUserIdentity(user_credentials, user_highest_role_hierarchy, user_grants)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(access_response)
}

func getVerifyUserHandler(response http.ResponseWriter, request *http.Request) {
	var user_claims_cookie *http.Cookie
	var user_claims_token string

	user_claims_cookie, err := request.Cookie(app_config.USER_CLAIMS_COOKIE_NAME)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In Handlers/getVerifyUserHandler, while getting user claims cookie, error: %s", err))
		dungeon_helpers.WriteBooleanResponse(response, false)
		return
	}

	user_claims_token = user_claims_cookie.Value
	_, err = dungeon_models.ParsePlatformUserClaims(user_claims_token, app_config.JWT_SECRET)

	dungeon_helpers.WriteBooleanResponse(response, err == nil)
}

func getLogoutUserHandler(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.SkyBlueFG, "Logging out user")

	var user_claims_cookie *http.Cookie = &http.Cookie{
		Name:     app_config.USER_CLAIMS_COOKIE_NAME,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(response, user_claims_cookie)

	response.WriteHeader(200)
}

func postUserAuthHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchUserAuthHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteUserAuthHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putUserAuthHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
