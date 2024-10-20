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
	DeleteTaxonomyCTX(ctx context.Context, taxonomy_uuid string) error
	DeleteTaxonomy(taxonomy_uuid string) error
	GetGlobalTaxonomiesCTX(ctx context.Context) ([]service_models.TagTaxonomy, error)
	GetGlobalTaxonomies() ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomiesCTX(ctx context.Context, cluster_uuid string) ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomies(cluster_uuid string) ([]service_models.TagTaxonomy, error)
	GetClusterTagsCTX(ctx context.Context, cluster_uuid string) ([]service_models.TaxonomyTags, error)
	GetClusterTags(cluster_uuid string) ([]service_models.TaxonomyTags, error)
	GetTagByIdCTX(ctx context.Context, tag_id int) (service_models.DungeonTag, error)
	GetTagById(tag_id int) (service_models.DungeonTag, error)
	GetTaxonomyTagsCTX(ctx context.Context, taxonomy_uuid string) ([]service_models.DungeonTag, error)
	GetTaxonomyTags(taxonomy_uuid string) ([]service_models.DungeonTag, error)
	GetEntityTaggingsCTX(ctx context.Context, entity_uuid string) ([]service_models.DungeonTagging, error)
	GetEntityTaggings(entity_uuid string) ([]service_models.DungeonTagging, error)
	RemoveTagFromEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error
	RemoveTagFromEntity(tag_id int, entity_uuid string) error
	TagEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error
	TagEntity(tag_id int, entity_uuid string) error
	UpdateTaxonomyNameCTX(ctx context.Context, taxonomy_uuid, new_name string) error
	UpdateTaxonomyName(taxonomy_uuid, new_name string) error
	UpdateTagNameCTX(ctx context.Context, tag_id int, new_name string) error
	UpdateTagName(tag_id int, new_name string) error
}

var DungeonTagsRepo DungeonTagsRepository

func SetDungeonTagsRepository(repo DungeonTagsRepository) {
	DungeonTagsRepo = repo
}
