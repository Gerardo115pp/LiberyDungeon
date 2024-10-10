package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type CategoriesClustersRepository interface {
	GetClusterByID(ctx context.Context, cluster_id string) (dungeon_models.CategoryCluster, error)
	GetCategoryCluster(ctx context.Context, category_id string) (*dungeon_models.CategoryCluster, error)
}

var CategoriesClustersRepo CategoriesClustersRepository

func SetCategoriesClustersImplementation(impl CategoriesClustersRepository) {
	CategoriesClustersRepo = impl
}
