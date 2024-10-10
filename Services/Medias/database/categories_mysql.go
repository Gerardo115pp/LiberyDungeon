package database

import (
	"context"
	"database/sql"
	dungeon_models "libery-dungeon-libs/models"

	_ "github.com/go-sql-driver/mysql"
)

type CategoriesMysql struct {
	db *sql.DB
}

func NewCategoriesMysql() (*CategoriesMysql, error) {
	db, err := sql.Open("mysql", createDSN())
	if err != nil {
		return nil, err
	}

	return &CategoriesMysql{db: db}, nil
}

func (db *CategoriesMysql) GetCategoryByID(ctx context.Context, category_id string) (dungeon_models.Category, error) {
	var category dungeon_models.Category

	stmt, err := db.db.Prepare("SELECT uuid, name, fullpath, parent, cluster FROM categorys WHERE uuid = ?")
	if err != nil {
		return category, err
	}

	err = stmt.QueryRowContext(ctx, category_id).Scan(&category.Uuid, &category.Name, &category.Fullpath, &category.Parent, &category.Cluster)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (db *CategoriesMysql) Close() error {
	return db.db.Close()
}
