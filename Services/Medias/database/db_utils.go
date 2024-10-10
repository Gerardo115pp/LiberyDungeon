package database

import (
	"fmt"
	app_config "libery_categories_service/Config"
)

func createDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", app_config.MYSQL_USER, app_config.MYSQL_PASSWORD, app_config.MYSQL_HOST, app_config.MYSQL_PORT, app_config.MYSQL_DB)
}
