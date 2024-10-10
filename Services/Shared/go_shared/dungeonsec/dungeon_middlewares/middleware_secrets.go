package dungeon_middlewares

import (
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"net/http"
)

func CheckDomainSecretMiddleware(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// TODO: Change this to an efficient(and by that i mean quick) signature check instead of directly passing the domain secret.
		var domain_secret string = dungeon_secrets.GetDungeonDomainSecret()
		if request.Header.Get(dungeon_secrets.DOMAIN_SECRET_HEADER) != domain_secret {
			response.WriteHeader(401)
			return
		}

		next(response, request)
	}
}

func checkMiddlewareCanRun() bool {
	var can_run bool = false

	if dungeon_secrets.GetDungeonJwtSecret() != "" {
		can_run = true
	}

	if dungeon_secrets.GetDungeonDomainSecret() != "" {
		can_run = true
	}

	return can_run
}
