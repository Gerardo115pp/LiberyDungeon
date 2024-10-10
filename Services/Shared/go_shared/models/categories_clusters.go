package dungeon_models

import "fmt"

type CategoryCluster struct {
	Uuid           string `json:"uuid"`
	Name           string `json:"name"`
	FsPath         string `json:"fs_path"`
	FilterCategory string `json:"filter_category"`
	RootCategory   string `json:"root_category"`
}

func (cc *CategoryCluster) ToWeakIdentity() *CategoryClusterWeakIdentity {
	return &CategoryClusterWeakIdentity{
		ClusterUUID:   cc.Uuid,
		ClusterFsPath: cc.FsPath,
	}
}

// no cluster identity a cluster only depends on itself to be identified both on the database and on the filesystem

// Weak identity is still useful as the purpose of weak identities is saving database requests. a CategoryClusterWeakIdentity can be
// generated from a CategoryIdentity, MediaIdentity and MediaWeakIdentity.
type CategoryClusterWeakIdentity struct {
	ClusterUUID   string `json:"cluster_uuid"`
	ClusterFsPath string `json:"cluster_fs_path"`
}

func (cc *CategoryCluster) UpdateClusterData(other_cluster *CategoryCluster) error {
	if other_cluster.Uuid != cc.Uuid {
		return fmt.Errorf("Clusters UUIDs do not match")
	}

	cc.Name = other_cluster.Name
	cc.FsPath = other_cluster.FsPath
	cc.FilterCategory = other_cluster.FilterCategory
	cc.RootCategory = other_cluster.RootCategory

	return nil
}

func (cc *CategoryCluster) String() string {
	var cluster_string string = ""

	cluster_string += fmt.Sprintf("UUID: %s\n", cc.Uuid)
	cluster_string += fmt.Sprintf("Name: %s\n", cc.Name)
	cluster_string += fmt.Sprintf("FsPath: %s\n", cc.FsPath)
	cluster_string += fmt.Sprintf("FilterCategory: %s\n", cc.FilterCategory)
	cluster_string += fmt.Sprintf("RootCategory: %s\n", cc.RootCategory)

	return cluster_string
}
