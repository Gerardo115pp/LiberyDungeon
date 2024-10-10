package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type CategoriesRepository interface {
	GetCategoryByID(ctx context.Context, category_id string) (dungeon_models.Category, error)
	Close() error
}

var CategoriesRepo CategoriesRepository

func SetCategoriesImplementation(impl CategoriesRepository) {
	CategoriesRepo = impl
}
