package repository

import "libery-metadata-service/models"

type CategoriesConfigurationRepository interface {
	CreateCategoryConfig(category_uuid string) *models.CategoryConfig // Creates a new category config with default values.
	GetCategoryConfig(category_uuid string) *models.CategoryConfig
	UpdateCategoryConfig(category_config *models.CategoryConfig) error
}

var CategoriesConfigRepo CategoriesConfigurationRepository

func SetCategoriesConfigRepository(repo CategoriesConfigurationRepository) {
	CategoriesConfigRepo = repo
}
