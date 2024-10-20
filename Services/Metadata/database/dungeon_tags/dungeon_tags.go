package dungeon_tags

import (
	"context"
	"database/sql"
	"libery-dungeon-libs/libs/dungeon_sqlite_opener"
	app_config "libery-metadata-service/Config"
	service_models "libery-metadata-service/models"
)

type DungeonTagsDB struct {
	db_conn *sql.DB
}

func NewDungeonTagsDB() *DungeonTagsDB {
	var dungeon_tags_dbd *DungeonTagsDB = new(DungeonTagsDB)

	var sqlite_opener *dungeon_sqlite_opener.DungeonSqliteOpener
	sqlite_opener = dungeon_sqlite_opener.NewDungeonSqliteOpener("dungeon_tags.db", "dungeon_tags.sql", app_config.OPERATION_DATA_PATH)

	db, err := sqlite_opener.OpenDB(true)
	if err != nil {
		panic(err)
	}

	dungeon_tags_dbd.db_conn = db

	return dungeon_tags_dbd
}

func (dt_db *DungeonTagsDB) CreateTaxonomyCTX(ctx context.Context, taxonomy *service_models.TagTaxonomy) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "INSERT INTO `tag_taxonomies`(`uuid`, `name`, `cluster_domain`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, taxonomy.UUID, taxonomy.Name, taxonomy.ClusterDomain)

	return err
}

func (dt_db *DungeonTagsDB) CreateTaxonomy(taxonomy *service_models.TagTaxonomy) error {
	return dt_db.CreateTaxonomyCTX(context.Background(), taxonomy)
}

func (dt_db *DungeonTagsDB) CreateTagCTX(ctx context.Context, tag *service_models.DungeonTag) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "INSERT INTO `dungeon_tags`(`name`, `taxonomy`, `name_taxonomy`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, tag.Name, tag.Taxonomy, tag.NameTaxonomy)

	return err
}

func (dt_db *DungeonTagsDB) CreateTag(tag *service_models.DungeonTag) error {
	return dt_db.CreateTagCTX(context.Background(), tag)
}

func (dt_db *DungeonTagsDB) DeleteTaxonomyCTX(ctx context.Context, taxonomy_uuid string) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "DELETE FROM `tag_taxonomies` WHERE `uuid`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, taxonomy_uuid)

	return err
}

func (dt_db *DungeonTagsDB) DeleteTaxonomy(taxonomy_uuid string) error {
	return dt_db.DeleteTaxonomyCTX(context.Background(), taxonomy_uuid)
}

func (dt_db *DungeonTagsDB) DeleteTagCTX(ctx context.Context, tag_id int) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "DELETE FROM `dungeon_tags` WHERE `id`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, tag_id)

	return err
}

func (dt_db *DungeonTagsDB) DeleteTag(tag_id int) error {
	return dt_db.DeleteTagCTX(context.Background(), tag_id)
}

func (dt_db *DungeonTagsDB) GetGlobalTaxonomiesCTX(ctx context.Context) ([]service_models.TagTaxonomy, error) {
	var taxonomies []service_models.TagTaxonomy

	rows, err := dt_db.db_conn.QueryContext(ctx, "SELECT uuid, name WHERE cluster_domain=''")
	if err != nil {
		return taxonomies, err
	}

	for rows.Next() {
		var taxonomy service_models.TagTaxonomy
		err = rows.Scan(&taxonomy.UUID, &taxonomy.Name)
		if err != nil {
			return taxonomies, err
		}

		taxonomies = append(taxonomies, taxonomy)
	}

	return taxonomies, nil
}

func (dt_db *DungeonTagsDB) GetGlobalTaxonomies() ([]service_models.TagTaxonomy, error) {
	return dt_db.GetGlobalTaxonomiesCTX(context.Background())
}

func (dt_db *DungeonTagsDB) GetClusterTaxonomiesCTX(ctx context.Context, cluster_uuid string) ([]service_models.TagTaxonomy, error) {
	var taxonomies []service_models.TagTaxonomy

	rows, err := dt_db.db_conn.QueryContext(ctx, "SELECT `uuid`, `name`, `cluster_domain` WHERE `cluster_domain`=?", cluster_uuid)
	if err != nil {
		return taxonomies, err
	}

	for rows.Next() {
		var taxonomy service_models.TagTaxonomy
		err = rows.Scan(&taxonomy.UUID, &taxonomy.Name, &taxonomy.ClusterDomain)
		if err != nil {
			return taxonomies, err
		}

		taxonomies = append(taxonomies, taxonomy)
	}

	return taxonomies, nil
}

func (dt_db *DungeonTagsDB) GetClusterTaxonomies(cluster_uuid string) ([]service_models.TagTaxonomy, error) {
	return dt_db.GetClusterTaxonomiesCTX(context.Background(), cluster_uuid)
}

func (db_db *DungeonTagsDB) GetClusterTagsCTX(ctx context.Context, cluster_uuid string) ([]service_models.TaxonomyTags, error) {
	var taxonomies_tags []service_models.TaxonomyTags = make([]service_models.TaxonomyTags, 0)

	taxonomies, err := db_db.GetClusterTaxonomiesCTX(ctx, cluster_uuid)
	if err != nil {
		return taxonomies_tags, err
	}

	for _, taxonomy := range taxonomies {
		var taxonomy_tags service_models.TaxonomyTags = service_models.TaxonomyTags{Taxonomy: &taxonomy}

		tags, err := db_db.GetTaxonomyTagsCTX(ctx, taxonomy.UUID)
		if err != nil {
			return taxonomies_tags, err
		}

		taxonomy_tags.Tags = tags

		taxonomies_tags = append(taxonomies_tags, taxonomy_tags)
	}

	return taxonomies_tags, nil
}

func (db_db *DungeonTagsDB) GetClusterTags(cluster_uuid string) ([]service_models.TaxonomyTags, error) {
	return db_db.GetClusterTagsCTX(context.Background(), cluster_uuid)
}

func (dt_db *DungeonTagsDB) GetTagByIdCTX(ctx context.Context, tag_id int) (service_models.DungeonTag, error) {
	var tag service_models.DungeonTag

	stmt, err := dt_db.db_conn.PrepareContext(ctx, "SELECT `id`, `name`, `taxonomy`, `name_taxonomy FROM `dungeon_tags` WHERE `id`=?")
	if err != nil {
		return tag, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, tag_id).Scan(&tag.ID, &tag.Name, &tag.Taxonomy, &tag.NameTaxonomy)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func (dt_db *DungeonTagsDB) GetTagById(tag_id int) (service_models.DungeonTag, error) {
	return dt_db.GetTagByIdCTX(context.Background(), tag_id)
}

func (dt_db *DungeonTagsDB) GetTagByNameCTX(ctx context.Context, tag_name, taxonomy string) (service_models.DungeonTag, error) {
	var tag service_models.DungeonTag

	stmt, err := dt_db.db_conn.PrepareContext(ctx, "SELECT `id`, `name`, `taxonomy`, `name_taxonomy FROM `dungeon_tags` WHERE `name`=? AND `taxonomy`=?")
	if err != nil {
		return tag, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, tag_name, taxonomy).Scan(&tag.ID, &tag.Name, &tag.Taxonomy, &tag.NameTaxonomy)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func (dt_db *DungeonTagsDB) GetTagByName(tag_name, taxonomy string) (service_models.DungeonTag, error) {
	return dt_db.GetTagByNameCTX(context.Background(), tag_name, taxonomy)
}

func (dt_db *DungeonTagsDB) GetTaxonomyTagsCTX(ctx context.Context, taxonomy_uuid string) ([]service_models.DungeonTag, error) {
	var tags []service_models.DungeonTag = make([]service_models.DungeonTag, 0)

	rows, err := dt_db.db_conn.QueryContext(ctx, "SELECT `id`, `name`, `taxonomy`, `name_taxonomy` WHERE `taxonomy`=?", taxonomy_uuid)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag service_models.DungeonTag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Taxonomy, &tag.NameTaxonomy)
		if err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (dt_db *DungeonTagsDB) GetTaxonomyTags(taxonomy_uuid string) ([]service_models.DungeonTag, error) {
	return dt_db.GetTaxonomyTagsCTX(context.Background(), taxonomy_uuid)
}

func (dt_db *DungeonTagsDB) GetEntityTaggingsCTX(ctx context.Context, entity_uuid string) ([]service_models.DungeonTagging, error) {
	var taggings []service_models.DungeonTagging = make([]service_models.DungeonTagging, 0)

	stmt, err := dt_db.db_conn.PrepareContext(ctx, "SELECT `tagging_id`, `tag`, `taggable_id` FROM `taggings` WHERE `taggable_id`=?")
	if err != nil {
		return taggings, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, entity_uuid)
	if err != nil {
		return taggings, err
	}

	for rows.Next() {
		var tagging service_models.DungeonTagging
		var tag_id int
		err = rows.Scan(&tagging.TaggingID, &tag_id, &tagging.TaggedEntityUUID)
		if err != nil {
			return taggings, err
		}

		var tag service_models.DungeonTag

		tag, err = dt_db.GetTagByIdCTX(ctx, tag_id)
		if err != nil {
			return taggings, err
		}

		tagging.Tag = &tag

		taggings = append(taggings, tagging)
	}

	return taggings, nil
}

func (dt_db *DungeonTagsDB) GetEntityTaggings(entity_uuid string) ([]service_models.DungeonTagging, error) {
	return dt_db.GetEntityTaggingsCTX(context.Background(), entity_uuid)
}

func (dt_db *DungeonTagsDB) RemoveTagFromEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "DELETE FROM `taggings` WHERE `tag`=? AND `taggable_id`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, tag_id, entity_uuid)

	return err
}

func (dt_db *DungeonTagsDB) RemoveTagFromEntity(tag_id int, entity_uuid string) error {
	return dt_db.RemoveTagFromEntityCTX(context.Background(), tag_id, entity_uuid)
}

func (dt_db *DungeonTagsDB) TagEntityCTX(ctx context.Context, tag_id int, entity_uuid string) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "INSERT INTO `taggings`(`tag_id`, `tagged_entity_uuid`) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, tag_id, entity_uuid)

	return err
}

func (dt_db *DungeonTagsDB) TagEntity(tag_id int, entity_uuid string) error {
	return dt_db.TagEntityCTX(context.Background(), tag_id, entity_uuid)
}

func (dt_db *DungeonTagsDB) UpdateTaxonomyNameCTX(ctx context.Context, taxonomy_uuid, new_name string) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "UPDATE `tag_taxonomies` SET `name`=? WHERE `uuid`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, new_name, taxonomy_uuid)

	return err
}

func (dt_db *DungeonTagsDB) UpdateTaxonomyName(taxonomy_uuid, new_name string) error {
	return dt_db.UpdateTaxonomyNameCTX(context.Background(), taxonomy_uuid, new_name)
}

func (dt_db *DungeonTagsDB) UpdateTagNameCTX(ctx context.Context, tag_id int, new_name string) error {
	stmt, err := dt_db.db_conn.PrepareContext(ctx, "UPDATE `dungeon_tags` SET `name`=? WHERE `id`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, new_name, tag_id)

	return err
}

func (dt_db *DungeonTagsDB) UpdateTagName(tag_id int, new_name string) error {
	return dt_db.UpdateTagNameCTX(context.Background(), tag_id, new_name)
}
