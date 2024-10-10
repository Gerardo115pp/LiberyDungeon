package workflows

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery_categories_service/models"
	"libery_categories_service/repository"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func RestoreTransactionToCategory(transaction service_models.TrashcanTransaction, category_identity *dungeon_models.CategoryIdentity) (lerr *dungeon_models.LabeledError) {
	var media_uuids []string = make([]string, 0)

	for _, media := range transaction.Content {
		media_uuids = append(media_uuids, media.Uuid)
	}

	_, err := repository.TrashRepo.EraseTransaction(transaction.TransactionID)
	if err != nil {
		lerr = dungeon_models.NewLabeledError(err, "In workflows.RestoreTransactionToCategory", dungeon_models.ErrDB_CouldNotCreateTX)
		return
	}

	for _, media_uuid := range media_uuids {
		lerr = RestoreMediaFromTransaction(*category_identity, transaction, media_uuid)
		if lerr != nil {
			lerr.AppendContext("In workflows.RestoreTransactionToCategory")
			return
		}
	}

	return
}

func RestoreMediaFromTransaction(parent_identity dungeon_models.CategoryIdentity, transaction service_models.TrashcanTransaction, media_uuid string) (lerr *dungeon_models.LabeledError) {
	var err error
	media := transaction.GetMediaByUuid(media_uuid)
	if media == nil {
		echo.Echo(echo.RedFG, "In workflows.RestoreMediaFromTransaction: Media not found in transaction")
		lerr = dungeon_models.NewLabeledError(nil, "In workflows.RestoreMediaFromTransaction", dungeon_models.ErrProcessError)
		return
	}

	new_parent_path := filepath.Join(parent_identity.ClusterPath, parent_identity.Category.Fullpath)
	current_parent_path := filepath.Join(repository.TrashRepo.GetTrashcanMediaLocation(), media.Name)

	media.MainCategory = parent_identity.Category.Uuid

	SetUniqueMediaName(media, new_parent_path)

	lerr = MoveMediaFile(media, new_parent_path, current_parent_path)
	if lerr != nil {
		echo.Echo(echo.RedFG, "In workflows.RestoreMediaFromTransaction: Error moving media file")
		lerr.AppendContext("In workflows.RestoreMediaFromTransaction")
		return
	}

	err = repository.MediasRepo.InsertMedia(context.Background(), media)
	if err != nil {
		echo.Echo(echo.RedFG, "In workflows.RestoreMediaFromTransaction: Error inserting media")
		lerr = dungeon_models.NewLabeledError(err, "In workflows.RestoreMediaTransaction", dungeon_models.ErrDB_CouldNotConnectToDB)
		return
	}

	return
}
