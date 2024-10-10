package request_parameters

import (
	"encoding/json"
	"net/http"
)

type SyncClusterPathRequest struct {
	ClusterUUID      string `json:"cluster_uuid"`
	SyncCategoryUUID string `json:"from_category_uuid"`
}

func NewSyncClusterPathFromRequest(request *http.Request) (sync_cluster_path_request *SyncClusterPathRequest, err error) {
	sync_cluster_path_request = &SyncClusterPathRequest{}
	err = json.NewDecoder(request.Body).Decode(sync_cluster_path_request)
	return
}
