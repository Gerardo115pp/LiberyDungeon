package repository

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
)

type MediasRepository interface {
	GetMedia(ctx context.Context, media_uuid string) (*dungeon_models.Media, error)
	GetMedias(ctx context.Context, media_uuids []string) ([]*dungeon_models.Media, error)
	DeleteMedia(ctx context.Context, media_uuid string) error
	InsertMedia(ctx context.Context, media *dungeon_models.Media) error
}

// This is only meant to be used by the cluster creation process. As making an API call for
// each media file would be too slow for an entire cluster of medias. Do not use this repo outside that
// context.
var MediasRepo MediasRepository

func SetMediasImplementation(repo MediasRepository) {
	MediasRepo = repo
}
