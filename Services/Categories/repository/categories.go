package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type CategoriesRepository interface {
	GetCategoryChildsByID(ctx context.Context, category_id string) ([]dungeon_models.ChildCategory, error)
	GetCategoryMedias(ctx context.Context, category_id string) ([]dungeon_models.Media, error)
	GetCategoryContent(ctx context.Context, category_id string) (*dungeon_models.CategoryLeaf, error)
	GetCategory(ctx context.Context, category_id string) (dungeon_models.Category, error)
	GetCategories(ctx context.Context, category_ids []string) ([]dungeon_models.Category, error)
	GetClusterCategories(ctx context.Context, cluster_id string) ([]dungeon_models.Category, error)
	GetAllCategories(ctx context.Context) ([]dungeon_models.Category, error)
	GetCategoryContentByFullpath(ctx context.Context, category_path string, category_cluster string) (*dungeon_models.Category, error)
	GetCategoryFSBranch(ctx context.Context, category_id string) ([]dungeon_models.MediaWeakIdentity, error)
	DeleteCategoryMedias(ctx context.Context, medias []dungeon_models.Media) error
	DeleteCategory(ctx context.Context, category_id string) error
	IsCategoryEmpty(ctx context.Context, category_id string) (bool, error)
	InsertCategory(ctx context.Context, category dungeon_models.Category) error
	UpdateMedia(ctx context.Context, media dungeon_models.Media) error
	UpdateMedias(ctx context.Context, medias []dungeon_models.Media) error
	UpdateCategoryName(ctx context.Context, category dungeon_models.Category, new_name string) error
	UpdateCategoryParent(ctx context.Context, category dungeon_models.Category, new_parent dungeon_models.Category) error
	UpdateCategoryThumbnail(ctx context.Context, category_uuid, new_thumbnail string) error
}

var CategoriesRepo CategoriesRepository

func SetCategoriesImplementation(impl CategoriesRepository) {
	CategoriesRepo = impl
}
