package resource_access

import (
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/dungeonsec"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"net/http"
)

func IsClusterPrivate(cluster_uuid string) (bool, error) {
	var cluster_is_private bool

	cluster_is_private, err := communication.Metadata.CheckClusterPrivate(cluster_uuid)
	if err != nil {
		return false, err
	}

	return cluster_is_private, nil
}

func UserClaimsCookieHasClusterAccess(cluster_uuid string, request *http.Request) (bool, error) {
	is_cluster_private, err := IsClusterPrivate(cluster_uuid)
	if err != nil {
		return false, err
	}

	var user_has_access bool = !is_cluster_private

	if !user_has_access {
		user_claims, err := dungeon_middlewares.GetUserClaims(request, dungeon_secrets.GetDungeonJwtSecret())
		if err != nil {
			return false, err
		}

		user_has_access = dungeonsec.CanViewPrivateClusters(user_claims.UserGrants)
	}

	return user_has_access, nil
}
