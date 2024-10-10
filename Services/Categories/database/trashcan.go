package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	service_helpers "libery_categories_service/helpers"
	service_models "libery_categories_service/models"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

var JOURNAL_NAME string = "journal.json"
var CATEGORIES_JOURNAL_NAME string = "categories_journal.json"
var MEDIAS_STORAGE_DIRECTORY string = "medias"
var datetime_format string = "2006-01-02 15:04:05"

type TrashcanDatabase struct {
	trashcan_location      string
	transaction_journal    map[string]*service_models.TrashcanTransaction
	current_transaction_id string // A datetime string
	current_transaction    *service_models.TrashcanTransaction
	categories_journal     map[string]dungeon_models.CategoryIdentity
}

func NewTrashcanDatabase(trashcan_path string) (*TrashcanDatabase, error) {

	transaction_journal, err := loadTransactionJournal(trashcan_path)
	if err != nil {
		return nil, err
	}

	categories_journal, err := loadCategoriesJournal(trashcan_path)
	if err != nil {
		return nil, err
	}

	return &TrashcanDatabase{
		trashcan_location:      trashcan_path,
		transaction_journal:    transaction_journal,
		current_transaction:    nil,
		current_transaction_id: "",
		categories_journal:     categories_journal,
	}, nil
}

func loadTransactionJournal(trashcan_path string) (map[string]*service_models.TrashcanTransaction, error) {
	transaction_journal_path := filepath.Join(trashcan_path, JOURNAL_NAME)

	if !service_helpers.FileExists(transaction_journal_path) {
		echo.EchoWarn(fmt.Sprintf("Transaction journal not found at %s", transaction_journal_path))
		return make(map[string]*service_models.TrashcanTransaction), nil
	}

	file_data, err := os.ReadFile(transaction_journal_path)
	if err != nil {
		return nil, err
	}

	data_buffer := bytes.NewBuffer(file_data)

	var transaction_journal map[string]*service_models.TrashcanTransaction = make(map[string]*service_models.TrashcanTransaction)
	err = json.NewDecoder(data_buffer).Decode(&transaction_journal)
	if err != nil {
		return nil, err
	}

	return transaction_journal, nil
}

func loadCategoriesJournal(trashcan_path string) (map[string]dungeon_models.CategoryIdentity, error) {
	categories_journal_path := filepath.Join(trashcan_path, CATEGORIES_JOURNAL_NAME)

	if !service_helpers.FileExists(categories_journal_path) {
		echo.EchoWarn(fmt.Sprintf("Categories journal not found at %s", categories_journal_path))
		return make(map[string]dungeon_models.CategoryIdentity), nil
	}

	file_data, err := os.ReadFile(categories_journal_path)
	if err != nil {
		return nil, err
	}

	data_buffer := bytes.NewBuffer(file_data)

	var categories_journal map[string]dungeon_models.CategoryIdentity = make(map[string]dungeon_models.CategoryIdentity)
	err = json.NewDecoder(data_buffer).Decode(&categories_journal)
	if err != nil {
		return nil, err
	}

	return categories_journal, nil
}

func (db *TrashcanDatabase) Commit() error {
	if db.current_transaction == nil {
		return fmt.Errorf("No transaction started")
	}

	db.transaction_journal[db.current_transaction_id] = db.current_transaction

	db.current_transaction = nil
	db.current_transaction_id = ""

	err := db.SaveTransactionsJournal()
	if err != nil {
		return err
	}

	return nil
}

func (db *TrashcanDatabase) DeleteEmptyCategory(category_identity dungeon_models.CategoryIdentity) error {
	if category_identity.Category.Fullpath == "" {
		return fmt.Errorf("Category fullpath is empty")
	}

	category_path := filepath.Join(category_identity.ClusterPath, category_identity.Category.Fullpath)

	// Check if path exists
	if !service_helpers.FileExists(category_path) {
		return fmt.Errorf("Category path %s does not exist", category_path)
	}

	// Check if path is empty
	empty, err := service_helpers.IsDirectoryEmpty(category_path)
	if err != nil {
		return err
	}

	if !empty {
		return fmt.Errorf("Category path %s is not empty", category_path)
	}

	// Delete path
	err = os.Remove(category_path)
	if err != nil {
		return err
	}

	db.categories_journal[category_identity.Category.Uuid] = category_identity
	db.SaveCategoriesJournal()

	return nil
}

func (db *TrashcanDatabase) CleanSingleTransaction(transaction_id string) error {
	transaction, exists := db.transaction_journal[transaction_id]
	if !exists {
		return fmt.Errorf("Transaction %s not found", transaction_id)
	}

	for _, media := range transaction.Content {
		media_trash_path := filepath.Join(db.trashcan_location, MEDIAS_STORAGE_DIRECTORY, media.Name)

		if !service_helpers.FileExists(media_trash_path) {
			echo.EchoWarn(fmt.Sprintf("Media '%s' not found in trash, skipping", media.Name))
			continue
		}

		err := os.Remove(media_trash_path)
		if err != nil {
			return err
		}
	}

	delete(db.transaction_journal, transaction_id)

	err := db.Save()

	return err
}

func (db *TrashcanDatabase) EmptyTrashcan() error {
	var trashcan_files_directory string = db.GetTrashcanMediaLocation()

	echo.Echo(echo.PinkBG, fmt.Sprintf("Erasing all files from '%s'", trashcan_files_directory))

	trashcan_files, err := os.ReadDir(trashcan_files_directory)
	if err != nil {
		return err
	}

	for _, file := range trashcan_files {
		file_path := filepath.Join(trashcan_files_directory, file.Name())

		err := os.Remove(file_path)
		if err != nil {
			return err
		}
	}

	db.transaction_journal = make(map[string]*service_models.TrashcanTransaction)

	err = db.Save()

	return err
}

// Deletes a transaction from the journal without deleting the files
func (db *TrashcanDatabase) EraseTransaction(transaction_id string) (*service_models.TrashcanTransaction, error) {
	transaction, exists := db.transaction_journal[transaction_id]
	if !exists {
		return nil, fmt.Errorf("Transaction %s not found", transaction_id)
	}

	delete(db.transaction_journal, transaction_id)

	err := db.Save()

	return transaction, err
}

// Deletes a media from a transaction without deleting the file
func (db *TrashcanDatabase) EraseMediaFromTransaction(transaction_id string, media_uuid string) (*dungeon_models.Media, error) {
	transaction, exists := db.transaction_journal[transaction_id]
	if !exists {
		return nil, fmt.Errorf("Transaction %s not found", transaction_id)
	}

	var media_to_erase_index int = -1
	var media_to_erase *dungeon_models.Media = nil

	for h, media := range transaction.Content {
		if media.Uuid == media_uuid {
			media_to_erase = &media
			media_to_erase_index = h
			break
		}
	}

	if media_to_erase_index == -1 {
		return nil, fmt.Errorf("Media %s not found in transaction %s", media_uuid, transaction_id)
	}

	transaction.Content = append(transaction.Content[:media_to_erase_index], transaction.Content[media_to_erase_index+1:]...)

	if len(transaction.Content) == 0 {
		delete(db.transaction_journal, transaction_id)
	}

	err := db.Save()

	return media_to_erase, err
}

func (db *TrashcanDatabase) GetTrashcanLocation() string {
	return db.trashcan_location
}

func (db *TrashcanDatabase) GetTrashcanMediaLocation() string {
	return filepath.Join(db.trashcan_location, MEDIAS_STORAGE_DIRECTORY)
}

func (db *TrashcanDatabase) GetTransactions() []*service_models.TrashcanTransaction {
	var transactions []*service_models.TrashcanTransaction = make([]*service_models.TrashcanTransaction, 0)

	for _, transaction := range db.transaction_journal {
		transactions = append(transactions, transaction)
	}

	return transactions
}

// return all transactions sorted by date
func (db *TrashcanDatabase) GetSortedTransactions() []*service_models.TrashcanTransaction {
	var transactions []*service_models.TrashcanTransaction = db.GetTransactions()

	sort.Slice(transactions, func(h, k int) bool {
		h_time, _ := time.Parse(datetime_format, transactions[h].TransactionID)
		k_time, _ := time.Parse(datetime_format, transactions[k].TransactionID)

		return h_time.After(k_time)
	})

	return transactions
}

func (db *TrashcanDatabase) GetSortedTransactionEntries() []*service_models.TrashcanTransactionEntry {
	var sorted_transactions []*service_models.TrashcanTransaction = db.GetSortedTransactions()
	var transaction_entries []*service_models.TrashcanTransactionEntry = make([]*service_models.TrashcanTransactionEntry, 0)

	for _, transaction := range sorted_transactions {
		transaction_entries = append(transaction_entries, transaction.ToEntry())
	}

	return transaction_entries
}

func (db *TrashcanDatabase) GetTransaction(transaction_id string) (*service_models.TrashcanTransaction, error) {
	transaction, exists := db.transaction_journal[transaction_id]
	if !exists {
		return nil, fmt.Errorf("Transaction %s not found", transaction_id)
	}

	return transaction, nil
}

func (db *TrashcanDatabase) GetTrashcanSize() int {
	var trashcan_size int = 0

	for _, transaction := range db.transaction_journal {
		trashcan_size += len(transaction.Content)
	}

	return trashcan_size
}

func (db *TrashcanDatabase) HasTransaction(transaction_id string) bool {
	_, exists := db.transaction_journal[transaction_id]
	return exists
}

func (db *TrashcanDatabase) MoveToTrash(media_identity *dungeon_models.MediaIdentity) error {
	if db.current_transaction_id == "" {
		return fmt.Errorf("No transaction started")
	}

	media_trash_path := filepath.Join(db.trashcan_location, MEDIAS_STORAGE_DIRECTORY, media_identity.Media.Name)
	current_media_abs_path := filepath.Join(media_identity.ClusterPath, media_identity.CategoryPath, media_identity.Media.Name)

	err := service_helpers.MoveFile(current_media_abs_path, media_trash_path)
	if err != nil {
		return err
	}

	db.current_transaction.AddMedia(*media_identity.Media)

	return nil
}

func (db *TrashcanDatabase) RestoreMediaFromTrash(rejected_media dungeon_models.Media, original_path string) error {
	media_trash_path := filepath.Join(db.trashcan_location, "medias", rejected_media.Name)
	original_media_path := filepath.Join(original_path, rejected_media.Name)

	err := service_helpers.MoveFile(media_trash_path, original_media_path)
	if err != nil {
		return err
	}

	return nil
}

func (db *TrashcanDatabase) Reload() error {
	transaction_journal, err := loadTransactionJournal(db.trashcan_location)
	if err != nil {
		return err
	}

	categories_journal, err := loadCategoriesJournal(db.trashcan_location)
	if err != nil {
		return err
	}

	db.transaction_journal = transaction_journal
	db.categories_journal = categories_journal

	return nil
}

// DELETE THIS FUNCTION
func (db *TrashcanDatabase) RestoreTo(new_parent_path, transaction_id, media_uuid string) (*dungeon_models.Media, error) {
	transaction, exists := db.transaction_journal[transaction_id]
	if !exists {
		return nil, fmt.Errorf("Transaction %s not found", transaction_id)
	}

	var media_to_restore_index int = -1
	var media_to_restore *dungeon_models.Media = nil

	for h, media := range transaction.Content {
		if media.Uuid == media_uuid {
			media_to_restore = &media
			media_to_restore_index = h
			break
		}
	}

	if media_to_restore == nil {
		return nil, fmt.Errorf("Media %s not found in transaction %s", media_uuid, transaction_id)
	}

	media_trash_path := filepath.Join(db.trashcan_location, MEDIAS_STORAGE_DIRECTORY, media_to_restore.Name)
	new_path := filepath.Join(new_parent_path, media_to_restore.Name)

	is_same_filesystem, err := dungeon_helpers.IsSameFilesystem(new_parent_path, db.trashcan_location)
	if err != nil {
		return nil, err
	}

	if !is_same_filesystem {
		err := dungeon_helpers.MoveFile(media_trash_path, new_path) // can handle different filesystems but is slower
		if err != nil {
			return nil, err
		}
	} else {
		err := os.Rename(media_trash_path, new_path) // faster but breaks if paths are on different filesystems
		if err != nil {
			return nil, err
		}
	}

	transaction.Content = append(transaction.Content[:media_to_restore_index], transaction.Content[media_to_restore_index+1:]...)

	if len(transaction.Content) == 0 {
		delete(db.transaction_journal, transaction_id)
	}

	err = db.Save()

	return media_to_restore, err
}

func (db *TrashcanDatabase) Rollback() error {
	var original_path string = filepath.Join(db.current_transaction.OriginIdentity.ClusterPath, db.current_transaction.OriginIdentity.CategoryPath)

	if db.current_transaction_id == "" {
		return fmt.Errorf("No transaction started")
	}

	for _, media := range db.current_transaction.Content {
		err := db.RestoreMediaFromTrash(media, original_path)
		if err != nil {
			return err
		}
	}

	db.current_transaction = nil
	db.current_transaction_id = ""

	return nil
}

func (db *TrashcanDatabase) StartTransaction(category_identity *dungeon_models.CategoryWeakIdentity) error {
	localtime, err := time.LoadLocation(app_config.LOCALTIME)
	if err != nil {
		return err
	}

	current_time := time.Now().In(localtime)

	db.current_transaction_id = current_time.Format(datetime_format)

	var new_transaction *service_models.TrashcanTransaction = service_models.NewTrashcanTransaction(db.current_transaction_id, *category_identity)

	db.current_transaction = new_transaction

	return nil
}

func (db *TrashcanDatabase) SaveTransactionsJournal() error {
	transaction_journal_path := filepath.Join(db.trashcan_location, JOURNAL_NAME)

	journal_json, err := json.Marshal(db.transaction_journal)
	if err != nil {
		return err
	}

	err = os.WriteFile(transaction_journal_path, journal_json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *TrashcanDatabase) SaveCategoriesJournal() error {
	categories_journal_path := filepath.Join(db.trashcan_location, CATEGORIES_JOURNAL_NAME)

	journal_json, err := json.Marshal(db.categories_journal)
	if err != nil {
		return err
	}

	err = os.WriteFile(categories_journal_path, journal_json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *TrashcanDatabase) Save() error {
	err := db.SaveTransactionsJournal()
	if err != nil {
		fmt.Printf("Error saving transactions journal: %s", err.Error())
		return err
	}

	err = db.SaveCategoriesJournal()
	if err != nil {
		fmt.Printf("Error saving categories journal: %s", err.Error())
		return err
	}

	return nil
}
