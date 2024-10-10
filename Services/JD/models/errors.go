package models

import dungeon_models "libery-dungeon-libs/models"

var (
	ErrService_NotSuchService        dungeon_models.ErrorLabel = "Service not found"
	ErrService_NoSuchEventRegistered dungeon_models.ErrorLabel = "An event with the given UUID is not registered"
)
