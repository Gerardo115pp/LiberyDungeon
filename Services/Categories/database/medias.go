package database

import (
	"context"
	"database/sql"
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
