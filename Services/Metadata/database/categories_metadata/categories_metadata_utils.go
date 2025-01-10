package categories_metadata

import (
	"encoding/json"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-metadata-service/models"
	"os"
	"path"
)

const CATEGORIES_METADATA_DIRECTORY_NAME string = "categories_metadata"

// prepares the given path to store the categories metadata. Returns the absolute path where the metadata will be stored
func bootCategoriesMetadataDB(config_path string) (string, error) {

	if config_path == "" || !dungeon_helpers.FileExists(config_path) {
		return "", fmt.Errorf("Config path <%s> does not exist", config_path)
	}

	var abs_categories_metadata_path string = path.Join(config_path, CATEGORIES_METADATA_DIRECTORY_NAME)

	if !dungeon_helpers.FileExists(abs_categories_metadata_path) {
		err := os.Mkdir(abs_categories_metadata_path, 0777)
		if err != nil {
			return "", fmt.Errorf("Error creating categories metadata directory: %s", err.Error())
		}
	}

	return abs_categories_metadata_path, nil
}

// Saves a category config with a non-nullish category uuid.
func saveCategoryConfig(category_config *models.CategoryConfig, save_path string) error {
	if category_config.CategoryUUID == "" {
		return fmt.Errorf("Category config has no category uuid")
	}

	if save_path == "" || !dungeon_helpers.FileExists(save_path) {
		return fmt.Errorf("Save path <%s> does not exist", save_path)
	}

	var category_config_file_path string = path.Join(save_path, category_config.CategoryUUID+".json")

	config_buffer, err := json.Marshal(category_config)
	if err != nil {
		return fmt.Errorf("Error marshalling category config: %s", err.Error())
	}

	err = os.WriteFile(category_config_file_path, config_buffer, 0777)
	if err != nil {
		return fmt.Errorf("Error writing category config file: %s", err.Error())
	}

	return nil
}

// Loads a category config file matching the given category uuid.
func loadCategoryConfig(category_uuid string, load_path string) (*models.CategoryConfig, error) {
	if category_uuid == "" {
		return nil, fmt.Errorf("Category uuid is empty")
	}

	if load_path == "" || !dungeon_helpers.FileExists(load_path) {
		return nil, fmt.Errorf("Load path <%s> does not exist", load_path)
	}

	var category_config_file_path string = path.Join(load_path, category_uuid+".json")

	config_file_bytes, err := os.ReadFile(category_config_file_path)
	if err != nil {
		return nil, fmt.Errorf("Error reading category config file: %s", err.Error())
	}

	var category_config *models.CategoryConfig = new(models.CategoryConfig)
	err = json.Unmarshal(config_file_bytes, category_config)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling category config: %s", err.Error())
	}

	return category_config, nil
}
