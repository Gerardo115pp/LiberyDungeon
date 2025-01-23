package access_sec

import (
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_models "libery-dungeon-libs/models"
	"net/http"
	"time"
)

const cluster_access_prefix string = "cluster-access"

func getClusterAccessCookieName(cluster *dungeon_models.CategoryCluster) string {
	if cluster == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s", cluster_access_prefix, cluster.Uuid)
}

func getClusterAccessTokenExpiration() time.Time {
	return time.Now().Add(time.Hour * 24)
}

func createClusterAccessCookie(cluster_instance *dungeon_models.CategoryCluster) (*http.Cookie, error) {
	if cluster_instance == nil {
		return nil, fmt.Errorf("cluster_instance is nil")
	}

	if !dungeon_secrets.CheckSecretsReady() {
		return nil, dungeon_secrets.DungeonSecretsErr_SecretsUndefined
	}

	expiration_time := getClusterAccessTokenExpiration()

	var access_token string
	access_token, err := dungeon_models.GenerateCategoriesClusterAccess(*cluster_instance, expiration_time, dungeon_secrets.GetDungeonJwtSecret())
	if err != nil {
		return nil, err
	}

	var access_cookie_name string = getClusterAccessCookieName(cluster_instance)

	var access_cookie http.Cookie = http.Cookie{
		Name:     access_cookie_name,
		Value:    access_token,
		Path:     "/",
		Expires:  expiration_time,
		HttpOnly: true,
	}

	return &access_cookie, nil
}

func SignClusterAccessOnRequest(cluster *dungeon_models.CategoryCluster, response http.ResponseWriter) error {
	if cluster == nil {
		return fmt.Errorf("cluster or request is nil")
	}

	var access_cookie *http.Cookie
	access_cookie, err := createClusterAccessCookie(cluster)
	if err != nil {
		return err
	}

	http.SetCookie(response, access_cookie)

	return nil
}
