package dungeon_tags

import (
	"database/sql"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	app_config "libery-metadata-service/Config"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	_ "github.com/mattn/go-sqlite3"
)

func getDungeonTagsSchemaPath(settings_directory string) string {
	var schema_file string = filepath.Join(settings_directory, "schemas", "dungeon_tags.sql")

	return schema_file
}

func dbSchemaExist(settings_directory string) bool {
	var schema_file string = getDungeonTagsSchemaPath(settings_directory)

	return dungeon_helpers.FileExists(schema_file)
}

func getSchema() (string, error) {
	var schema_file_path string = getDungeonTagsSchemaPath(app_config.OPERATION_DATA_PATH)

	if !dbSchemaExist(app_config.OPERATION_DATA_PATH) {
		return "", fmt.Errorf("Schema directory or files not found at %s", schema_file_path)
	}

	schema, err := os.ReadFile(schema_file_path)
	if err != nil {
		return "", fmt.Errorf("Error reading schema file: %s", err.Error())
	}

	return string(schema), nil
}

func openDB() (*sql.DB, error) {
	var err error

	var db_files_directory string = filepath.Join(app_config.OPERATION_DATA_PATH, "databases")
	var db_file_path string = filepath.Join(db_files_directory, "dungeon_tags.db")

	var db *sql.DB

	var database_exists bool = dungeon_helpers.FileExists(db_file_path)

	if !database_exists {

		if !dungeon_helpers.FileExists(db_files_directory) {
			err = os.Mkdir(db_files_directory, 0777)
			if err != nil {
				return nil, fmt.Errorf("Error creating database directory: %s", err.Error())
			}
		}

		db_file, err := os.OpenFile(db_file_path, os.O_CREATE, 0777)
		if err != nil {
			return nil, fmt.Errorf("Error creating database file: %s", err.Error())
		}

		db_file.Close()
	}

	echo.Echo(echo.PinkFG, fmt.Sprintf("Opening database at %s", db_file_path))

	db, err = sql.Open("sqlite3", db_file_path)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %s", err.Error())
	}

	if !database_exists {
		err = writeSchema(db)
		if err != nil {
			return nil, fmt.Errorf("Error writing schema: %s", err.Error())
		}
	}

	return db, nil
}

// NOTE: This is an amazing candidate for libery-dungeon-libs. as it is a common pattern
func writeSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("trying to write schema to nil database")
	}

	echo.Echo(echo.PurpleFG, "Writing schema to database")

	schema, err := getSchema()
	if err != nil {
		return fmt.Errorf("Error getting schema: %s", err.Error())
	}

	statements := strings.Split(schema, ";")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %s", err.Error())
	}

	for _, statement := range statements {
		cmd := strings.TrimSpace(statement)
		if cmd == "" {
			continue
		}

		echo.Echo(echo.PurpleFG, fmt.Sprintf("Executing statement: %s", cmd))

		_, err = tx.Exec(cmd)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error executing statement: %s", err.Error())
		}
	}

	return tx.Commit()
}
