package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	dungeon_models "libery-dungeon-libs/models"

	"github.com/Gerardo115pp/patriots_lib/echo"
	_ "github.com/go-sql-driver/mysql"
)

type MediasMysql struct {
	db *sql.DB
}

func NewMediasMysql() (*MediasMysql, error) {
	db, err := sql.Open("mysql", createDSN())
	if err != nil {
		return nil, err
	}

	return &MediasMysql{db: db}, nil
}

func (db *MediasMysql) GetRandomMedia(ctx context.Context, cluster_id string, category_id string, only_image bool) (*dungeon_models.Media, *dungeon_models.Category, error) {
	var media dungeon_models.Media = dungeon_models.Media{}
	var category dungeon_models.Category = dungeon_models.Category{}

	sql_query := "SELECT `medias`.`uuid`, `medias`.`name`, `medias`.`last_seen`, `medias`.`main_category`, `medias`.`type`, `medias`.`downloaded_from`, `categorys`.`name` as `category_name`, `categorys`.`fullpath` as `category_fullpath`, `categorys`.`parent` as `category_parent` FROM `medias` LEFT JOIN `categorys` ON `medias`.`main_category` = `categorys`.`uuid`"

	if category_id != "" || only_image {
		sql_query += " WHERE"
	}

	if category_id != "" {
		sql_query += " `categorys`.`uuid` = ?"

		if only_image {
			sql_query += " AND "
		}
	}

	if only_image {
		sql_query += " `medias`.`type` = 'image'"
	}

	sql_query += " ORDER BY -LOG(1.0 - RAND()) / 0.4326 LIMIT 1"

	echo.EchoDebug(fmt.Sprintf("Query: %s", sql_query))

	stmt, err := db.db.Prepare(sql_query)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()

	var random_row *sql.Row

	var downloaded_from sql.NullInt64
	if category_id != "" {
		echo.EchoDebug(fmt.Sprintf("Querying with category_id: %s", category_id))
		random_row = stmt.QueryRowContext(ctx, category_id)
	} else {
		echo.EchoDebug("Querying without category_id")
		random_row = stmt.QueryRowContext(ctx)
	}

	err = random_row.Scan(&media.Uuid, &media.Name, &media.LastSeen, &media.MainCategory, &media.Type, &downloaded_from, &category.Name, &category.Fullpath, &category.Parent)
	if err != nil {
		return nil, nil, err
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	category.Uuid = media.MainCategory
	category.Cluster = cluster_id

	return &media, &category, nil
}

func (db *MediasMysql) GetMediaByID(ctx context.Context, media_id string) (*dungeon_models.Media, error) {
	var media *dungeon_models.Media = new(dungeon_models.Media)

	stmt, err := db.db.Prepare("SELECT `uuid`, `name`, `last_seen`, `main_category`, `media_thumbnail`, `type`, `downloaded_from` FROM `medias` WHERE `uuid` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var downloaded_from sql.NullInt64
	var nullish_media_thumbnail sql.NullString

	err = stmt.QueryRowContext(ctx, media_id).Scan(&media.Uuid, &media.Name, &media.LastSeen, &media.MainCategory, &nullish_media_thumbnail, &media.Type, &downloaded_from)
	if err != nil {
		return nil, err
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	if nullish_media_thumbnail.Valid {
		media.MediaThumbnail = nullish_media_thumbnail.String
	}

	return media, nil
}

func (db *MediasMysql) GetMediaIdentity(ctx context.Context, media_uuid string) (*dungeon_models.MediaIdentity, error) {
	var media_identity *dungeon_models.MediaIdentity

	stmt, err := db.db.PrepareContext(ctx, `
		SELECT
			m.uuid, m.name, m.last_seen, m.main_category, m.media_thumbnail, m.type, m.downloaded_from,
			c.uuid, c.fullpath,
			cc.uuid, cc.fs_path
		FROM medias m
		INNER JOIN categorys c ON m.main_category=c.uuid
		INNER JOIN categories_clusters cc ON c.cluster=cc.uuid
		WHERE m.uuid=?
	`)
	if err != nil {
		return media_identity, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentity: Error preparing statement"))
	}
	defer stmt.Close()

	var time_reciever sql.NullTime
	var downloaded_from_reciever sql.NullInt64

	row := stmt.QueryRowContext(ctx, media_uuid)

	media_identity = new(dungeon_models.MediaIdentity)
	media_identity.Media = new(dungeon_models.Media)

	var media_thumbnail_reciever sql.NullString

	err = row.Scan(
		&media_identity.Media.Uuid,
		&media_identity.Media.Name,
		&time_reciever,
		&media_identity.Media.MainCategory,
		&media_thumbnail_reciever,
		&media_identity.Media.Type,
		&downloaded_from_reciever,
		&media_identity.CategoryUUID,
		&media_identity.CategoryPath,
		&media_identity.ClusterUUID,
		&media_identity.ClusterPath,
	)
	if err != nil {
		return media_identity, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentity: Error scanning row"))
	}

	if time_reciever.Valid {
		media_identity.Media.LastSeen = time_reciever.Time
	}

	media_identity.Media.MediaThumbnail = ""

	if media_thumbnail_reciever.Valid {
		media_identity.Media.MediaThumbnail = media_thumbnail_reciever.String
	}

	if downloaded_from_reciever.Valid {
		media_identity.Media.DownloadedFrom = downloaded_from_reciever.Int64
	}

	return media_identity, nil
}

func (db *MediasMysql) GetMediaByName(ctx context.Context, media_name string, main_category_id string) (*dungeon_models.Media, error) {
	var media *dungeon_models.Media = new(dungeon_models.Media)

	stmt, err := db.db.Prepare("SELECT `uuid`, `name`, `last_seen`, `main_category`, `media_thumbnail`, `type`, `downloaded_from` FROM `medias` WHERE `name` = ? AND `main_category` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var downloaded_from sql.NullInt64
	var nullish_media_thumbnail sql.NullString

	err = stmt.QueryRowContext(ctx, media_name, main_category_id).Scan(&media.Uuid, &media.Name, &media.LastSeen, &media.MainCategory, &nullish_media_thumbnail, &media.Type, &downloaded_from)
	if err != nil {
		return nil, err
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	if nullish_media_thumbnail.Valid {
		media.MediaThumbnail = nullish_media_thumbnail.String
	}

	return media, nil
}

func (db *MediasMysql) InsertMedia(ctx context.Context, media *dungeon_models.Media) error {
	stmt, err := db.db.Prepare("INSERT INTO `medias` (`uuid`, `name`, `last_seen`, `main_category`, `type`, `downloaded_from`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var downloaded_from sql.NullInt64
	if media.DownloadedFrom != 0 {
		downloaded_from.Int64 = media.DownloadedFrom
		downloaded_from.Valid = true
	}

	_, err = stmt.ExecContext(ctx, media.Uuid, media.Name, media.LastSeen, media.MainCategory, media.Type, downloaded_from)
	if err != nil {
		return err
	}

	return nil
}

func (db *MediasMysql) UpdateMediaName(ctx context.Context, media_id string, new_name string) error {
	stmt, err := db.db.Prepare("UPDATE `medias` SET `name` = ? WHERE `uuid` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, new_name, media_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *MediasMysql) Close() error {
	return db.db.Close()
}
