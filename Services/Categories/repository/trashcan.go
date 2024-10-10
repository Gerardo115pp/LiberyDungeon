package repository

import (
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery_categories_service/models"
)

type TrashRepository interface {
	Commit() error
	DeleteEmptyCategory(category_identity dungeon_models.CategoryIdentity) error
	EmptyTrashcan() error
	EraseTransaction(transaction_id string) (*service_models.TrashcanTransaction, error)
	EraseMediaFromTransaction(transaction_id string, media_uuid string) (*dungeon_models.Media, error)
	GetTransactions() []*service_models.TrashcanTransaction
	GetSortedTransactions() []*service_models.TrashcanTransaction
	GetTransaction(transaction_id string) (*service_models.TrashcanTransaction, error)
	GetSortedTransactionEntries() []*service_models.TrashcanTransactionEntry
	GetTrashcanLocation() string
	GetTrashcanMediaLocation() string
	GetTrashcanSize() int
	HasTransaction(transaction_id string) bool
	MoveToTrash(media_identity *dungeon_models.MediaIdentity) error
	RestoreMediaFromTrash(rejected_media dungeon_models.Media, original_path string) error
	CleanSingleTransaction(transaction_id string) error
	Rollback() error
	StartTransaction(category_identity *dungeon_models.CategoryWeakIdentity) error
	Save() error
	Reload() error
}

var TrashRepo TrashRepository

func SetTrashImplementation(impl TrashRepository) {
	TrashRepo = impl
}
