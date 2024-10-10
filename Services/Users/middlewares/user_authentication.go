package middlewares

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	app_config "libery_users_service/Config"
	"libery_users_service/repository"
	"libery_users_service/workflows/common_workflows"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CheckUserCanCreateUsers(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		initial_setup_secret := request.URL.Query().Get("initial-setup-secret")
		if app_config.INITIAL_SETUP_SECRET != "" && initial_setup_secret == app_config.INITIAL_SETUP_SECRET {
			is_initial_setup_done, err := common_workflows.InitialSetupDone(request.Context())
			if err != nil {
				response.WriteHeader(500)
				return
			}

			if is_initial_setup_done {
				echo.EchoWarn("Initial setup secret provided, but initial setup is already done")
				response.WriteHeader(403)
				return
			}

			next(response, request)
			return
		}

		user_claims, err := dungeon_middlewares.GetUserClaims(request, app_config.JWT_SECRET)
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("Error getting user claims: %s", err))
			response.WriteHeader(401)
			return
		}

		_, err = repository.UsersRepo.GetUserByUuidCTX(request.Context(), user_claims.UserUUID) // Check if user exists
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("Error getting user by uuid: %s", err))
			response.WriteHeader(401)
			return
		}

		var authorized bool = dungeonsec.CanCreateUsers(user_claims.UserGrants)

		if !authorized {
			echo.Echo(echo.RedBG, "User is not authorized to create users")
			response.WriteHeader(403)
			return
		}

		next(response, request)
	}
}
