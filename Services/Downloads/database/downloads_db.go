package database

import (
	"database/sql"
	"fmt"
	"libery_downloads_service/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DownloadDB struct {
	db *sql.DB
}

func NewDownloadDB() (*DownloadDB, error) {
	var download_db *DownloadDB = new(DownloadDB)

	db, err := openDownloadDB()
	if err != nil {
		return nil, err
	}

	download_db.db = db

	return download_db, nil
}

func (download_db *DownloadDB) DownloadExists(download_uuid string) (bool, error) {
	stmt, err := download_db.db.Prepare("SELECT count(*) FROM downloads WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("Error preparing statement for download exists: %s", err)
	}

	defer stmt.Close()

	var count int
	err = stmt.QueryRow(download_uuid).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error getting download exists: %s", err)
	}

	return count > 0, nil
}

func (download_db *DownloadDB) GetDownload(download_uuid string) (*models.RegisteredDownload, error) {
	var download *models.RegisteredDownload = new(models.RegisteredDownload)

	var download_timestamp int64

	stmt, err := download_db.db.Prepare("SELECT download_timestamp, category_id FROM downloads WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("Error preparing statement for get download: %s", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(download_uuid).Scan(&download_timestamp, &download.CategoryUuid)
	if err != nil {
		return nil, fmt.Errorf("Error getting download: %s", err)
	}

	download_files, err := download_db.GetDownloadFiles(download_uuid)
	if err != nil {
		return nil, err
	}

	download.DownloadUuid = download_uuid
	download.DownloadTimestamp = time.Unix(download_timestamp, 0)
	download.DownloadFiles = download_files

	return download, nil
}

func (download_db *DownloadDB) GetDownloadFiles(download_uuid string) ([]models.RegisteredDownloadFile, error) {
	var download_files []models.RegisteredDownloadFile = make([]models.RegisteredDownloadFile, 0)

	stmt, err := download_db.db.Prepare("SELECT id, url FROM download_files WHERE download = ?")
	if err != nil {
		return nil, fmt.Errorf("Error preparing statement for get download files: %s", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(download_uuid)
	if err != nil {
		return nil, fmt.Errorf("Error getting download files: %s", err)
	}

	defer rows.Close()

	var download_file models.RegisteredDownloadFile
	for rows.Next() {
		download_file = models.RegisteredDownloadFile{}

		err = rows.Scan(&download_file.Id, &download_file.Url)
		if err != nil {
			return nil, fmt.Errorf("Error getting download file: %s", err)
		}

		download_files = append(download_files, download_file)
	}

	return download_files, nil
}

func (download_db *DownloadDB) InsertDownload(download *models.DownloadRequest) error {
	var err error

	download_timestamp := time.Now().Unix()

	stmt, err := download_db.db.Prepare("INSERT INTO downloads(id, download_timestamp, category_id) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement for insert download: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(download.DownloadUuid, download_timestamp, download.CategoryUuid)
	if err != nil {
		return fmt.Errorf("Error inserting download: %s", err)
	}

	return download_db.InsertDownloadFiles(download)
}

func (download_db *DownloadDB) InsertDownloadFiles(download *models.DownloadRequest) error {
	stmt, err := download_db.db.Prepare("INSERT INTO download_files(download, url) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement for insert download files: %s", err)
	}

	defer stmt.Close()

	for h := 0; h < download.Len(); h++ {
		var download_file *models.DownloadFile = download.Get(h)

		_, err = stmt.Exec(download.DownloadUuid, download_file.Url)
		if err != nil {
			return fmt.Errorf("Error inserting download file '%s': %s", download_file.Url, err)
		}
	}

	return nil
}

func (download_db *DownloadDB) UpdateDownloadFiles(download_uuid string, new_file_urls []models.DownloadFile) error {
	tx, err := download_db.db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %s", err)
	}

	stmt, err := tx.Prepare("INSERT INTO download_files(download, url) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement for update download files: %s", err)
	}

	defer stmt.Close()

	for _, new_file_url := range new_file_urls {
		_, err = stmt.Exec(download_uuid, new_file_url.Url)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting download file '%s': %s", new_file_url.Url, err)
		}
	}

	tx.Commit()

	return download_db.updateDownloadTimestamp(download_uuid)
}

func (download_db *DownloadDB) updateDownloadTimestamp(download_uuid string) error {
	download_timestamp := time.Now().Unix()

	stmt, err := download_db.db.Prepare("UPDATE downloads SET download_timestamp = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("Error preparing statement for update download timestamp: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(download_timestamp, download_uuid)
	if err != nil {
		return fmt.Errorf("Error updating download timestamp: %s", err)
	}

	return nil
}

func (download_db *DownloadDB) Close() error {
	return download_db.db.Close()
}
