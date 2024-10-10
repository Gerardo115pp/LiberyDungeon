package models

type DownloadProgressMessage struct {
	DownloadUuid    string `json:"download_uuid"`
	TotalFiles      int    `json:"total_files"`
	DownloadedFiles int    `json:"downloaded_files"`
	Completed       bool   `json:"completed"`
}

func CreateDownloadProgressMessage(download_uuid string, total_files int, downloaded_files int, completed bool) *DownloadProgressMessage {
	return &DownloadProgressMessage{
		DownloadUuid:    download_uuid,
		TotalFiles:      total_files,
		DownloadedFiles: downloaded_files,
		Completed:       completed,
	}
}
