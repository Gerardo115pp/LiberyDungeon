package repository

import (
	"context"
	service_models "libery-metadata-service/models"
)

type DungeonTagsRepository interface {
	CreateTaxonomyCTX(ctx context.Context, taxonomy *service_models.TagTaxonomy) error
	CreateTaxonomy(taxonomy *service_models.TagTaxonomy) error
	CreateTagCTX(ctx context.Context, tag *service_models.DungeonTag) error
	CreateTag(tag *service_models.DungeonTag) error
	TagEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error
	TagEntity(tag_id int, entity_uuid string) error
}

var DungeonTagsRepo DungeonTagsRepository

func SetDungeonTagsRepository(repo DungeonTagsRepository) {
	DungeonTagsRepo = repo
}
