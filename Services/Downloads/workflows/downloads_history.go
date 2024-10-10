package workflows

import (
	"fmt"
	"libery_downloads_service/repository"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func ClearRepeatedDownloadFiles(download_uuid string, images_urls []string) []string {
	download_register, err := repository.Downloads.GetDownload(download_uuid)
	if err != nil {
		echo.EchoDebug(fmt.Sprintf("Error getting download register: %s", err.Error()))
		return images_urls
	}
	unique_image_urls := make([]string, 0)
	downloaded_files := download_register.DownloadedFilesMap()

	for _, image_url := range images_urls {
		if _, exists := downloaded_files[image_url]; !exists {
			unique_image_urls = append(unique_image_urls, image_url)
		}

	}

	return unique_image_urls
}
