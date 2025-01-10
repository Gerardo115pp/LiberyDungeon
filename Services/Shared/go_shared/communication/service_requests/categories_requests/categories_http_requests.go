package categories_requests

import (
	dungeon_helpers "libery-dungeon-libs/helpers"
	"net/http"
)

type GetMediaListRequest struct {
	ClusterUUID string   `json:"cluster_uuid"`
	MediaUUIDs  []string `json:"media_uuids"`
}

func ParseMediaListParams(request *http.Request) (*GetMediaListRequest, error) {
	const (
		cluster_uuid_key string = "cluster_uuid"
		media_uuids_key  string = "media_uuids"
	)

	var request_body *GetMediaListRequest = new(GetMediaListRequest)

	var cluster_uuid string = request.URL.Query().Get(cluster_uuid_key)

	media_uuids, err := dungeon_helpers.ParseQueryParameterAsStringSlice(request, media_uuids_key)
	if err != nil {
		return nil, err
	}

	request_body.ClusterUUID = cluster_uuid
	request_body.MediaUUIDs = media_uuids

	return request_body, nil
}
