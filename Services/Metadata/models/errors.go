package models

import (
	dungeon_models "libery-dungeon-libs/models"
)

const (
	ErrInvalidWatchPointSize dungeon_models.ErrorLabel = "Invalid watch point size"
	ErrEndOfStream           dungeon_models.ErrorLabel = "End of stream"
	ErrWatchPointNotFound    dungeon_models.ErrorLabel = "Watch point not found"
)
