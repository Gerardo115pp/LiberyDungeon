package dungeon_sqlite_opener

import (
	"database/sql"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	_ "github.com/mattn/go-sqlite3"
)

type DungeonSqliteOpener struct {
	SchemaFileName     string
	DatabaseFileName   string
	SchemasDirectory   string
	DatabasesDirectory string
	OperationDataPath  string
}

func NewDungeonSqliteOpener(database_filename, schema_filename, operation_data_path string) *DungeonSqliteOpener {
	var ds_opener *DungeonSqliteOpener = new(DungeonSqliteOpener)

	ds_opener.DatabaseFileName = database_filename
	ds_opener.SchemaFileName = schema_filename
	ds_opener.SchemasDirectory = "schemas"
	ds_opener.DatabasesDirectory = "databases"
	ds_opener.OperationDataPath = operation_data_path

	return ds_opener
}

func (ds_opener DungeonSqliteOpener) getDungeonTagsSchemaPath() string {
	var schema_file string = filepath.Join(ds_opener.OperationDataPath, ds_opener.SchemasDirectory, ds_opener.SchemaFileName)

	return schema_file
}

func (ds_opener *DungeonSqliteOpener) dbSchemaExist() bool {
	var schema_file string = ds_opener.getDungeonTagsSchemaPath()

	return dungeon_helpers.FileExists(schema_file)
}

func (ds_opener DungeonSqliteOpener) getSchema() (string, error) {
	var schema_file_path string = ds_opener.getDungeonTagsSchemaPath()

	if !ds_opener.dbSchemaExist() {
		return "", fmt.Errorf("Schema directory or files not found at %s", schema_file_path)
	}

	schema, err := os.ReadFile(schema_file_path)
	if err != nil {
		return "", fmt.Errorf("Error reading schema file: %s", err.Error())
	}

	return string(schema), nil
}

func (ds_opener DungeonSqliteOpener) OpenDB(with_foreign_keys bool) (*sql.DB, error) {
	var err error

	var db_files_directory string = filepath.Join(ds_opener.OperationDataPath, ds_opener.DatabasesDirectory)
	var db_file_path string = filepath.Join(db_files_directory, ds_opener.DatabaseFileName)

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
		err = ds_opener.writeSchema(db)
		if err != nil {
			return nil, fmt.Errorf("Error writing schema: %s", err.Error())
		}
	}

	if with_foreign_keys {
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return nil, fmt.Errorf("Error enabling foreign keys: %s", err)
		}
	}

	return db, nil
}

func (ds_opener *DungeonSqliteOpener) SetDatabaseFileName(database_filename string) {
	ds_opener.DatabaseFileName = database_filename
}

func (ds_opener *DungeonSqliteOpener) SetSchemaFileName(schema_filename string) {
	ds_opener.SchemaFileName = schema_filename
}

func (ds_opener *DungeonSqliteOpener) SetDatabasesDirectory(databases_directory string) {
	ds_opener.DatabasesDirectory = databases_directory
}

func (ds_opener *DungeonSqliteOpener) SetSchemasDirectory(schemas_directory string) {
	ds_opener.SchemasDirectory = schemas_directory
}

func (ds_opener DungeonSqliteOpener) writeSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("trying to write schema to nil database")
	}

	echo.Echo(echo.PurpleFG, "Writing schema to database")

	schema, err := ds_opener.getSchema()
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
