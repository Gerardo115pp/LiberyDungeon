package clusters_metadata_database

import (
	"encoding/json"
	"fmt"
	dungeons_helpers "libery-dungeon-libs/helpers"
	app_config "libery-metadata-service/Config"
	"os"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func getClustersMetadataFilename() string {
	var filename string = fmt.Sprintf("%s.clusters_metadata.json", app_config.SERVICE_ID)

	return filepath.Join(app_config.OPERATION_DATA_PATH, filename)
}

func loadClustersMetadata() (*clusterMetadataStorage, error) {
	var clusters_metadata_filename string = getClustersMetadataFilename()

	if !dungeons_helpers.FileExists(clusters_metadata_filename) {
		echo.EchoWarn(fmt.Sprintf("Clusters metadata file<%s> not found", clusters_metadata_filename))
		return nil, nil
	}

	clusters_metadata, err := os.ReadFile(clusters_metadata_filename)
	if err != nil {
		return nil, err
	}

	var storage clusterMetadataStorage
	err = json.Unmarshal(clusters_metadata, &storage)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func saveClustersMetadata(storage *clusterMetadataStorage) error {
	var clusters_metadata_filename string = getClustersMetadataFilename()

	clusters_metadata, err := json.Marshal(storage)
	if err != nil {
		return err
	}

	err = os.WriteFile(clusters_metadata_filename, clusters_metadata, 0644)
	if err != nil {
		return err
	}

	return nil
}
