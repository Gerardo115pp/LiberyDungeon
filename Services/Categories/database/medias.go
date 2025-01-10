package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"

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

func (medias_repo *MediasMysql) GetMedia(ctx context.Context, media_uuid string) (*dungeon_models.Media, error) {
	stmt, err := medias_repo.db.PrepareContext(ctx, "SELECT `uuid`, `name`, `last_seen`, `main_category`, `media_thumbnail`, `type`, `downloaded_from` FROM `medias` WHERE `uuid` = ?")
	if err != nil {
		return nil, errors.Join(errors.New("In database/medias.GetMedia: while preparing statement: "), err)
	}
	defer stmt.Close()

	var media *dungeon_models.Media = new(dungeon_models.Media)

	var downloaded_from sql.NullInt64
	var time_reciever sql.NullTime
	var media_thumbnail_reciever sql.NullString

	err = stmt.QueryRowContext(ctx, media_uuid).Scan(&media.Uuid, &media.Name, &time_reciever, &media.MainCategory, &media_thumbnail_reciever, &media.Type, &downloaded_from)
	if err != nil {
		return nil, errors.Join(errors.New("In database/medias.GetMedia: while scanning row: "), err)
	}

	if time_reciever.Valid {
		media.LastSeen = time_reciever.Time
	}

	if media_thumbnail_reciever.Valid {
		media.MediaThumbnail = media_thumbnail_reciever.String
	}

	if downloaded_from.Valid {
		media.DownloadedFrom = downloaded_from.Int64
	}

	return media, nil
}

func (medias_repo *MediasMysql) GetMedias(ctx context.Context, media_uuids []string) ([]*dungeon_models.Media, error) {
	if len(media_uuids) < 1 {
		return make([]*dungeon_models.Media, 0), nil
	}

	var stmt_placeholderrs string = dungeon_helpers.GetPreparedListPlaceholders(len(media_uuids))

	stmt, err := medias_repo.db.PrepareContext(ctx, fmt.Sprintf("SELECT `uuid`, `name`, `last_seen`, `main_category`, `media_thumbnail`, `type`, `downloaded_from` FROM `medias` WHERE `uuid` IN (%s)", stmt_placeholderrs))
	if err != nil {
		return nil, errors.Join(errors.New("In database/medias.GetMedias: while preparing statement: "), err)
	}
	defer stmt.Close()

	var medias []*dungeon_models.Media = make([]*dungeon_models.Media, len(media_uuids))

	var downloaded_from sql.NullInt64
	var time_reciever sql.NullTime
	var media_thumbnail_reciever sql.NullString

	args := make([]interface{}, len(media_uuids))
	for h, media_uuid := range media_uuids {
		args[h] = media_uuid
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Join(errors.New("In database/medias.GetMedias: while querying rows: "), err)
	}
	defer rows.Close()

	for h := 0; rows.Next(); h++ {
		var new_media *dungeon_models.Media = new(dungeon_models.Media)

		err = rows.Scan(&new_media.Uuid, &new_media.Name, &time_reciever, &new_media.MainCategory, &media_thumbnail_reciever, &new_media.Type, &downloaded_from)
		if err != nil {
			return nil, errors.Join(errors.New("In database/medias.GetMedias: while scanning row: "), err)
		}

		if time_reciever.Valid {
			new_media.LastSeen = time_reciever.Time
		}

		if media_thumbnail_reciever.Valid {
			new_media.MediaThumbnail = media_thumbnail_reciever.String
		}

		if downloaded_from.Valid {
			new_media.DownloadedFrom = downloaded_from.Int64
		}

		medias[h] = new_media
	}

	return medias, nil
}

func (medias_repo *MediasMysql) InsertMedia(ctx context.Context, media *dungeon_models.Media) error {
	stmt, err := medias_repo.db.PrepareContext(ctx, "INSERT INTO medias (uuid, name, last_seen, main_category, type, downloaded_from) VALUES (?, ?, ?, ?, ?, ?)")
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

func (medias_repo *MediasMysql) DeleteMedia(ctx context.Context, media_uuid string) error {
	stmt, err := medias_repo.db.PrepareContext(ctx, "DELETE FROM medias WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, media_uuid)
	if err != nil {
		return err
	}

	return nil
}
