package metadata_requests

// -------------------- Dungeon tags --------------------

type MultiTagEntitiesRequest struct {
	DungeonTags []int    `json:"dungeon_tags"`
	EntityUUIDS []string `json:"entity_uuids"`
	EntityType  string   `json:"entity_type"`
}

type MultiTagEntityRequest struct {
	DungeonTags []int  `json:"dungeon_tags"`
	EntityUUID  string `json:"entity_uuid"`
	EntityType  string `json:"entity_type"`
}

type TagEntitiesRequest struct {
	DungeonTagID  int      `json:"tag_id"`
	EntityType    string   `json:"entity_type"`
	EntitiesUUIDs []string `json:"entities_uuids"`
}

// -------------------- Categories config --------------------

type PatchCategoryBillboardTagsRequest struct {
	CategoryUUID         string `json:"category_uuid"`
	BillboardDungeonTags []int  `json:"billboard_tags"`
}

type PatchCategoryBillboardMediasRequest struct {
	CategoryUUID        string   `json:"category_uuid"`
	BillboardMediaUUIDs []string `json:"billboard_media_uuids"`
}
