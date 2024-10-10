package clusters_metadata_database

import (
	"fmt"
	dungeons_helpers "libery-dungeon-libs/helpers"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type clusterMetadataStorage struct {
	PrivateClusters []string `json:"private_clusters"` // List of cluster UUIDs that are private
}

type ClusterMetadataDB struct {
	storage  *clusterMetadataStorage
	filename string
}

func NewClusterMetadataDB() *ClusterMetadataDB {
	var storage_content *clusterMetadataStorage = new(clusterMetadataStorage)
	var filename string = getClustersMetadataFilename()
	var err error

	if dungeons_helpers.FileExists(filename) {
		storage_content, err = loadClustersMetadata()
		if err != nil {
			echo.EchoFatal(fmt.Errorf("In MetadataService/databases/clusters_metadata/clusters_metadata.go while loading clusters metadata file: %s", err))
		}
	} else {
		err = saveClustersMetadata(storage_content)
		if err != nil {
			echo.EchoFatal(fmt.Errorf("In MetadataService/databases/clusters_metadata/clusters_metadata.go while creating clusters metadata file: %s", err))
		}
	}

	return &ClusterMetadataDB{
		storage:  storage_content,
		filename: filename,
	}
}

func (db *ClusterMetadataDB) AddPrivateCluster(cluster_uuid string) {
	db.storage.PrivateClusters = append(db.storage.PrivateClusters, cluster_uuid)
	db.save()
}

func (db *ClusterMetadataDB) GetPrivateClusters() []string {
	return db.storage.PrivateClusters
}

func (db *ClusterMetadataDB) IsPrivateCluster(cluster_uuid string) bool {
	var is_private bool = false

	for _, cluster := range db.storage.PrivateClusters {
		if cluster == cluster_uuid {
			is_private = true
			break
		}
	}

	return is_private

}

func (db *ClusterMetadataDB) RemovePrivateCluster(cluster_uuid string) {
	var new_private_clusters []string

	for _, cluster := range db.storage.PrivateClusters {
		if cluster != cluster_uuid {
			new_private_clusters = append(new_private_clusters, cluster)
		}
	}

	db.storage.PrivateClusters = new_private_clusters
	db.save()
}

func (db *ClusterMetadataDB) save() {
	err := saveClustersMetadata(db.storage)
	if err != nil {
		echo.EchoFatal(fmt.Errorf("In MetadataService/databases/clusters_metadata/clusters_metadata.go while saving clusters metadata file: %s", err))
	}
}
