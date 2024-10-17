package database

import (
	"context"
	"database/sql"
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

	stmt, err := db.db.Prepare("SELECT `uuid`, `name`, `last_seen`, `main_category`, `type`, `downloaded_from` FROM `medias` WHERE `uuid` = ?")
	if err != nil {
		return nil, err
	}

	var downloaded_from sql.NullInt64

	err = stmt.QueryRowContext(ctx, media_id).Scan(&media.Uuid, &media.Name, &media.LastSeen, &media.MainCategory, &media.Type, &downloaded_from)
	if err != nil {
		return nil, err
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	return media, nil
}

func (db *MediasMysql) GetMediaByName(ctx context.Context, media_name string, main_category_id string) (*dungeon_models.Media, error) {
	var media *dungeon_models.Media = new(dungeon_models.Media)

	stmt, err := db.db.Prepare("SELECT `uuid`, `name`, `last_seen`, `main_category`, `type`, `downloaded_from` FROM `medias` WHERE `name` = ? AND `main_category` = ?")
	if err != nil {
		return nil, err
	}

	var downloaded_from sql.NullInt64

	err = stmt.QueryRowContext(ctx, media_name, main_category_id).Scan(&media.Uuid, &media.Name, &media.LastSeen, &media.MainCategory, &media.Type, &downloaded_from)
	if err != nil {
		return nil, err
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	return media, nil
}

func (db *MediasMysql) InsertMedia(ctx context.Context, media *dungeon_models.Media) error {
	stmt, err := db.db.Prepare("INSERT INTO medias (uuid, name, last_seen, main_category, type, downloaded_from) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

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

func (db *MediasMysql) Close() error {
	return db.db.Close()
}
