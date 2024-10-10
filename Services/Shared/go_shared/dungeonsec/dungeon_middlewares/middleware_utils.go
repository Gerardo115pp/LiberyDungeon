package dungeon_middlewares

import (
	dungeon_models "libery-dungeon-libs/models"
	"net/http"
)

const USER_CLAIMS_COOKIE_NAME string = "auth_user_claims"

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

func GetUserClaims(request *http.Request, sk string) (*dungeon_models.PlatformUserClaims, error) {
	user_cookie, err := request.Cookie(USER_CLAIMS_COOKIE_NAME)
	if err != nil {
		return nil, err
	}

	user_claims, err := dungeon_models.ParsePlatformUserClaims(user_cookie.Value, sk)
	if err != nil {
		return nil, err
	}

	// TODO: Add a check to see if the user exists. this is a distributed system, so if a user did not changed the default jwt secret, a malicious user could create a jwt token
	// on their system that has full access and send it to another user's system(that uses the default jwt secret) and get full access to that system. This another good argument
	// to force users to set a jwt secret and not have a default one at all.

	return user_claims, nil
}
