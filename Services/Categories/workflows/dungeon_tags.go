package workflows

import (
	"errors"
	"fmt"
	"libery-dungeon-libs/communication"
	dungeon_models "libery-dungeon-libs/models"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

// Receives a map of category_uuid -> medias, and applies all dungeon tags of the category to the media list.
func ApplyCategoryTags(new_categories_medias map[string][]dungeon_models.Media, medias_cluster *dungeon_models.CategoryCluster) error {

	start_time := time.Now()

	for category_uuid, media_list := range new_categories_medias {
		medias_uuids := make([]string, len(media_list))
		for h, media := range media_list {
			medias_uuids[h] = media.Uuid
		}

		_, err := communication.Metadata.CopyEntityTagsToEntities(category_uuid, medias_cluster.Uuid, dungeon_models.ENTITY_TYPE_MEDIA, medias_uuids)
		if err != nil {
			return errors.Join(fmt.Errorf("In workflows/dungeons_tags.ApplyCategoryTags: Couldn't copy tags from category<%s> to medias", category_uuid), err)
		}
	}

	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Function took %s", time.Since(start_time)))

	return nil
}
