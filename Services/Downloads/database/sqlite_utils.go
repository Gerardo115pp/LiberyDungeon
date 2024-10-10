package database

import (
	"database/sql"
	"fmt"
	"libery-dungeon-libs/helpers"
	app_config "libery_downloads_service/Config"
	"os"
	"path"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func getDownloadsSchemaPath(settings_directory string) string {
	var schema_file string = path.Join(settings_directory, "schemas", "downloads.sql")

	return schema_file
}

func dbSchemaExist(settings_directory string) bool {
	var schema_file string = getDownloadsSchemaPath(settings_directory)

	return helpers.FileExists(schema_file)
}

func getSchema() (string, error) {
	var schema_file string = getDownloadsSchemaPath(app_config.OPERATION_DATA_PATH)
	if !dbSchemaExist(app_config.OPERATION_DATA_PATH) {
		return "", fmt.Errorf("Schema directory or files not found at %s", schema_file)
	}

	schema, err := os.ReadFile(schema_file)
	if err != nil {
		return "", fmt.Errorf("Error reading schema file: %s", err.Error())
	}

	return string(schema), nil
}

func openDownloadDB() (*sql.DB, error) {
	var err error

	var db_files_directory string = path.Join(app_config.OPERATION_DATA_PATH, "databases")
	var db_file_path string = path.Join(db_files_directory, "downloads.db")

	var db *sql.DB

	var database_exists bool = helpers.FileExists(db_file_path)

	if !database_exists {

		if !helpers.FileExists(db_files_directory) {
			err = os.Mkdir(db_files_directory, 0777)
			if err != nil {
				return nil, fmt.Errorf("Error creating database directory: %s", err.Error())
			}
		}

		db_file, err := os.OpenFile(db_file_path, os.O_CREATE, 0777)
		if err != nil {
			return nil, fmt.Errorf("Error opening database file: %s", err.Error())
		}

		db_file.Close()
	}

	db, err = sql.Open("sqlite3", db_file_path)
	if err != nil {
		echo.EchoFatal(fmt.Errorf("Error opening database: %s", err.Error())) // because of the defer writeSchema(db) line, even if we just return the error, the writeSchema would panic because db is nil
	}

	if !database_exists {
		err = writeSchema(db)
		if err != nil {
			return nil, fmt.Errorf("Error writing schema: %s", err.Error())
		}
	}

	return db, nil
}

func writeSchema(db *sql.DB) error {
	if db == nil {
		echo.EchoFatal(fmt.Errorf("trying to write schema to nil database"))
	}

	echo.Echo(echo.PurpleFG, "Writing schema to database")

	schema, err := getSchema()
	if err != nil {
		return err
	}

	commands := strings.Split(schema, ";")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %s", err.Error())
	}

	for _, command := range commands {
		cmd := strings.TrimSpace(command)
		if cmd == "" {
			continue
		}

		echo.Echo(echo.BlueFG, fmt.Sprintf("Executing command: %s", cmd))

		if _, err := tx.Exec(cmd); err != nil {
			tx.Rollback()
			return fmt.Errorf("Error executing command: %s", err.Error())
		}
	}

	return tx.Commit()
}
