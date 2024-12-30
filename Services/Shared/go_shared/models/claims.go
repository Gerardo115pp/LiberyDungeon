package dungeon_models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtSigningMethod *jwt.SigningMethodHMAC = jwt.SigningMethodHS256

type ClaimsStandardResponse struct {
	RedirectURL string `json:"redirect_url"`
	Granted     bool   `json:"granted"`
}

type CategoriesClusterAccessGrant struct {
	jwt.StandardClaims
	CategoryCluster
}

func GenerateCategoriesClusterAccess(categories_cluster CategoryCluster, expires_at time.Time, sk string) (string, error) {
	claims := &CategoriesClusterAccessGrant{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at.Unix(),
		},
		CategoryCluster: categories_cluster,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString([]byte(sk))
}

func ParseCategoriesClusterAccess(token string, sk string) (*CategoriesClusterAccessGrant, error) {
	claims := &CategoriesClusterAccessGrant{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return claims, err
}

func GetCategoriesClusterFromToken(token string, sk string) (*CategoryCluster, error) {
	claims, err := ParseCategoriesClusterAccess(token, sk)
	if err != nil {
		return nil, err
	}

	return &claims.CategoryCluster, nil
}

type PlatformUserClaims struct {
	jwt.StandardClaims
	UserUUID             string   `json:"user_uuid"`
	UserName             string   `json:"username"`
	UserHighestHierarchy int      `json:"user_highest_hierarchy"` // Lower values mean higher hierarchy
	UserGrants           []string `json:"user_grants"`
}

func GeneratePlatformUserClaims(user_uuid, username string, user_highest_hierarchy int, user_grants []string, expires_at time.Time, sk string) (string, error) {
	claims := &PlatformUserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at.Unix(),
		},
		UserUUID:             user_uuid,
		UserName:             username,
		UserHighestHierarchy: user_highest_hierarchy,
		UserGrants:           user_grants,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString([]byte(sk))
}

func ParsePlatformUserClaims(token string, sk string) (*PlatformUserClaims, error) {
	claims := &PlatformUserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return claims, err
}

// =========================
// RESOURCE SHARING
// =========================

type MediaShareToken struct {
	jwt.StandardClaims
	MediaIdentity *MediaIdentity `json:"media_identity"`
}

func GenerateMediaShareToken(media_identity *MediaIdentity, expires_at time.Time, sk string) (string, error) {
	claims := &MediaShareToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at.Unix(),
		},
		MediaIdentity: media_identity,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString([]byte(sk))
}

func ParseMediaShareToken(token string, sk string) (*MediaShareToken, error) {
	claims := &MediaShareToken{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return claims, err
}
