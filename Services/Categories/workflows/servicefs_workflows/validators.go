package servicefs_workflows

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	service_models "libery_categories_service/models"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"os"
	"strings"
)

/**
 * Verifies that the subdirectory path is a valid path to scan for new cluster options.
 * Rules:
 * 1. The path cannot be the root directory aka "/"
 * 2. It must be under but not equal to a directory defined by the app_config.SERVICE_CLUSTERS_ROOT setting. If is_final_directory is false this rule is relaxed to allow SERVICE_CLUSTERS_ROOT
 * to be scanned for directories options.
 * 3. The path must exist or it's parent must exist.
 * 4. If the path exists, it must be a directory.
 * 5. The path cannot be equal or inside another cluster directory tree.
 */
func VerifyValidSubdirectoryScanPath(subdirectory_path string, is_final_directory bool) *dungeon_models.LabeledError {
	var normalized_path string = dungeon_helpers.NormalizePath(subdirectory_path)
	var labeled_err *dungeon_models.LabeledError
	var err error

	if normalized_path == "/" {
		err = fmt.Errorf("Refusing to add a cluster in the root directory")
		labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while is root directory", workflows.ErrForbiddenDirectoryScan)
		return labeled_err
	}

	if normalized_path == app_config.SERVICE_CLUSTERS_ROOT && is_final_directory {
		err = fmt.Errorf("Refusing to add the SERVICE_CLUSTERS_ROOT as a cluster as this would block the system from creating new clusters")
		labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while is service clusters root directory", workflows.ErrForbiddenDirectoryScan)
		return labeled_err
	}

	if !strings.HasPrefix(normalized_path, app_config.SERVICE_CLUSTERS_ROOT) {
		err = fmt.Errorf("Refusing to scan directory '%s' as is not a child of the service clusters root '%s'", subdirectory_path, app_config.SERVICE_CLUSTERS_ROOT)
		labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while is child of service clusters root", workflows.ErrForbiddenDirectoryScan)
		return labeled_err
	}

	var path_exists bool = dungeon_helpers.FileExists(normalized_path)

	if path_exists {
		subdirectory_stat, err := os.Stat(normalized_path)
		if err != nil {
			labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while getting directory stat", dungeon_models.ErrProcessError)
			return labeled_err
		}

		// If path exists, it must be a directory.
		if !subdirectory_stat.IsDir() {
			err = fmt.Errorf("The path '%s' is not a directory", subdirectory_path)
			labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while checking path is a directory", workflows.ErrNoSuchDirectory)
			return labeled_err
		}
	} else {
		// if path does not exist, it's parent must exist.
		parent_path := dungeon_helpers.GetParentDirectory(subdirectory_path)
		parent_path = dungeon_helpers.NormalizePath(parent_path)
		if !dungeon_helpers.FileExists(parent_path) {
			err = fmt.Errorf("The parent directory '%s' does not exist", parent_path)
			labeled_err = dungeon_models.NewLabeledError(err, "in verifyValidSubdirectoryScanPath, while checking parent directory exists", workflows.ErrNoSuchDirectory)
			return labeled_err
		}
	}

	_, labeled_err = isPathUsedInPlatform(subdirectory_path)

	return labeled_err
}

/**
 * Verifies that the subdirectory path is not 'in use' by another cluster. this means:
 * 1. The path is not equal to another cluster's path.
 * 2. The path is not a parent of another cluster's path.
 * 3. The path is not a child of another cluster's path.
 */
func isPathUsedInPlatform(subdirectory_path string) (bool, *dungeon_models.LabeledError) {
	var normalized_path string = dungeon_helpers.NormalizePath(subdirectory_path)
	var category_clusters []dungeon_models.CategoryCluster
	var is_path_used bool = false
	var labeled_err *dungeon_models.LabeledError

	category_clusters, err := repository.CategoriesClustersRepo.GetClusters(context.Background())
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "in isPathUsedInPlatfrom, while getting category clusters", service_models.ErrDB_CouldNotFindCategoryCluster)
		return is_path_used, labeled_err
	}

	for _, cluster := range category_clusters {
		path_is_exact_match_of_cluster := normalized_path == cluster.FsPath
		path_is_parent_of_cluster := strings.HasPrefix(cluster.FsPath, normalized_path)
		path_is_child_of_cluster := strings.HasPrefix(normalized_path, cluster.FsPath)

		if path_is_exact_match_of_cluster || path_is_parent_of_cluster || path_is_child_of_cluster {
			switch {
			case path_is_exact_match_of_cluster:
				is_path_used = true
				labeled_err = dungeon_models.NewLabeledError(fmt.Errorf("The path '%s' is an exact match of another cluster's path '%s'", subdirectory_path, cluster.Name), "in isPathUsedInPlatform, while checking exact match", workflows.ErrPathValidation_PathIsACluster)
			case path_is_parent_of_cluster:
				is_path_used = true
				labeled_err = dungeon_models.NewLabeledError(fmt.Errorf("The path '%s' is a parent of another cluster's path '%s'", subdirectory_path, cluster.Name), "in isPathUsedInPlatform, while checking parent of cluster", workflows.ErrPathValidation_PathIsAClusterAncestor)
			case path_is_child_of_cluster:
				is_path_used = true
				labeled_err = dungeon_models.NewLabeledError(fmt.Errorf("The path '%s' is a child of another cluster's path '%s'", subdirectory_path, cluster.Name), "in isPathUsedInPlatform, while checking child of cluster", workflows.ErrPathValidation_PathIsAClusterChild)
			}

			labeled_err.StoreVariable("cluster_id", cluster.Uuid)
			labeled_err.StoreVariable("cluster_name", cluster.Name)

			break
		}

	}

	return is_path_used, labeled_err
}
