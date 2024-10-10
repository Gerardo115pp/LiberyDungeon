package workflows

import (
	"fmt"
	"libery_downloads_service/workflows/jobs"
	"net/http"
)

func RegisterProgressListener(download_uuid string, response http.ResponseWriter, request *http.Request) error {
	upgrader := jobs.DownloadWorker.GetUpgrader()
	if upgrader == nil {
		return fmt.Errorf("While registering progress listener: No upgrader available")
	}

	connection, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		return fmt.Errorf("While upgrading connection: %s", err.Error())
	}

	err = jobs.DownloadWorker.RegisterDownloadListener(download_uuid, connection)
	if err != nil {
		connection.Close()
		return fmt.Errorf("While registering download listener: %s.", err.Error())
	}

	return nil
}
