package workflows

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/helpers"
	"libery_categories_service/repository"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/google/uuid"
)

func ProcessRejectedMedias(rejected_medias []dungeon_models.Media, rejected_from dungeon_models.Category, cluster_rejected_from *dungeon_models.CategoryCluster) error {
	var category_path string = filepath.Join(cluster_rejected_from.FsPath, rejected_from.Fullpath)
	var err error

	category_identity := dungeon_models.CreateNewCategoryIdentity(&rejected_from, cluster_rejected_from)

	echo.Echo(echo.PinkBG, fmt.Sprintf("Moving %d medias from %s to trash", len(rejected_medias), category_path))

	err = repository.TrashRepo.StartTransaction(category_identity.ToWeakIdentity())
	if err != nil {
		return err
	}

	for _, media := range rejected_medias {
		media_identity := dungeon_models.CreateNewMediaIdentity(&media, &rejected_from, cluster_rejected_from)
		err = repository.TrashRepo.MoveToTrash(media_identity)
		if err != nil {
			restoration_error := repository.TrashRepo.Rollback()
			if restoration_error != nil {
				echo.Echo(echo.RedFG, fmt.Sprintf("Error restoring medias from trash: %s", restoration_error.Error()))
			}
			return err
		}
	}

	err = repository.CategoriesRepo.DeleteCategoryMedias(context.Background(), rejected_medias)
	if err != nil {
		restoration_error := repository.TrashRepo.Rollback()
		if restoration_error != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error restoring medias from trash: %s", restoration_error.Error()))
		}
		return err
	}

	err = repository.TrashRepo.Commit()

	return err
}

// ProcessMovedMedias moves medias from one category to several other categories. and calles the repository to update the database. if one error occurs, it rolls back all the changes.
// Parameters:
//
//	moved_medias: map of new_category_id -> medias to move to that category. new category doesn't mean a category that doesn't exist, it means the category to which the medias are being moved to.
//	current_category: category from which the medias are being moved
//	medias_cluster: cluster that containes the medias. for now we don't support moving medias from one cluster to another.
//
// Returns:
//
//	error: nil if everything went ok, error otherwise
func ProcessMovedMedias(moved_medias map[string][]dungeon_models.Media, current_category dungeon_models.Category, medias_cluster *dungeon_models.CategoryCluster) error {
	var err error
	var moved_to map[string]string = make(map[string]string) // where the media was moved to. uuid -> new_path. used to rollback in case of error
	var updated_medias []dungeon_models.Media = make([]dungeon_models.Media, 0)
	var new_category dungeon_models.Category
	var old_path string
	var new_path string

	echo.Echo(echo.PinkBG, fmt.Sprintf("Moving %d medias", len(moved_medias)))

	for category_id, medias := range moved_medias {
		new_category, err = repository.CategoriesRepo.GetCategory(context.Background(), category_id)
		if err != nil {
			echo.EchoDebug(fmt.Sprintf("[getting new_category]Error getting category %s: %s", category_id, err.Error()))
			return err
		}

		for _, media := range medias {
			old_path = filepath.Join(medias_cluster.FsPath, current_category.Fullpath, media.Name)
			new_path = filepath.Join(medias_cluster.FsPath, new_category.Fullpath, media.Name)

			// if helpers.FileExists(new_path) {
			// 	return fmt.Errorf("File %s already exists in category %s", media.Name, new_category.Name) // TODO: Create new file name
			// }

			for helpers.FileExists(new_path) {
				GetUuidMediaName(&media)
				new_path = filepath.Join(medias_cluster.FsPath, new_category.Fullpath, media.Name)
			}

			err = os.Rename(old_path, new_path)
			if err != nil {
				new_err := rollbackMovedMedias(moved_to, updated_medias, old_path)
				if new_err != nil {
					echo.EchoWarn(fmt.Sprintf("[After os.Rename]Error restoring medias from %s to %s: %s", old_path, new_path, new_err.Error()))
				}
				return err
			}

			moved_to[media.Uuid] = new_path
			media.MainCategory = new_category.Uuid

			updated_medias = append(updated_medias, media)
		}
	}

	echo.EchoDebug(fmt.Sprintf("Updating %d medias", len(updated_medias)))

	err = repository.CategoriesRepo.UpdateMedias(context.Background(), updated_medias)
	if err != nil {
		new_err := rollbackMovedMedias(moved_to, updated_medias, old_path)
		if new_err != nil {
			echo.EchoWarn(fmt.Sprintf("[After db update]Error restoring medias from %s to %s: %s", old_path, new_path, new_err.Error()))
		}
		return err
	}

	return nil
}

func rollbackMovedMedias(moved_to map[string]string, medias []dungeon_models.Media, original_path string) error {
	var err error

	for _, media := range medias {
		new_path, ok := moved_to[media.Uuid]
		if !ok {
			continue
		}

		err = os.Rename(new_path, original_path)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Error restoring media %s from %s to %s: %s", media.Name, new_path, original_path, err.Error()))
			return err
		}

		delete(moved_to, media.Uuid)
	}

	if len(moved_to) > 0 {
		echo.EchoWarn(fmt.Sprintf("Some medias were not restored: %v", moved_to))
	}

	return nil
}

func GetUuidMediaName(media *dungeon_models.Media) {
	var media_extension string = filepath.Ext(media.Name)

	uuid_name := uuid.New().String()

	media.Name = uuid_name + media_extension
}

func SetUniqueMediaName(media *dungeon_models.Media, parent_path string) {
	var media_path string = filepath.Join(parent_path, media.Name)
	var media_path_exists bool = dungeon_helpers.FileExists(media_path)

	if !media_path_exists {
		return
	}

	media_name_ext := filepath.Ext(media.Name)
	media_name := strings.TrimSuffix(media.Name, media_name_ext)
	var new_media_name string

	for h := 1; media_path_exists; h++ {
		new_media_name = fmt.Sprintf("%s(%d)%s", media_name, h, media_name_ext)
		media_path = filepath.Join(parent_path, new_media_name)

		media_path_exists = dungeon_helpers.FileExists(media_path)
	}

	media.Name = new_media_name

	return
}

func MoveMediaFile(media *dungeon_models.Media, new_parent_path, current_path string) (lerr *dungeon_models.LabeledError) {

	current_parent_path := dungeon_helpers.GetParentDirectory(current_path)

	is_same_filesystem, err := dungeon_helpers.IsSameFilesystem(new_parent_path, current_parent_path)
	if err != nil {
		lerr = dungeon_models.NewLabeledError(err, "In workflows.MoveMediaFile", dungeon_models.ErrProcessError)
		return lerr
	}

	new_media_path := filepath.Join(new_parent_path, media.Name)

	if !is_same_filesystem {
		err := dungeon_helpers.MoveFile(current_path, new_media_path) // can handle different filesystems but is slower
		if err != nil {
			lerr = dungeon_models.NewLabeledError(err, "In workflows.MoveMediaFile while using dungeon_helpers.MoveFile", dungeon_models.ErrProcessError)
			return lerr
		}
	} else {
		err := os.Rename(new_media_path, current_path) // faster but breaks if paths are on different filesystems
		if err != nil {
			lerr = dungeon_models.NewLabeledError(err, "In workflows.MoveMediaFile while using os.Rename", dungeon_models.ErrProcessError)
			return lerr
		}
	}

	return nil
}
