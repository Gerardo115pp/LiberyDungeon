package workflows

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	service_helpers "libery_categories_service/helpers"
	"libery_categories_service/repository"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CreateNewCategory(ctx context.Context, name string, parent_id string, cluster_id string) (category dungeon_models.Category, err error) {
	parent_category, err := repository.CategoriesRepo.GetCategory(context.Background(), parent_id)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("Couldn't get parent category: %s", parent_id))
		return
	}

	new_category := dungeon_models.CreateNewCategory(name, parent_category.Uuid, parent_category.Fullpath, cluster_id)

	// Create category path
	created, err := CreateNewCategoryPath(&new_category)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("Couldn't create category path: %s", new_category.Fullpath))
		return
	}

	if !created {
		err = fmt.Errorf("Couldn't create category path: %s", new_category.Fullpath)
		return
	}

	// Create category on database
	err = repository.CategoriesRepo.InsertCategory(ctx, new_category)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("Couldn't create category on database: %s", new_category.Fullpath))
		return
	}

	return new_category, err
}

func CreateNewCategoryPath(category *dungeon_models.Category) (created bool, err error) {
	category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(context.Background(), category.Cluster)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("Couldn't get category cluster: %s", category.Cluster))
		return
	}

	echo.Echo(echo.GreenFG, fmt.Sprintf("Creating category path: %s on %s", category.Fullpath, category_cluster.FsPath))

	if existsOnDungeonFS(category.Fullpath, category_cluster.FsPath) {
		setUniqueCategoryPath(category, &category_cluster)
	}

	// Ensure parent path exists
	parent_path := dungeon_helpers.GetParentDirectory(category.Fullpath)
	if !existsOnDungeonFS(parent_path, category_cluster.FsPath) {
		err = fmt.Errorf("Parent path doesn't exists: %s\ncategory path: %s", parent_path, category.Fullpath)
		return
	}

	system_path := getDungeonFSPath(category.Fullpath, category_cluster.FsPath)

	err = os.Mkdir(system_path, 0777)
	if err == nil {
		created = true
	}

	return
}

func DeleteCategory(category_uuid string, delete_content bool) (deleted bool, err error) {
	var category dungeon_models.Category
	var category_cluster dungeon_models.CategoryCluster

	// Get category
	category, err = repository.CategoriesRepo.GetCategory(context.Background(), category_uuid)
	if err != nil {
		return
	}

	// Get the category cluster
	category_cluster, err = repository.CategoriesClustersRepo.GetClusterByID(context.Background(), category.Cluster)
	if err != nil {
		return
	}

	category_identity := dungeon_models.CreateNewCategoryIdentity(&category, &category_cluster)

	// Get category childs(aka subcategories)
	childs, err := repository.CategoriesRepo.GetCategoryChildsByID(context.Background(), category_identity.Category.Uuid)
	if err != nil {
		return
	}

	// Get category medias
	medias, err := repository.CategoriesRepo.GetCategoryMedias(context.Background(), category_identity.Category.Uuid)
	if err != nil {
		return
	}

	if !delete_content && (len(childs) > 0 || len(medias) > 0) {
		return
	}

	var sub_category_deleted bool

	for _, child := range childs {
		sub_category_deleted, err = DeleteCategory(child.Uuid, delete_content)
		if err != nil || !sub_category_deleted {
			return
		}
	}

	// Delete category medias
	err = ProcessRejectedMedias(medias, category, &category_cluster)
	if err != nil {
		return
	}

	// Delete category
	err = repository.CategoriesRepo.DeleteCategory(context.Background(), category_identity.Category.Uuid)
	if err != nil {
		return
	}

	err = repository.TrashRepo.DeleteEmptyCategory(*category_identity)
	if err != nil {
		return
	}

	deleted = true

	return
}

func existsOnDungeonFS(relative_path string, cluster_path string) bool {
	var system_path string

	system_path = getDungeonFSPath(relative_path, cluster_path)

	return service_helpers.FileExists(system_path)
}

func getDungeonFSPath(relative_path string, cluster_path string) string {
	return path.Join(cluster_path, relative_path)
}

// Check if a provided category name is available(no other category with the same name) in th provided parent category. If match_casing is true, the comparison is case sensitive. else
// names will be converted to lowercase before comparison.
func IsCategoryNameAvailable(category_name string, parent_id string, match_casing bool) (available bool, err error) {
	siblings, err := repository.CategoriesRepo.GetCategoryChildsByID(context.Background(), parent_id)
	if err != nil {
		return
	}

	if !match_casing {
		category_name = strings.ToLower(category_name)
	}

	for _, sibling := range siblings {
		sibling_name := sibling.Name
		if !match_casing {
			sibling_name = strings.ToLower(sibling_name)
		}

		if sibling_name == category_name {
			return
		}
	}

	available = true

	return
}

func MoveCategory(category_uuid string, new_parent_uuid string, ctx context.Context) (err error) {
	// Get necessary data
	var moved_category dungeon_models.Category
	var new_parent_category dungeon_models.Category
	var category_cluster dungeon_models.CategoryCluster

	moved_category, err = repository.CategoriesRepo.GetCategory(ctx, category_uuid)
	if err != nil {
		return
	}

	new_parent_category, err = repository.CategoriesRepo.GetCategory(ctx, new_parent_uuid)
	if err != nil {
		return
	}

	category_cluster, err = repository.CategoriesClustersRepo.GetClusterByID(ctx, moved_category.Cluster)
	if err != nil {
		return
	}

	// Check that the receiver category is not a child of the moved category or the same category
	if new_parent_category.Uuid == moved_category.Uuid || strings.Contains(new_parent_category.Fullpath, moved_category.Fullpath) {
		err = fmt.Errorf("Receiver category is a child of the moved category or the same category: %s -> %s", new_parent_category.Uuid, moved_category.Uuid)
		return
	}

	// Check if new parent has already a category with the same name.
	old_path := getDungeonFSPath(moved_category.Fullpath, category_cluster.FsPath)
	new_path := getDungeonFSPath(fmt.Sprintf("%s/%s", new_parent_category.Fullpath, moved_category.Name), category_cluster.FsPath)

	directory_stat, err := os.Stat(new_path)
	if !os.IsNotExist(err) && directory_stat.IsDir() {
		err = fmt.Errorf("New parent category already has a category with the same name: %s. or there is a directory with the same name even if not related to a category", moved_category.Name)
		return
	}

	err = repository.CategoriesRepo.UpdateCategoryParent(ctx, moved_category, new_parent_category)
	if err != nil {
		return
	}

	// Move category on filesystem

	err = os.Rename(old_path, new_path)
	if err != nil {
		original_parent_category, err := repository.CategoriesRepo.GetCategory(ctx, moved_category.Parent)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Couldn't get original parent category: %s", moved_category.Parent))
		}

		err = repository.CategoriesRepo.UpdateCategoryParent(ctx, moved_category, original_parent_category)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Couldn't rollback category(%s) parent change from '%s' to '%s'", moved_category.Uuid, new_parent_category.Uuid, original_parent_category.Uuid))
		} else {
			echo.Echo(echo.OrangeBG, fmt.Sprintf("Couldn't move category(%s) path from '%s' to '%s'. Rolled back parent change", moved_category.Uuid, old_path, new_path))
		}
	}

	return
}

func RenameCategory(category_uuid string, new_name string) (err error) {
	// Get necessary data
	var category_cluster dungeon_models.CategoryCluster
	var category dungeon_models.Category

	category, err = repository.CategoriesRepo.GetCategory(context.Background(), category_uuid)
	if err != nil {
		return
	}

	category_cluster, err = repository.CategoriesClustersRepo.GetClusterByID(context.Background(), category.Cluster)
	if err != nil {
		return
	}

	// Check if new name is available
	available, err := IsCategoryNameAvailable(new_name, category.Parent, true)
	if err != nil {
		return
	}

	if !available {
		err = fmt.Errorf("Category name is not available: %s", new_name)
		return
	}

	// Rename category on database

	err = repository.CategoriesRepo.UpdateCategoryName(context.Background(), category, new_name)
	if err != nil {
		return
	}

	old_path := getDungeonFSPath(category.Fullpath, category_cluster.FsPath)
	new_path := fmt.Sprintf("%s/%s", filepath.Dir(old_path), new_name)

	err = os.Rename(old_path, new_path)
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("In workflows.RenameCategory: Couldn't rename category(%s) path from %s to %s. because: %s", category_uuid, old_path, new_path, err.Error()))

		new_category_state := new(dungeon_models.Category)
		new_category_state.CopyContent(category)
		new_category_state.Name = new_name
		new_category_state.Fullpath = new_path

		err = repository.CategoriesRepo.UpdateCategoryName(context.Background(), *new_category_state, category.Name)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Couldn't rollback category(%s) name change from %s to %s", category_uuid, new_name, category.Name))
		}

	}

	return
}

func setUniqueCategoryPath(category *dungeon_models.Category, cluster *dungeon_models.CategoryCluster) {
	var path_exists bool
	var new_path string = category.Fullpath

	path_exists = existsOnDungeonFS(category.Fullpath, cluster.FsPath)

	for path_exists {
		path_version_counter := 1
		new_path = fmt.Sprintf("%s_v%d", category.Fullpath, path_version_counter)
		path_exists = existsOnDungeonFS(new_path, cluster.FsPath)
	}

	category.Fullpath = new_path

	return
}
