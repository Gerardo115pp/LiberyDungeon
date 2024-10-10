package watch_points_database

import (
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	app_config "libery-metadata-service/Config"
	"os"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func getWatchPointsFilename() string {
	var filename string = fmt.Sprintf("%s.watch_points", app_config.SERVICE_ID)
	return filepath.Join(app_config.OPERATION_DATA_PATH, filename)
}

func loadWatchPoints() ([]byte, error) {
	var watch_points_filename string = getWatchPointsFilename()

	if !dungeon_helpers.FileExists(watch_points_filename) {
		echo.EchoWarn(fmt.Sprintf("Watch points file<%s> not found", watch_points_filename))
		return nil, nil
	}

	watch_points, err := os.ReadFile(watch_points_filename)
	if err != nil {
		return nil, err
	}

	return watch_points, nil
}
