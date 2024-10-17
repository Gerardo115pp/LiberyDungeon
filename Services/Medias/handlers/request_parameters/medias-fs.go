package request_parameters

import (
	"encoding/json"
	"net/http"
)

type SequenceRenameParams struct {
	CategoryUUID    string            `json:"category_uuid"`    // All medias must share the same main category
	SequenceMembers map[string]string `json:"sequence_members"` // media_uuid -> new_name
	NamePrefix      string            `json:"name_prefix"`
}

func NewSequenceRenameParamsFromRequest(request *http.Request) (sequence_rename_params *SequenceRenameParams, err error) {
	sequence_rename_params = &SequenceRenameParams{}
	err = json.NewDecoder(request.Body).Decode(sequence_rename_params)
	return
}
