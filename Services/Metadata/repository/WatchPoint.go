package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
	"libery-metadata-service/models"
)

type WatchPointRepository interface {
	InsertWatchPoint(ctx context.Context, media_uuid string, start_time uint32) *dungeon_models.LabeledError
	GetWatchPointByMediaID(ctx context.Context, media_uuid string) (*models.WatchPoint, *dungeon_models.LabeledError)
	GetWatchPointByORD(ctx context.Context, ord uint16) (*models.WatchPoint, *dungeon_models.LabeledError)
	Close() error
}

var WatchPointRepo WatchPointRepository

func SetWatchPointRepository(repo WatchPointRepository) {
	WatchPointRepo = repo
}
