package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type CategoriesClustersRepository interface {
	GetClusterByID(ctx context.Context, cluster_id string) (dungeon_models.CategoryCluster, error)
	GetCategoryCluster(ctx context.Context, category_id string) (dungeon_models.CategoryCluster, error)
	GetClusters(ctx context.Context) ([]dungeon_models.CategoryCluster, error)
	InsertCluster(ctx context.Context, cluster dungeon_models.CategoryCluster, root_category dungeon_models.Category, filter_category dungeon_models.Category) error
	UpdateCluster(ctx context.Context, cluster dungeon_models.CategoryCluster) error
	DeleteCluster(ctx context.Context, cluster_id string) *dungeon_models.LabeledError
}

var CategoriesClustersRepo CategoriesClustersRepository

func SetCategoriesClustersImplementation(impl CategoriesClustersRepository) {
	CategoriesClustersRepo = impl
}
