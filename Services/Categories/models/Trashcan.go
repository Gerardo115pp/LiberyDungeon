package models

import dungeon_models "libery-dungeon-libs/models"

type TrashcanTransaction struct {
	TransactionID  string                              `json:"transaction_id"`
	Content        []dungeon_models.Media              `json:"content"`
	OriginIdentity dungeon_models.CategoryWeakIdentity `json:"origin_identity"`
}

type TrashcanTransactionEntry struct {
	Timestamp      string `json:"timestamp"`
	AffectedMedias int    `json:"affected_medias"`
}

func NewTrashcanTransaction(transaction_id string, identity dungeon_models.CategoryWeakIdentity) *TrashcanTransaction {
	var new_transaction *TrashcanTransaction = new(TrashcanTransaction)

	new_transaction.TransactionID = transaction_id
	new_transaction.OriginIdentity = identity
	new_transaction.Content = make([]dungeon_models.Media, 0)

	return new_transaction
}

func (tt *TrashcanTransaction) AddMedia(media dungeon_models.Media) {
	tt.Content = append(tt.Content, media)
}

func (tt *TrashcanTransaction) ToEntry() *TrashcanTransactionEntry {
	return &TrashcanTransactionEntry{
		Timestamp:      tt.TransactionID,
		AffectedMedias: len(tt.Content),
	}
}

func (tt *TrashcanTransaction) GetMediaByUuid(media_uuid string) *dungeon_models.Media {
	for _, media := range tt.Content {
		if media.Uuid == media_uuid {
			return &media
		}
	}

	return nil
}
