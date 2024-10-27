package models

import dungeon_helpers "libery-dungeon-libs/helpers"

type TagTaxonomy struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	ClusterDomain string `json:"cluster_domain"` // If not an empty string, defines the cluster where the tags on this taxonomy are applied
	IsInternal    bool   `json:"is_internal"`    // If true, this taxonomy is only used internally by the system and not exposed to the user
}

type TaxonomyTags struct {
	Taxonomy *TagTaxonomy `json:"taxonomy"`
	Tags     []DungeonTag `json:"tags"`
}

type DungeonTag struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Taxonomy     string `json:"taxonomy"`      // The UUID of the taxonomy the tag belongs to
	NameTaxonomy string `json:"name_taxonomy"` // A hash of the Name and Taxonomy fields. Meant to keep name uniqueness within a taxonomy
}

func (dt *DungeonTag) RecalculateNameTaxonomy() {
	new_source := dt.Name + dt.Taxonomy

	dt.NameTaxonomy = dungeon_helpers.GenerateSha1ID(new_source)
}

type DungeonTagging struct {
	TaggingID        int64       `json:"tagging_id"`
	EntityType       string      `json:"entity_type"`
	Tag              *DungeonTag `json:"tag"`
	TaggedEntityUUID string      `json:"tagged_entity_uuid"`
}
