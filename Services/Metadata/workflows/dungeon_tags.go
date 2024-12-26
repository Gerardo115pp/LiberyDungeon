package workflows

import (
	"context"
	"errors"
	"fmt"
	"libery-metadata-service/repository"
)

func CopyEntityTagsToTagListCTX(ctx context.Context, entity_uuid string, entities_uuids []string, cluster_domain, entities_type string) error {
	entity_taggings, err := repository.DungeonTagsRepo.GetEntityTaggingsCTX(ctx, entity_uuid, cluster_domain)
	if err != nil {
		return errors.Join(fmt.Errorf("In workflows/dungeons_tags.CopyEntityTagsToTagList: Couldn't get entity<%s> taggings in domain '%s'", entity_uuid, cluster_domain), err)
	}

	for _, tagging := range entity_taggings {
		err := repository.DungeonTagsRepo.TagEntitiesCTX(ctx, int(tagging.Tag.ID), entities_uuids, entities_type)
		if err != nil {
			return errors.Join(fmt.Errorf("In workflows/dungeons_tags.CopyEntityTagsToTagList: Couldn't tag entities with tag<%d> in domain '%s'", tagging.Tag.ID, cluster_domain), err)
		}
	}

	return nil
}
