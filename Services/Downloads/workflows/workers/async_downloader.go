package workers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_downloads_service/Config"
	"libery_downloads_service/models"
	"libery_downloads_service/repository"
	"net/http"
	"path"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/gorilla/websocket"
)

type AsyncDownloader struct {
	download_queue     *models.DownloadQueue
	download_enqueued  chan bool
	listener_upgrader  *websocket.Upgrader
	progress_listeners map[string]*websocket.Conn
}

func NewAsyncDownloader() *AsyncDownloader {
	var download_queue *models.DownloadQueue = new(models.DownloadQueue)
	var download_enqueued chan bool = make(chan bool)

	var async_downloader *AsyncDownloader = new(AsyncDownloader)

	async_downloader.download_queue = download_queue
	async_downloader.download_enqueued = download_enqueued

	async_downloader.listener_upgrader = &websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	async_downloader.progress_listeners = make(map[string]*websocket.Conn)

	go async_downloader.monitorDownloadQueue()

	return async_downloader
}

func (ad *AsyncDownloader) GetUpgrader() *websocket.Upgrader {
	return ad.listener_upgrader
}

func (ad *AsyncDownloader) GetCurrentDownloadUUID() string {
	if ad.download_queue.IsEmpty() {
		return ""
	}

	return ad.download_queue.Peek().DownloadUuid
}

func (ad *AsyncDownloader) DownloadImagesBatch(category_uuid string, image_urls []string, custom_download_id string, category_cluster *dungeon_models.CategoryCluster) (string, error) {
	new_download_request := models.CreateNewDownloadRequest(image_urls, category_uuid, category_cluster)

	if custom_download_id != "" {
		new_download_request.DownloadUuid = custom_download_id
	}

	if len(image_urls) > 0 {
		ad.download_queue.Enqueue(new_download_request)

		ad.download_enqueued <- true
	}

	return new_download_request.DownloadUuid, nil
}

func (ad *AsyncDownloader) monitorDownloadQueue() {
	var new_download_received bool

	for {

		if ad.download_queue.IsEmpty() {
			select {
			case new_download_received = <-ad.download_enqueued:
				if !new_download_received {
					return
				}
			}
		}

		echo.Echo(echo.GreenFG, "Downloading batch")

		download_request := ad.download_queue.Peek()
		if download_request == nil {
			echo.EchoWarn("Download received but queue is empty")
			continue
		}

		// time.Sleep(500 * time.Millisecond) // Simulate download time

		ad.downloadImagesRequest(download_request)

		_, err := ad.download_queue.Dequeue()
		if err != nil {
			echo.EchoErr(err)
		}

		err = ad.registerDownload(download_request)

		echo.Echo(echo.GreenFG, fmt.Sprintf("Downloaded batch: %s", download_request.DownloadUuid))
	}
}

func (ad *AsyncDownloader) downloadImagesRequest(download_request *models.DownloadRequest) {
	var current_file *models.DownloadFile

	communication.Medias.GetUploadStreamTicket(download_request.DownloadUuid, download_request.CategoryUuid, download_request.Len())

	for h := 0; h < download_request.Len(); h++ {
		current_file = download_request.Get(h)

		if current_file.IsDownloaded {
			continue
		}

		downloaded := ad.downloadImageFile(current_file.Url, download_request.CategoryUuid)

		for !downloaded && current_file.Trys < 3 {
			if !downloaded {
				time.Sleep(500 * time.Millisecond) // give the server some time to recover
			}

			downloaded = ad.downloadImageFile(current_file.Url, download_request.CategoryUuid)

			current_file.Trys++
		}

		// Update download progress for listeners
		ad.updateProgress(download_request.DownloadUuid, h+1, download_request.Len())

		current_file.IsDownloaded = true
	}
}

func (ad *AsyncDownloader) downloadImageFile(image_url string, target_category string) bool {
	client := &http.Client{}

	image_http_request, err := http.NewRequest("GET", image_url, nil)
	if err != nil {
		echo.EchoErr(fmt.Errorf("Error creating image request: %s", err))
		return false
	}

	image_http_request.Header.Set("User-Agent", app_config.DOWNLOAD_USER_AGENT)

	image_http_response, err := client.Do(image_http_request)
	if err != nil {
		echo.EchoErr(fmt.Errorf("Error requesting image: %s", err))
		return false
	}
	defer image_http_response.Body.Close()

	if image_http_response.StatusCode != 200 {
		echo.EchoErr(fmt.Errorf("Error requesting image: %s", image_http_response.Status))
		return false
	}

	filename := path.Base(image_url)
	echo.Echo(echo.CyanFG, fmt.Sprintf("Downloading image: %s", filename))

	err = communication.Medias.UploadMediaFile(image_http_response.Body, filename)
	if err != nil {
		echo.EchoErr(fmt.Errorf("Error uploading image: %s", err))
		return false
	}

	return true
}

func (ad *AsyncDownloader) registerDownload(download_request *models.DownloadRequest) (err error) {
	download_exists, err := repository.Downloads.DownloadExists(download_request.DownloadUuid)
	if err != nil {
		return fmt.Errorf("Error checking if download exists: %s", err)
	}

	if !download_exists {
		err = repository.Downloads.InsertDownload(download_request)
		if err != nil {
			echo.EchoErr(err)
		}
	} else {
		repository.Downloads.UpdateDownloadFiles(download_request.DownloadUuid, download_request.DownloadFiles())
	}

	return
}

func (ad *AsyncDownloader) RegisterDownloadListener(download_uuid string, ws *websocket.Conn) error {
	if _, exists := ad.progress_listeners[download_uuid]; exists {
		return fmt.Errorf("Download progress connection already exists")
	}

	ad.progress_listeners[download_uuid] = ws

	return nil
}

func (ad *AsyncDownloader) updateProgress(download_uuid string, downloaded_files int, total_files int) {

	peer, has_listener := ad.progress_listeners[download_uuid]
	if !has_listener {
		echo.EchoWarn(fmt.Sprintf("No listener for download: %s", download_uuid))
		return
	}

	is_completed := downloaded_files == total_files

	progress := models.CreateDownloadProgressMessage(download_uuid, total_files, downloaded_files, is_completed)
	progress_json, err := json.Marshal(progress)
	if err != nil {
		echo.EchoErr(err)
		return
	}

	err = peer.WriteMessage(websocket.TextMessage, progress_json)
	if err != nil {
		// Peer disconnected
		peer.Close()
		delete(ad.progress_listeners, download_uuid)
		return
	}

	if is_completed {
		peer.Close()
		delete(ad.progress_listeners, download_uuid)
	}
}

func (ad *AsyncDownloader) Stop() {
	ad.download_enqueued <- false
}
