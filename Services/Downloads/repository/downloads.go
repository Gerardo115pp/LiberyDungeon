package repository

import "libery_downloads_service/models"

type DownloadsRepository interface {
	DownloadExists(download_uuid string) (bool, error)
	GetDownload(download_uuid string) (*models.RegisteredDownload, error)
	GetDownloadFiles(download_uuid string) ([]models.RegisteredDownloadFile, error)
	InsertDownload(download *models.DownloadRequest) error
	InsertDownloadFiles(download *models.DownloadRequest) error
	UpdateDownloadFiles(download_uuid string, new_file_urls []models.DownloadFile) error
	Close() error
}

var Downloads DownloadsRepository

func SetDownloadsRepository(downloads_impl DownloadsRepository) {
	Downloads = downloads_impl
}
