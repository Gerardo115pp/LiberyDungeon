package repository

import (
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery_JD_service/models"
)

type PlatformFeatures interface {
	SetServiceOnline(service_name string, online bool) (lerr *dungeon_models.LabeledError)
	IsServiceOnline(service_name string) bool
	SetServiceMetadata(new_service service_models.PlatformService) (lerr *dungeon_models.LabeledError)
}

var PlatformFeaturesRepo PlatformFeatures

func SetPlatformFeaturesRepo(repo PlatformFeatures) {
	PlatformFeaturesRepo = repo
}
