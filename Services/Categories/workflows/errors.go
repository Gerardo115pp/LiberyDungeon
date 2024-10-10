package workflows

import dungeon_models "libery-dungeon-libs/models"

const (
	ErrForbiddenDirectoryScan                dungeon_models.ErrorLabel = "Forbidden directory scan"
	ErrNoSuchFileOrDirectory                 dungeon_models.ErrorLabel = "No such file or directory"
	ErrNoSuchDirectory                       dungeon_models.ErrorLabel = "A path does not exist or is not a directory"
	ErrPath_AlreadyExists                    dungeon_models.ErrorLabel = "Path already exists"
	ErrPathValidation_PathIsACluster         dungeon_models.ErrorLabel = "Path is a cluster"
	ErrPathValidation_PathIsAClusterAncestor dungeon_models.ErrorLabel = "Path is a cluster ancestor"
	ErrPathValidation_PathIsAClusterChild    dungeon_models.ErrorLabel = "Path is a cluster child"
)
