package database

import (
	"context"
	"database/sql"
	dungeon_models "libery-dungeon-libs/models"

	_ "github.com/go-sql-driver/mysql"
)

type CategoriesClustersMysql struct {
	db *sql.DB
}

func NewCategoriesClustersMysql() (*CategoriesClustersMysql, error) {
	db, err := sql.Open("mysql", createDSN())
	if err != nil {
		return nil, err
	}

	return &CategoriesClustersMysql{db: db}, nil
}

func (categories_clusters_repo *CategoriesClustersMysql) GetClusterByID(ctx context.Context, cluster_id string) (dungeon_models.CategoryCluster, error) {
	var cluster dungeon_models.CategoryCluster

	stmt, err := categories_clusters_repo.db.Prepare("SELECT `uuid`, `name`, `fs_path`, `filter_category`, `root_category` FROM `categories_clusters` WHERE `uuid`=?")
	if err != nil {
		return cluster, err
	}

	err = stmt.QueryRowContext(ctx, cluster_id).Scan(&cluster.Uuid, &cluster.Name, &cluster.FsPath, &cluster.FilterCategory, &cluster.RootCategory)
	if err != nil {
		return cluster, err
	}

	stmt.Close()

	return cluster, nil
}

// Returns the cluster that contains the category of which the id is passed as parameter
func (categories_clusters_repo *CategoriesClustersMysql) GetCategoryCluster(ctx context.Context, category_id string) (*dungeon_models.CategoryCluster, error) {
	var cluster dungeon_models.CategoryCluster

	stmt, err := categories_clusters_repo.db.PrepareContext(ctx, "SELECT `uuid`, `name`, `fs_path`, `filter_category`, `root_category` FROM `categories_clusters` WHERE `uuid`=(SELECT `cluster` FROM `categorys` WHERE `uuid`=?)")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, category_id).Scan(&cluster.Uuid, &cluster.Name, &cluster.FsPath, &cluster.FilterCategory, &cluster.RootCategory)
	if err != nil {
		return nil, err
	}

	stmt.Close()

	return &cluster, nil
}
