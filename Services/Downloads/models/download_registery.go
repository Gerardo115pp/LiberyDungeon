package models

import (
	"encoding/json"
	"strconv"
	"time"
)

type RegisteredDownload struct {
	DownloadUuid      string                   `json:"download_uuid"`
	DownloadTimestamp time.Time                `json:"download_timestamp"`
	CategoryUuid      string                   `json:"category_uuid"`
	DownloadFiles     []RegisteredDownloadFile `json:"download_files"`
}

type RegisteredDownloadFile struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

// ------ RegisteredDownload Methods ------

func (rd *RegisteredDownload) UnixTimestamp() int64 {
	return rd.DownloadTimestamp.Unix()
}

func (rd *RegisteredDownload) MarshalJSON() ([]byte, error) {

	json_files, err := json.Marshal(rd.DownloadFiles)
	if err != nil {
		return nil, err
	}

	download_timestamp := rd.UnixTimestamp()

	json_timestamp := strconv.FormatInt(download_timestamp, 10)

	return []byte(`{"download_uuid":"` + rd.DownloadUuid + `","download_timestamp":` + json_timestamp + `,"category_uuid":"` + rd.CategoryUuid + `","download_files":` + string(json_files) + `}`), nil
}

func (rd *RegisteredDownload) DownloadedFilesMap() map[string]bool {
	var downloaded_files_map map[string]bool = make(map[string]bool)

	for _, download_file := range rd.DownloadFiles {
		downloaded_files_map[download_file.Url] = true
	}

	return downloaded_files_map
}
