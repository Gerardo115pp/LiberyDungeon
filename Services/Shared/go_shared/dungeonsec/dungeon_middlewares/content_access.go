package dungeon_middlewares

import (
	"libery-dungeon-libs/dungeonsec"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"net/http"
)

func CheckUserCan_ShareContent(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	if !checkMiddlewareCanRun() {
		panic("You missed one step of the middleware setup, steps are: set the jwt secret")
	}

	return func(response http.ResponseWriter, request *http.Request) {
		user_claims, err := GetUserClaims(request, dungeon_secrets.GetDungeonJwtSecret())
		if err != nil {
			response.WriteHeader(401)
			return
		}

		var the_user_can bool = false

		the_user_can = dungeonsec.CanShareContent(user_claims.UserGrants)

		if the_user_can {
			next(response, request)
			return
		}

		response.WriteHeader(403)
	}
}
