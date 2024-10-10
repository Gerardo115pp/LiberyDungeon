package servicefs_workflows

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery_categories_service/models"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Returns a list of valid directory options from which to create a new category cluster.
 */
func GetDirectoryOptionsFromPath(subdirectory_path string) ([]service_models.DirectoryOption, *dungeon_models.LabeledError) {
	var directory_options []service_models.DirectoryOption = make([]service_models.DirectoryOption, 0)
	var labeled_err *dungeon_models.LabeledError
	var err error

	if labeled_err = VerifyValidSubdirectoryScanPath(subdirectory_path, false); labeled_err != nil {
		fmt.Printf("Error label: %s\n", labeled_err.Label)
		if labeled_err.Label != workflows.ErrPathValidation_PathIsAClusterAncestor {
			labeled_err.AppendContext("in GetDirectoyOptionsFromPath, while verifying valid subdirectory scan path")
			return directory_options, labeled_err
		}
	}

	subdirectories, err := os.ReadDir(subdirectory_path)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "in GetDirectoyOptionsFromPath, while calling os.ReadDir", dungeon_models.ErrProcessError)
		return directory_options, labeled_err
	}

	var new_cluster_option service_models.DirectoryOption
	for _, subdirectory := range subdirectories {

		if subdirectory.IsDir() && !strings.HasPrefix(subdirectory.Name(), ".") {
			directory_path := filepath.Join(subdirectory_path, subdirectory.Name())

			new_cluster_option = service_models.DirectoryOption{
				Name: subdirectory.Name(),
				Path: directory_path,
			}

			directory_options = append(directory_options, new_cluster_option)
		}
	}

	return directory_options, nil
}

/**
 * Creates a new category cluster. The cluster will be created in the directory specified by the new_cluster.FsPath field.
 */
func CreateNewCategoryCluster(new_cluster *dungeon_models.CategoryCluster) (*dungeon_models.CategoryCluster, *dungeon_models.LabeledError) {
	var labeled_err *dungeon_models.LabeledError
	var err error

	labeled_err = VerifyValidSubdirectoryScanPath(new_cluster.FsPath, true)
	if labeled_err != nil {
		labeled_err.AppendContext("in CreateNewCategoryCluster, while verifying valid subdirectory scan path")
		return nil, labeled_err
	}

	if !dungeon_helpers.FileExists(new_cluster.FsPath) {
		err = os.MkdirAll(new_cluster.FsPath, 0755)
		if err != nil {
			labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling os.MkdirAll", dungeon_models.ErrProcessError)
			return nil, labeled_err
		}
	}

	root_category := createMainCategory(new_cluster.Uuid)
	new_cluster.RootCategory = root_category.Uuid

	// Right now, the filter category field in new_cluster is the name of the category to use as filter category. not the uuid.
	filter_category := dungeon_models.CreateNewCategory(new_cluster.FilterCategory, root_category.Uuid, root_category.Fullpath, new_cluster.Uuid)
	new_cluster.FilterCategory = filter_category.Uuid

	filter_category_fs_path := filepath.Join(new_cluster.FsPath, filter_category.Fullpath)
	filter_category_path_existed := dungeon_helpers.FileExists(filter_category_fs_path)

	if !filter_category_path_existed {
		err = os.MkdirAll(filter_category_fs_path, 0755)
		if err != nil {
			labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling os.MkdirAll for filter category", dungeon_models.ErrProcessError)
			return nil, labeled_err
		}
	}

	err = repository.CategoriesClustersRepo.InsertCluster(context.Background(), *new_cluster, *root_category, filter_category)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, fmt.Sprintf("in CreateNewCategoryCluster, while calling CategoriesClustersRepo.InsertCluster on:\n%s", new_cluster), dungeon_models.ErrProcessError)
		return nil, labeled_err
	}

	if filter_category_path_existed {
		filter_category_content, err := os.ReadDir(filter_category_fs_path)
		if err != nil {
			labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling os.ReadDir for filter category", dungeon_models.ErrProcessError)
			return nil, labeled_err
		}

		for _, f := range filter_category_content {
			if f.IsDir() {
				cluster_identity := new_cluster.ToWeakIdentity()
				labeled_err := CreateClusterTree(cluster_identity, filepath.Join(filter_category.Fullpath, f.Name()), filter_category.Uuid)
				if labeled_err != nil {
					labeled_err.AppendContext(fmt.Sprintf(":%s:", f.Name()))
					return nil, labeled_err
				}
			} else if dungeon_helpers.IsSupportedFileExtension(f.Name()) {
				new_media := dungeon_models.CreateNewMedia(f.Name(), filter_category.Uuid, dungeon_helpers.IsVideoFile(f.Name()), 0)

				err = repository.MediasRepo.InsertMedia(context.Background(), new_media)
				if err != nil {
					labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling MediasRepo.InsertMedia for filter category", dungeon_models.ErrProcessError)
					return nil, labeled_err
				}
			}
		}
	}

	new_cluster_contents, err := os.ReadDir(new_cluster.FsPath)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling os.ReadDir for new cluster", dungeon_models.ErrProcessError)
		return nil, labeled_err
	}

	for _, f := range new_cluster_contents {
		if f.IsDir() && filter_category.Name != f.Name() {
			cluster_identity := new_cluster.ToWeakIdentity()
			labeled_err := CreateClusterTree(cluster_identity, f.Name(), root_category.Uuid)
			if labeled_err != nil {
				labeled_err.AppendContext(fmt.Sprintf(":%s:", f.Name()))
				return nil, labeled_err
			}
		} else if dungeon_helpers.IsSupportedFileExtension(f.Name()) {
			new_media := dungeon_models.CreateNewMedia(f.Name(), root_category.Uuid, dungeon_helpers.IsVideoFile(f.Name()), 0)

			err = repository.MediasRepo.InsertMedia(context.Background(), new_media)
			if err != nil {
				labeled_err = dungeon_models.NewLabeledError(err, "in CreateNewCategoryCluster, while calling MediasRepo.InsertMedia for root category", dungeon_models.ErrProcessError)
				return nil, labeled_err
			}
		}
	}

	return new_cluster, nil
}

func CreateClusterTree(new_cluster_identity *dungeon_models.CategoryClusterWeakIdentity, subdirectory_path string, parent_uuid string) (labeled_err *dungeon_models.LabeledError) {
	var category_name string = filepath.Base(subdirectory_path)
	var fs_path string = filepath.Join(new_cluster_identity.ClusterFsPath, subdirectory_path)

	var category_parent_directory string = dungeon_helpers.GetParentDirectory(subdirectory_path)

	var new_category dungeon_models.Category = dungeon_models.CreateNewCategory(category_name, parent_uuid, category_parent_directory, new_cluster_identity.ClusterUUID)

	err := repository.CategoriesRepo.InsertCategory(context.Background(), new_category)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "in createClusterTree, while calling CategoriesRepo.InsertCategory", dungeon_models.ErrProcessError)
		return labeled_err
	}

	directory_contents, err := os.ReadDir(fs_path)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "in createClusterTree, while calling os.ReadDir", dungeon_models.ErrProcessError)
		return labeled_err
	}

	for _, f := range directory_contents {
		if f.IsDir() {
			child_path := filepath.Join(subdirectory_path, f.Name())
			labeled_err = CreateClusterTree(new_cluster_identity, child_path, new_category.Uuid)
			if labeled_err != nil {
				labeled_err.AppendContext(fmt.Sprintf(":%s:", child_path))
				return labeled_err
			}
		} else {
			if dungeon_helpers.IsSupportedFileExtension(f.Name()) {
				new_media := dungeon_models.CreateNewMedia(f.Name(), new_category.Uuid, dungeon_helpers.IsVideoFile(f.Name()), 0)
				err = repository.MediasRepo.InsertMedia(context.Background(), new_media)
				if err != nil {
					labeled_err = dungeon_models.NewLabeledError(err, "in createClusterTree, while calling MediasRepo.InsertMedia", dungeon_models.ErrProcessError)
					return labeled_err
				}
			}
		}
	}

	return nil
}

func createMainCategory(cluster_id string) *dungeon_models.Category {
	var new_category dungeon_models.Category = dungeon_models.CreateNewCategory("main", "", "/", cluster_id)

	new_category.Fullpath = "/"

	return &new_category
}
