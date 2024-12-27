package metadata_requests

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
