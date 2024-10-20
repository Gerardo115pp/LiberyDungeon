package dungeon_middlewares

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"net/http"
)

func factory_grantOnClaimCheckMiddleware(checker_func dungeonsec.UserCanChecker) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if !checkMiddlewareCanRun() {
			panic("You missed one step of the middleware setup, steps are: set the jwt secret")
		}

		return func(response http.ResponseWriter, request *http.Request) {
			var is_authenticated bool = false

			user_claims, err := GetUserClaims(request, dungeon_secrets.GetDungeonJwtSecret())
			if err == nil {
				is_authenticated = true
			} else {
				fmt.Printf("Failed to get user claims: %s\nAssuming the user is not authenticated. Will attempt to validate call with domain secret.\n", err)
			}

			if !is_authenticated {
				domain_secret_wrapper := CheckDomainSecretMiddleware(next)
				domain_secret_wrapper(response, request)
				return
			}

			var the_user_can bool = false

			the_user_can = checker_func(user_claims.UserGrants)

			if the_user_can {
				next(response, request)
				return
			}

			response.WriteHeader(403)
		}
	}
}

var checkUserCan_Grant MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanGrant)

func CheckUserCan_Grant(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_Grant(next)
}

var checkUserCan_ReadUsers MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanReadUsers)

func CheckUserCan_ReadUsers(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_ReadUsers(next)
}

var checkUserCan_ModifyUsers MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanModifyUsers)

func CheckUserCan_ModifyUsers(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_ModifyUsers(next)
}

var checkUserCan_DeleteUsers MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanDeleteUsers)

func CheckUserCan_DeleteUsers(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_DeleteUsers(next)
}

var checkUserCan_UploadFiles MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanUploadFiles)

func CheckUserCan_UploadFiles(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_UploadFiles(next)
}

var checkUserCan_ContentAlter MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanContentAlter)

func CheckUserCan_ContentAlter(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_ContentAlter(next)
}

var checkUserCan_DungeonTagsCreate MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanDungeonTagsCreate)

func CheckUserCan_DungeonTagsCreate(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_DungeonTagsCreate(next)
}

var checkUserCan_DungeonTagsTag MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanDungeonTagsTag)

func CheckUserCan_DungeonTagsTag(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_DungeonTagsTag(next)
}

var checkUserCan_DungeonTagsUntag MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanDungeonTagsUntag)

func CheckUserCan_DungeonTagsUntag(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_DungeonTagsUntag(next)
}

var checkUserCan_DungeonTagsTaxonomyCreate MiddlewareFunc = factory_grantOnClaimCheckMiddleware(dungeonsec.CanDungeonTagsTaxonomyCreate)

func CheckUserCan_DungeonTagsTaxonomyCreate(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return checkUserCan_DungeonTagsTaxonomyCreate(next)
}
