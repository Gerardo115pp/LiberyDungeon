package models

import (
	dungeon_models "libery-dungeon-libs/models"

	"github.com/google/uuid"
)

type DownloadRequest struct {
	DownloadUuid    string
	download_batch  []DownloadFile
	CategoryCluster *dungeon_models.CategoryCluster
	CategoryUuid    string
}

type DownloadFile struct {
	Url          string
	IsDownloaded bool
	Trys         int
}

// ------ DownloadRequest Methods ------

func CreateNewDownloadRequest(download_images []string, category_uuid string, category_cluster *dungeon_models.CategoryCluster) *DownloadRequest {
	var new_download_request *DownloadRequest = new(DownloadRequest)
	new_download_request.DownloadUuid = uuid.New().String()
	new_download_request.CategoryUuid = category_uuid
	new_download_request.CategoryCluster = category_cluster

	new_download_request.download_batch = make([]DownloadFile, 0)
	var new_download_file *DownloadFile = new(DownloadFile)

	for _, image := range download_images {
		new_download_file.Url = image
		new_download_file.IsDownloaded = false

		new_download_request.download_batch = append(new_download_request.download_batch, *new_download_file)
	}

	return new_download_request
}

func (dr *DownloadRequest) IsDownloaded() bool {
	for _, file := range dr.download_batch {
		if !file.IsDownloaded {
			return false
		}
	}
	return true
}

func (dr *DownloadRequest) Get(h int) *DownloadFile {
	return &dr.download_batch[h]
}

func (dr *DownloadRequest) DownloadFiles() []DownloadFile {
	return dr.download_batch
}

func (dr *DownloadRequest) Len() int {
	return len(dr.download_batch)
}
