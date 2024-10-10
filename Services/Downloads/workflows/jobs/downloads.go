package jobs

import (
	dungeon_models "libery-dungeon-libs/models"

	"github.com/gorilla/websocket"
)

type DownloadManager interface {
	DownloadImagesBatch(category_uuid string, image_urls []string, custom_download_id string, category_cluster *dungeon_models.CategoryCluster) (string, error)
	GetCurrentDownloadUUID() string
	RegisterDownloadListener(download_uuid string, ws *websocket.Conn) error
	GetUpgrader() *websocket.Upgrader
}

var DownloadWorker DownloadManager

func SetDownloadWorkerImplementation(impl DownloadManager) {
	DownloadWorker = impl
}
