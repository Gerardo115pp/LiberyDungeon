package models

type CategoryConfig struct {
	CategoryUUID         string   `json:"category_uuid"`
	BillboardMediaUUIDs  []string `json:"billboard_media_uuids"`
	BillboardDungeonTags []int    `json:"billboard_dungeon_tags"`
}

// Copies the non-default and non-nullish values from another category config except for the category uuid.
func (category_config *CategoryConfig) CopyNonDefaultValues(other_category_config *CategoryConfig) {
	if other_category_config.BillboardMediaUUIDs != nil && len(other_category_config.BillboardMediaUUIDs) > 0 {
		category_config.BillboardMediaUUIDs = make([]string, len(other_category_config.BillboardMediaUUIDs))

		copy(category_config.BillboardMediaUUIDs, other_category_config.BillboardMediaUUIDs)
	}

	if other_category_config.BillboardDungeonTags != nil && len(other_category_config.BillboardDungeonTags) > 0 {
		category_config.BillboardDungeonTags = make([]int, len(other_category_config.BillboardDungeonTags))

		copy(category_config.BillboardDungeonTags, other_category_config.BillboardDungeonTags)
	}
}

// Returns the default values used for a category configs.
func getDefaultCategoryConfig() *CategoryConfig {
	return &CategoryConfig{
		BillboardMediaUUIDs:  []string{},
		BillboardDungeonTags: []int{},
	}
}

// Creates a new category config.
func newCategoryConfig() *CategoryConfig {
	return getDefaultCategoryConfig()
}
