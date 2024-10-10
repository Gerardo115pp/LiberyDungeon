package database

import (
	"context"
	"database/sql"
	"fmt"
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery_categories_service/models"

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
func (categories_clusters_repo *CategoriesClustersMysql) GetCategoryCluster(ctx context.Context, category_id string) (dungeon_models.CategoryCluster, error) {
	var cluster dungeon_models.CategoryCluster

	stmt, err := categories_clusters_repo.db.PrepareContext(ctx, "SELECT `uuid`, `name`, `fs_path`, `filter_category`, `root_category` FROM `categories_clusters` WHERE `uuid`=(SELECT `cluster` FROM `categorys` WHERE `uuid`=?)")
	if err != nil {
		return cluster, err
	}

	err = stmt.QueryRowContext(ctx, category_id).Scan(&cluster.Uuid, &cluster.Name, &cluster.FsPath, &cluster.FilterCategory, &cluster.RootCategory)
	if err != nil {
		return cluster, err
	}

	stmt.Close()

	return cluster, nil
}

func (categories_clusters_repo *CategoriesClustersMysql) GetClusters(ctx context.Context) ([]dungeon_models.CategoryCluster, error) {
	var clusters []dungeon_models.CategoryCluster = make([]dungeon_models.CategoryCluster, 0)
	var err error

	rows, err := categories_clusters_repo.db.QueryContext(ctx, "SELECT `uuid`, `name`, `fs_path`, `filter_category`, `root_category` FROM `categories_clusters`")
	if err != nil {
		return clusters, err
	}

	for rows.Next() {
		var cluster dungeon_models.CategoryCluster
		err = rows.Scan(&cluster.Uuid, &cluster.Name, &cluster.FsPath, &cluster.FilterCategory, &cluster.RootCategory)
		if err != nil {
			return clusters, err
		}

		clusters = append(clusters, cluster)
	}

	rows.Close()

	return clusters, nil
}

func (categories_clusters_repo *CategoriesClustersMysql) InsertCluster(ctx context.Context, cluster dungeon_models.CategoryCluster, root_category dungeon_models.Category, filter_category dungeon_models.Category) error {
	tx, err := categories_clusters_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction because: %s", err.Error())
	}

	_, err = tx.ExecContext(ctx, "ALTER TABLE `categorys` DROP FOREIGN KEY `cluster_fk`")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to drop foreign key cluster_fk because: %s", err.Error())
	}

	category_insert_stmt, err := tx.PrepareContext(ctx, "INSERT INTO `categorys` (`uuid`, `name`, `fullpath`, `cluster`) VALUES (?, ?, ?, ?)")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to prepare category insert statement because: %s", err.Error())
	}
	defer category_insert_stmt.Close()

	_, err = category_insert_stmt.ExecContext(ctx, root_category.Uuid, root_category.Name, root_category.Fullpath, root_category.Cluster)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to insert root category because: %s", err.Error())
	}

	insert_filter_stmt, err := tx.PrepareContext(ctx, "INSERT INTO `categorys` (`uuid`, `name`, `fullpath`, `cluster`, `parent`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to prepare filter category insert statement because: %s", err.Error())
	}
	defer insert_filter_stmt.Close()

	_, err = insert_filter_stmt.ExecContext(ctx, filter_category.Uuid, filter_category.Name, filter_category.Fullpath, filter_category.Cluster, filter_category.Parent)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to insert filter category because: %s", err.Error())
	}

	insert_cluster_stmt, err := tx.PrepareContext(ctx, "INSERT INTO `categories_clusters` (`uuid`, `name`, `fs_path`, `filter_category`, `root_category`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return err
	}
	defer insert_cluster_stmt.Close()

	_, err = insert_cluster_stmt.ExecContext(ctx, cluster.Uuid, cluster.Name, cluster.FsPath, cluster.FilterCategory, cluster.RootCategory)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return err
	}

	_, err = tx.ExecContext(ctx, "ALTER TABLE `categorys` ADD CONSTRAINT `cluster_fk` FOREIGN KEY (`cluster`) REFERENCES `categories_clusters`(`uuid`) ON DELETE CASCADE")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Failed to rollback transaction because: %s", err.Error())
		}

		return fmt.Errorf("Failed to add foreign key cluster_fk because: %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction because: %s", err.Error())
	}

	return nil
}

func (categories_clusters_repo *CategoriesClustersMysql) UpdateCluster(ctx context.Context, cluster dungeon_models.CategoryCluster) error {
	stmt, err := categories_clusters_repo.db.Prepare("UPDATE `categories_clusters` SET `name`=?, `fs_path`=?, `filter_category`=? WHERE `uuid`=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, cluster.Name, cluster.FsPath, cluster.FilterCategory, cluster.Uuid)
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func (categories_clusters_repo *CategoriesClustersMysql) DeleteCluster(ctx context.Context, cluster_id string) *dungeon_models.LabeledError {
	var labeled_err *dungeon_models.LabeledError
	tx, err := categories_clusters_repo.db.BeginTx(ctx, nil)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "On DeleteCluster, failed to start transaction", dungeon_models.ErrDB_CouldNotCreateTX)
		return labeled_err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM `categories_clusters` WHERE `uuid`=?", cluster_id)
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "On DeleteCluster, failed to delete cluster", service_models.ErrDB_CouldNotFindCategoryCluster)
		return labeled_err
	}

	err = tx.Commit()
	if err != nil {
		labeled_err = dungeon_models.NewLabeledError(err, "On DeleteCluster, failed to commit transaction", dungeon_models.ErrDB_FailedToCommitTX)
		return labeled_err
	}

	return nil
}
