package categories_metadata

import (
	"errors"
	"fmt"
	service_models "libery-metadata-service/models"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type categoryConfigDB struct {
	// The path where the category configs are stored
	categoriesConfigPath string
	// a map of category uuids -> category configs
	categoryConfigs map[string]*service_models.CategoryConfig
}

func NewCategoryConfigDB(config_path string) *categoryConfigDB {
	var new_category_config_db *categoryConfigDB = new(categoryConfigDB)

	categories_config_path, err := bootCategoriesMetadataDB(config_path)
	if err != nil {
		panic(errors.Join(fmt.Errorf("In database/categories_metadata/categories_metadata.NewCategoryConfigDB: while booting categories metadata database: "), err))
	}

	new_category_config_db.categoriesConfigPath = categories_config_path
	new_category_config_db.categoryConfigs = make(map[string]*service_models.CategoryConfig)

	return new_category_config_db
}

func (category_config_db *categoryConfigDB) CreateCategoryConfig(category_uuid string) *service_models.CategoryConfig {
	category_config := service_models.NewCategoryConfig() // Creates a new category config with default values.

	category_config.CategoryUUID = category_uuid

	category_config_db.categoryConfigs[category_uuid] = category_config

	err := saveCategoryConfig(category_config, category_config_db.categoriesConfigPath)
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("In database/categories_metadata/categories_metadata.CreateCategoryConfig: while saving category config to fs.\n\n%s", err))
	}

	return category_config
}

func (category_config_db *categoryConfigDB) GetCategoryConfig(category_uuid string) *service_models.CategoryConfig {
	var category_config, exists = category_config_db.categoryConfigs[category_uuid]

	if exists {
		return category_config
	}

	category_config, err := loadCategoryConfig(category_uuid, category_config_db.categoriesConfigPath)
	if err == nil {
		category_config_db.categoryConfigs[category_uuid] = category_config
		return category_config
	} else {
		echo.EchoWarn(fmt.Sprintf("In database/categories_metadata/categories_metadata.GetCategoryConfig: while loading category config from fs.\n\n%s", err))
	}

	return category_config_db.CreateCategoryConfig(category_uuid)
}

func (category_config_db *categoryConfigDB) UpdateCategoryConfig(category_config *service_models.CategoryConfig) error {
	if category_config.CategoryUUID == "" {
		return fmt.Errorf("In database/categories_metadata/categories_metadata.UpdateCategoryConfig: category config has no category uuid")
	}

	err := saveCategoryConfig(category_config, category_config_db.categoriesConfigPath)
	if err != nil {
		return fmt.Errorf("In database/categories_metadata/categories_metadata.UpdateCategoryConfig: while saving category config to fs.\n\n%s", err)
	}

	category_config_db.categoryConfigs[category_config.CategoryUUID] = category_config

	return nil
}
