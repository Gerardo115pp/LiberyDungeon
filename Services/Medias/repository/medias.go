package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type MediasRepository interface {
	InsertMedia(ctx context.Context, media *dungeon_models.Media) error
	GetMediaByID(ctx context.Context, media_id string) (*dungeon_models.Media, error)
	GetMediaByName(ctx context.Context, media_name string, main_category_id string) (*dungeon_models.Media, error)
	GetRandomMedia(ctx context.Context, cluster_id string, category_id string, only_image bool) (*dungeon_models.Media, *dungeon_models.Category, error)
	Close() error
}

var MediasRepo MediasRepository

func SetMediasImplementation(impl MediasRepository) {
	MediasRepo = impl
}
