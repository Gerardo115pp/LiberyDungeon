package metadata_requests

import (
	dungeon_helpers "libery-dungeon-libs/helpers"
	"net/http"
	"strconv"
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

/* ------------------------------ Video Moments ----------------------------- */

type VideoMoments_VideoIdentifier struct {
	VideoUUID    string `json:"video_uuid"`
	VideoCluster string `json:"video_cluster"`
}

type VideoMoments_VideoMomentData struct {
	MomentTime  int    `json:"moment_time"`
	MomentTitle string `json:"moment_title"`
}

type VideoMoments_NewVideoMoment struct {
	VideoMoments_VideoIdentifier
	VideoMoments_VideoMomentData
}

type VideoMoments_MomentIdentifier struct {
	MomentID int `json:"id"`
}

type VideoMoments_VideoMoment struct {
	VideoMoments_MomentIdentifier
	VideoMoments_VideoMomentData
}

func ParseVideoIdentifierParams(request *http.Request) *VideoMoments_VideoIdentifier {
	const (
		video_uuid_key    string = "video_uuid"
		video_cluster_key string = "video_cluster"
	)

	var request_body *VideoMoments_VideoIdentifier = new(VideoMoments_VideoIdentifier)

	request_body.VideoUUID = request.URL.Query().Get(video_uuid_key)
	request_body.VideoCluster = request.URL.Query().Get(video_cluster_key)

	return request_body
}

func ParseMomentIdentifierParams(request *http.Request) (*VideoMoments_MomentIdentifier, error) {
	var new_request_body *VideoMoments_MomentIdentifier = new(VideoMoments_MomentIdentifier)

	const (
		moment_id_key string = "id"
	)

	var moment_id_str = request.URL.Query().Get(moment_id_key)

	moment_id, err := strconv.Atoi(moment_id_str)
	if err != nil {
		return nil, err
	}

	new_request_body.MomentID = moment_id

	return new_request_body, nil
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
