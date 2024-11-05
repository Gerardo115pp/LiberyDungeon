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
	DeleteTagCTX(ctx context.Context, tag_id int) error
	DeleteTag(tag_id int) error
	GetGlobalTaxonomiesCTX(ctx context.Context) ([]service_models.TagTaxonomy, error)
	GetGlobalTaxonomies() ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomiesCTX(ctx context.Context, cluster_uuid string) ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomies(cluster_uuid string) ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomiesByInternalValueCTX(ctx context.Context, cluster_uuid string, internal bool) ([]service_models.TagTaxonomy, error)
	GetClusterTaxonomiesByInternalValue(cluster_uuid string, internal bool) ([]service_models.TagTaxonomy, error)
	GetClusterTagsCTX(ctx context.Context, cluster_uuid string) ([]service_models.TaxonomyTags, error)
	GetClusterTags(cluster_uuid string) ([]service_models.TaxonomyTags, error)
	GetClusterTagsByInternalValueCTX(ctx context.Context, cluster_uuid string, internal bool) ([]service_models.TaxonomyTags, error)
	GetClusterTagsByInternalValue(cluster_uuid string, internal bool) ([]service_models.TaxonomyTags, error)
	GetTagByIdCTX(ctx context.Context, tag_id int) (service_models.DungeonTag, error)
	GetTagById(tag_id int) (service_models.DungeonTag, error)
	GetTagByNameCTX(ctx context.Context, tag_name, taxonomy string) (service_models.DungeonTag, error)
	GetTagByName(tag_name, taxonomy string) (service_models.DungeonTag, error)
	GetTaxonomyTagsCTX(ctx context.Context, taxonomy_uuid string) ([]service_models.DungeonTag, error)
	GetTaxonomyTags(taxonomy_uuid string) ([]service_models.DungeonTag, error)
	GetTagTaxonomyCTX(ctx context.Context, taxonomy_uuid string) (service_models.TagTaxonomy, error)
	GetTagTaxonomy(taxonomy_uuid string) (service_models.TagTaxonomy, error)
	GetEntityTaggingsCTX(ctx context.Context, entity_uuid, cluster_domain string) ([]service_models.DungeonTagging, error)
	GetEntityTaggings(entity_uuid, cluster_domain string) ([]service_models.DungeonTagging, error)
	GetEntitiesWithTaggingsCTX(ctx context.Context, tags []int) ([]string, error)
	GetEntitiesWithTaggings(tags []int) ([]string, error)
	RemoveTagFromEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error
	RemoveTagFromEntity(tag_id int, entity_uuid string) error
	TagEntityCTX(ctx context.Context, tag_id int, entity_uuid, entity_type string) (int64, error)
	TagEntity(tag_id int, entity_uuid, entity_type string) (int64, error)
	TagEntitiesCTX(ctx context.Context, tag_id int, entities_uuids []string, entity_type string) error
	TagEntities(tag_id int, entities_uuids []string, entity_type string) error
	UpdateTaxonomyNameCTX(ctx context.Context, taxonomy_uuid, new_name string) error
	UpdateTaxonomyName(taxonomy_uuid, new_name string) error
	UpdateTagNameCTX(ctx context.Context, tag_id int, new_name string) error
	UpdateTagName(tag_id int, new_name string) error
}

var DungeonTagsRepo DungeonTagsRepository

func SetDungeonTagsRepository(repo DungeonTagsRepository) {
	DungeonTagsRepo = repo
}
