package sqlite_users

import (
	"database/sql"
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	app_config "libery_users_service/Config"
	service_models "libery_users_service/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	_ "github.com/mattn/go-sqlite3"
)

func getUsersSchemaPath(settings_directory string) string {
	var schema_file string = filepath.Join(settings_directory, "schemas", "users_sqlite.sql")

	return schema_file
}

func dbSchemaExist(settings_directory string) bool {
	var schema_file string = getUsersSchemaPath(settings_directory)

	return dungeon_helpers.FileExists(schema_file)
}

func getSchema() (string, error) {
	var schema_file_path string = getUsersSchemaPath(app_config.OPERATION_DATA_PATH)

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
	var db_file_path string = filepath.Join(db_files_directory, "users.db")

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

		err = createDefaultRoles(db)
		if err != nil {
			return nil, fmt.Errorf("Error creating default roles: %s", err.Error())
		}
	}

	return db, nil
}

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

func createDefaultRoles(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %s", err.Error())
	}

	err = insertDefaultRoles(tx)
	if err != nil {
		return fmt.Errorf("Error inserting default roles: %s", err.Error())
	}

	err = insertDefaultRoleGrants(tx)
	if err != nil {
		return fmt.Errorf("Error inserting default role grants: %s", err.Error())
	}

	return tx.Commit()
}

func insertDefaultRoles(tx *sql.Tx) error {
	super_admin_stmt, err := tx.Prepare("INSERT INTO `roles`(`role_label`, `role_hierarchy`) VALUES (?, 0)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error preparing statement: %s", err.Error())
	}
	defer super_admin_stmt.Close()

	_, err = super_admin_stmt.Exec(service_models.SUPER_ADMIN_ROLE)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error inserting super admin role: %s", err.Error())
	}

	admin_stmt, err := tx.Prepare("INSERT INTO `roles`(`role_label`, `role_hierarchy`) VALUES (?, 1)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error preparing statement: %s", err.Error())
	}
	defer admin_stmt.Close()

	_, err = admin_stmt.Exec(service_models.ADMIN_ROLE)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error inserting admin role: %s", err.Error())
	}

	visitor_stmt, err := tx.Prepare("INSERT INTO `roles`(`role_label`, `role_hierarchy`) VALUES (?, 20)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error preparing statement: %s", err.Error())
	}
	defer visitor_stmt.Close()

	_, err = visitor_stmt.Exec(service_models.VISITOR_ROLE)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error inserting visitor role: %s", err.Error())
	}

	return nil
}

func insertDefaultRoleGrants(tx *sql.Tx) error {
	all_grants := make([]string, 0)
	all_grants = append(all_grants, service_models.SuperAdminGrants...)
	all_grants = append(all_grants, service_models.AdminGrants...)
	all_grants = append(all_grants, service_models.VisitorGrants...)

	insert_grant_stmt, err := tx.Prepare("INSERT INTO `grants`(`grant_label`) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error preparing statement: %s", err.Error())
	}
	defer insert_grant_stmt.Close()

	insert_role_grant_stmt, err := tx.Prepare("INSERT INTO `role_grants`(`role`, `grant`) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error preparing statement: %s", err.Error())
	}
	defer insert_role_grant_stmt.Close()

	for _, grant := range all_grants {
		_, err = insert_grant_stmt.Exec(grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting grant: %s", err.Error())
		}
	}

	for _, visitor_grant := range service_models.VisitorGrants {
		_, err = insert_role_grant_stmt.Exec(service_models.VISITOR_ROLE, visitor_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting visitor role grant: %s", err.Error())
		}

		_, err = insert_role_grant_stmt.Exec(service_models.ADMIN_ROLE, visitor_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting admin role grant: %s", err.Error())
		}

		_, err = insert_role_grant_stmt.Exec(service_models.SUPER_ADMIN_ROLE, visitor_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting super admin role grant: %s", err.Error())
		}
	}

	for _, admin_grant := range service_models.AdminGrants {
		_, err = insert_role_grant_stmt.Exec(service_models.ADMIN_ROLE, admin_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting admin role grant: %s", err.Error())
		}

		_, err = insert_role_grant_stmt.Exec(service_models.SUPER_ADMIN_ROLE, admin_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting super admin role grant: %s", err.Error())
		}
	}

	for _, super_admin_grant := range service_models.SuperAdminGrants {
		_, err = insert_role_grant_stmt.Exec(service_models.SUPER_ADMIN_ROLE, super_admin_grant)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error inserting super admin role grant: %s", err.Error())
		}
	}

	return nil
}
