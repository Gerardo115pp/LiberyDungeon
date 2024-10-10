package models

type ClusterMetadata struct {
	ClusterUUID      string `json:"cluster_uuid"`
	IsClusterPrivate bool   `json:"is_cluster_private"`
}
