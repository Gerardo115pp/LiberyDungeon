package models

import (
	dungeon_models "libery-dungeon-libs/models"

	"github.com/golang-jwt/jwt"
)

type ClusterTokenClaims struct {
	jwt.StandardClaims
	Uuid           string `json:"uuid"`
	Name           string `json:"name"`
	FsPath         string `json:"fs_path"`
	FilterCategory string `json:"filter_category"`
	RootCategory   string `json:"root_category"`
}

func ParseClusterToken(cluster_token string, sk string) (ClusterTokenClaims, error) {
	claims := ClusterTokenClaims{}
	_, err := jwt.ParseWithClaims(cluster_token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return claims, err
}

func GetCategoriesClusterFromToken(cluster_token string, sk string) (*dungeon_models.CategoryCluster, error) {
	var new_cluster *dungeon_models.CategoryCluster = new(dungeon_models.CategoryCluster)

	cluster_claims, err := ParseClusterToken(cluster_token, sk)
	if err != nil {
		return nil, err
	}

	new_cluster.Uuid = cluster_claims.Uuid
	new_cluster.Name = cluster_claims.Name
	new_cluster.FsPath = cluster_claims.FsPath
	new_cluster.FilterCategory = cluster_claims.FilterCategory
	new_cluster.RootCategory = cluster_claims.RootCategory

	return new_cluster, nil
}
