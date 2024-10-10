package repository

import "context"

type CategoriesRepository interface {
	CreateCategory(ctx context.Context, name string, parent_uuid string) (category_uuid string, err error)
}

var Categories CategoriesRepository

func SetCategoriesImplementation(implementation CategoriesRepository) {
	Categories = implementation
}
