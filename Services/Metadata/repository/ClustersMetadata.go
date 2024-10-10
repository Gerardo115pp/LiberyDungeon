package repository

type ClusterMetadataRepository interface {
	AddPrivateCluster(cluster_uuid string)
	GetPrivateClusters() []string
	IsPrivateCluster(cluster_uuid string) bool
	RemovePrivateCluster(cluster_uuid string)
}

var ClusterMetadataRepo ClusterMetadataRepository

func SetClusterMetadataRepository(repo ClusterMetadataRepository) {
	ClusterMetadataRepo = repo
}
