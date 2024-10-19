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
