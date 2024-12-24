package workflows

import (
	"fmt"
	"libery-dungeon-libs/communication"
	dungeon_models "libery-dungeon-libs/models"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

// Cleans resources in other services associated to a list of deleted medias
func ProcessDeletedMedias(medias []dungeon_models.Media) error {

	var media_uuids []string = make([]string, len(medias))

	for h, media := range medias {
		media_uuids[h] = media.Uuid
	}

	_, err := communication.Metadata.RemoveAllTaggingsForEntities(media_uuids)
	if err != nil {
		echo.EchoErr(fmt.Errorf("In workflows/external_resource_cleaners.ProcessDeletedMedias: While calling communication.Metadata.RemoveAllTaggingsForEntities\n\n%s", err))
	}

	return err
}

func ProcessDeletedCategories(categories []dungeon_models.Category) error {
	return nil
}
