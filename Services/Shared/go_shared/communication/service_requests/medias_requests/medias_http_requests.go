package medias_http_requests

type RenameMediaRequest struct {
	NewName   string `json:"new_name"`
	MediaUUID string `json:"media_uuid"`
}
