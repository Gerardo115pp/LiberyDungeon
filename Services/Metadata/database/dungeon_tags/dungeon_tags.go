package dungeon_tags

import "database/sql"

type DungeonTagsDB struct {
	db_conn *sql.DB
}

func NewDungeonTagsDB() (*DungeonTagsDB, error) {
	var users_db *DungeonTagsDB = new(DungeonTagsDB)

	db, err := openDB()
	if err != nil {
		return nil, err
	}

	users_db.db_conn = db

	db.Exec("PRAGMA foreign_keys = ON")

	return users_db, nil
}
