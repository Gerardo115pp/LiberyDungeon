package fs_sync

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	"libery_categories_service/repository"
	"libery_categories_service/workflows/servicefs_workflows"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

// Verfies the fs state of a give category matches its database state. Each supported file found existing in the fs but not in the db will be inserted into the db.
// Each file an category found in the db but not in the fs will be removed from the db.
func SyncCategoryBranch(category_identity *dungeon_models.CategoryIdentity) *dungeon_models.LabeledError {
	var branch_content []dungeon_models.MediaWeakIdentity
	var lerr *dungeon_models.LabeledError

	echo.EchoDebug(fmt.Sprintf("Syncing category<%s> branch: '%s'", category_identity.Category.Uuid, category_identity.Category.Fullpath))

	branch_content, err := repository.CategoriesRepo.GetCategoryFSBranch(context.Background(), category_identity.Category.Uuid)
	if err != nil {
		return dungeon_models.NewLabeledError(err, "in SyncCategoryBranch, while calling CategoriesRepo.GetCategoryFSBranch", dungeon_models.ErrProcessError)
	}

	var branch_path string = filepath.Join(category_identity.ClusterPath, category_identity.Category.Fullpath)
	branch_path = dungeon_helpers.NormalizePath(branch_path)

	db_state_map := buildFsPathLookupTable(&branch_content, category_identity.ClusterPath)
	if app_config.DEBUG_MODE {
		// buildFsPathLookupTable modifies the branch_content slice, that why we wait until now to print it
		echo.EchoDebug(fmt.Sprintf("Category<%s> branch content: %d", category_identity.Category.Uuid, len(branch_content)))
		for _, media := range branch_content {
			echo.EchoDebug(fmt.Sprintf("%s", media.String()))
		}

	}

	sync_errors, lerr := scanSyncErrors(branch_path, db_state_map)
	if lerr != nil {
		lerr.AppendContext(fmt.Sprintf("In SyncCategoryBranch, while scanning category branch: '%s'", branch_path))
		return lerr
	}

	sync_errors.reportGhostFiles(category_identity, branch_content)

	if app_config.DEBUG_MODE {
		echo.Echo(echo.WhiteFG, sync_errors.String())
	}

	lerr = amendSyncErrors(sync_errors, category_identity)

	return lerr
}

func amendSyncErrors(sync_errors *stateSyncErrors, category_identity *dungeon_models.CategoryIdentity) *dungeon_models.LabeledError {
	var err error
	var lerr *dungeon_models.LabeledError

	echo.EchoDebug(fmt.Sprintf("%sAmending sync errors%s", echo.BlueFG, echo.CyanFG))
	lerr = syncUnregisteredContent(sync_errors, category_identity)
	if lerr != nil {
		return lerr
	}

	echo.EchoDebug(fmt.Sprintf("%sAmending ghost identities%s", echo.BlueFG, echo.CyanFG))
	err = syncGhostIdentities(sync_errors)
	if err != nil {
		lerr = dungeon_models.NewLabeledError(err, "in amendSyncErrors, while calling syncGhostIdentities", dungeon_models.ErrDB_CouldNotConnectToDB)
		return lerr
	}

	return nil
}

func syncUnregisteredContent(sync_errors *stateSyncErrors, root_identity *dungeon_models.CategoryIdentity) *dungeon_models.LabeledError {
	var err error
	var lerr *dungeon_models.LabeledError

	for _, unregistered_file := range sync_errors.UnregisteredFiles {
		err = syncUnregisteredMedia(unregistered_file)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Error could not sync unregistered media<%s> because: %s", unregistered_file.FilePath, err.Error()))
		}
	}

	for _, unregistered_category := range sync_errors.UnregisteredCategoriesPaths {
		lerr = syncUnregisteredCategory(unregistered_category, root_identity)
		if lerr != nil {
			echo.EchoWarn(fmt.Sprintf("Error could not sync unregistered category<%s> because: %s", unregistered_category.DirectoryPath, lerr.Error()))
		}
	}

	return nil
}

func syncUnregisteredMedia(unregistered_file unregisteredFile) error {
	file_stat, err := os.Stat(unregistered_file.FilePath)
	if err != nil {
		return err
	}

	var filename string = file_stat.Name()

	new_media := dungeon_models.CreateNewMedia(filename, unregistered_file.CategoryUUID, dungeon_helpers.IsVideoFile(filename), 0)

	echo.EchoDebug(fmt.Sprintf("-> Inserting new media: %s", new_media.Name))

	err = repository.MediasRepo.InsertMedia(context.Background(), new_media)

	return err
}

func syncUnregisteredCategory(unregistered_category unregisteredCategory, root_identity *dungeon_models.CategoryIdentity) *dungeon_models.LabeledError {
	cluster_identity := root_identity.ToClusterWeakIdentity()
	file_relative_path := strings.TrimPrefix(unregistered_category.DirectoryPath, cluster_identity.ClusterFsPath)
	file_relative_path = dungeon_helpers.NormalizePath(file_relative_path)

	echo.EchoDebug(fmt.Sprintf("-> Creating new category tree: %s", file_relative_path))
	err := servicefs_workflows.CreateClusterTree(cluster_identity, file_relative_path, unregistered_category.ParentUUID)

	return err
}

func syncGhostIdentities(sync_errors *stateSyncErrors) error {
	var err error

	for _, ghost_identity := range sync_errors.GhostIdentities {
		var is_media bool = ghost_identity.MediaName != ""

		if is_media {
			err = syncGhostMedia(ghost_identity)
		} else {
			err = syncGhostCategory(ghost_identity)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func syncGhostMedia(identity dungeon_models.MediaWeakIdentity) error {
	err := repository.MediasRepo.DeleteMedia(context.Background(), identity.MediaUUID)

	return err
}

func syncGhostCategory(identity dungeon_models.MediaWeakIdentity) error {
	err := repository.CategoriesRepo.DeleteCategory(context.Background(), identity.CategoryUUID)

	return err
}
