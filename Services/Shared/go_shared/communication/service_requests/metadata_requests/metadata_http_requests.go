package metadata_requests

import (
	dungeon_helpers "libery-dungeon-libs/helpers"
	"net/http"
)

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

type TagListRequest struct {
	TagList []int `json:"tag_list"`
}

func ParseTagIDsListParams(request *http.Request) (*TagListRequest, error) {
	const (
		tag_list_key string = "tags"
	)

	var request_body *TagListRequest = new(TagListRequest)

	tag_list, err := dungeon_helpers.ParseQueryParameterAsIntSlice(request, tag_list_key)
	if err != nil {
		return nil, err
	}

	request_body.TagList = tag_list

	return request_body, nil
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
