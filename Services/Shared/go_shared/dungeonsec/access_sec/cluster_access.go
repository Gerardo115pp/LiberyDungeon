package access_sec

import (
	"errors"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_models "libery-dungeon-libs/models"
	"net/http"
	"time"
)

const cluster_access_prefix string = "cluster-access"

var (
	MissingClusterAccessErr error = errors.New("Request has no access to the given cluster uuid")
)

// Creates a new cluster_access_token cookie for a given cluster
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

// Returns the name a cluster_access_token cookie for a given cluster must have
func getClusterAccessCookieName(cluster *dungeon_models.CategoryCluster) string {
	if cluster == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s", cluster_access_prefix, cluster.Uuid)
}

// Returns a cluster access cookie for the given cluster id in the given request.
func getClusterAccessCookie(cluster_uuid string, request *http.Request) *http.Cookie {
	if cluster_uuid == "" || request == nil {
		return nil
	}

	access_cookie_name := getClusterAccessCookieName(&dungeon_models.CategoryCluster{Uuid: cluster_uuid})

	access_cookie, err := request.Cookie(access_cookie_name)
	if err != nil {
		return nil
	}

	return access_cookie
}

// Returns the time at which a cluster_access_token should expire from the time of invocation
func getClusterAccessTokenExpiration() time.Time {
	return time.Now().Add(time.Hour * 24)
}

/** Returns cluster object for the given cluster_id that was signed on a given request. Or nil if:
 * The request had no matching cookie
 * The token on the cookie was expired
 * The signature on the token as invalid
 */
func GetSignedClusterOnRequest(cluster_uuid string, request *http.Request) (*dungeon_models.CategoryCluster, error) {
	if cluster_uuid == "" || request == nil {
		return nil, fmt.Errorf("cluster_uuid or request is nil")
	}

	cluster_access_cookie := getClusterAccessCookie(cluster_uuid, request)
	if cluster_access_cookie == nil {
		return nil, MissingClusterAccessErr
	}

	category_cluster, err := dungeon_models.GetCategoriesClusterFromToken(cluster_access_cookie.Value, dungeon_secrets.GetDungeonJwtSecret())
	if err != nil {
		return nil, errors.Join(err, MissingClusterAccessErr)
	}

	return category_cluster, nil
}

// Returns whether the given request has cluster_access_token matching the cluster_id.
func RequestHasClusterAccess(cluster_id string, request *http.Request) bool {
	if cluster_id == "" || request == nil {
		return false
	}

	access_cookie := getClusterAccessCookie(cluster_id, request)
	if access_cookie == nil {
		return false
	}

	return true
}

// Adds a brand new signed cluster_access_token for the given cluster as a cookie to the given response writer.
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
