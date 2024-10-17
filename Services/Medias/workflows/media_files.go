package workflows

import (
	"context"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"libery_medias_service/repository"
	"os"
	"path/filepath"

	"errors"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func RenameSequence(ctx context.Context, medias []dungeon_models.Media, media_filenames map[string]string) error {
	if len(medias) == 0 {
		return nil
	}

	main_category, err := repository.CategoriesRepo.GetCategoryByID(ctx, medias[0].MainCategory)
	if err != nil {
		errors.Join(err, errors.New("Error getting main category"))
		return err
	}

	category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(ctx, main_category.Cluster)
	if err != nil {
		errors.Join(err, errors.New("Error getting category cluster"))
		return err
	}

	for _, media := range medias {
		if media.MainCategory != main_category.Uuid {
			return errors.New("All medias must share the same main category")
		}

		new_stem, exists := media_filenames[media.Uuid]
		if !exists {
			return fmt.Errorf("Media<%s> with uuid '%s' does not have a new name", media.Name, media.Uuid)
		}

		media_identity := dungeon_models.CreateNewMediaIdentity(&media, &main_category, &category_cluster)

		new_media_name := dungeon_helpers.RenameFilename(media_identity.Media.Name, new_stem)
		fmt.Printf("New media<%s> name: %s\n", media.Name, new_media_name)

		if new_media_name == media.Name {
			echo.Echo(echo.YellowFG, fmt.Sprintf("Media<%s> with uuid '%s' already has the new name: '%s'", media.Name, media.Uuid, new_media_name))
			continue
		}

		err = RenameMedia(ctx, *media_identity, new_media_name)
		if err != nil {
			errors.Join(err, fmt.Errorf("Error renaming media with uuid '%s'", media.Uuid))
			return err
		}

		echo.Echo(echo.GreenFG, fmt.Sprintf("Media<%s> with uuid '%s' renamed to '%s'", media.Name, media.Uuid, new_media_name))
	}

	return nil
}

func RenameMedia(ctx context.Context, media_identity dungeon_models.MediaIdentity, new_name string) error {
	category_path := filepath.Join(media_identity.ClusterPath, media_identity.CategoryPath)

	abs_new_name := filepath.Join(category_path, new_name)
	abs_current_name := filepath.Join(category_path, media_identity.Media.Name)

	if dungeon_helpers.FileExists(abs_new_name) {
		name_holder_media, err := repository.MediasRepo.GetMediaByName(ctx, new_name, media_identity.CategoryUUID)
		if err != nil {
			return fmt.Errorf("Error getting media with name '%s'", new_name)
		}

		name_holder_identity := dungeon_models.MediaIdentity{
			Media:        name_holder_media,
			CategoryUUID: media_identity.CategoryUUID,
			ClusterPath:  media_identity.ClusterPath,
			CategoryPath: media_identity.CategoryPath,
			ClusterUUID:  media_identity.ClusterUUID,
		}

		name_holder_new_name, err := GetUniquieMediaName(name_holder_media.Name, media_identity.CategoryPath, media_identity.ClusterPath)
		if err != nil {
			return fmt.Errorf("Error getting unique media name: %s", err.Error())
		}

		err = RenameMedia(ctx, name_holder_identity, name_holder_new_name)
		if err != nil {
			return fmt.Errorf("Error renaming media with name '%s'", new_name)
		}

	}

	err := os.Rename(abs_current_name, abs_new_name)
	if err != nil {
		return fmt.Errorf("Error renaming media file `%s` to `%s`: %s", abs_current_name, abs_new_name, err.Error())
	}

	err = repository.MediasRepo.UpdateMediaName(ctx, media_identity.Media.Uuid, new_name)
	if err != nil {
		return fmt.Errorf("Error updating media name in database: %s", err.Error())
	}

	return nil
}
