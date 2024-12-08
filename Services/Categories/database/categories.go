package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"path/filepath"

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

func (categories_repo *CategoriesMysql) GetCategoryChildsByID(ctx context.Context, category_id string) ([]dungeon_models.ChildCategory, error) {
	var childs []dungeon_models.ChildCategory = make([]dungeon_models.ChildCategory, 0)
	var err error

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `fullpath`, `category_thumbnail` FROM `categorys` WHERE `parent`=? ORDER BY `name`")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var child dungeon_models.ChildCategory
		var null_string_category_thumbnail sql.NullString

		err = rows.Scan(&child.Uuid, &child.Name, &child.Fullpath, &null_string_category_thumbnail)
		if err != nil {
			return nil, err
		}

		child.CategoryThumbnail = "" // The default value but explicit is better than implicit

		if null_string_category_thumbnail.Valid {
			child.CategoryThumbnail = null_string_category_thumbnail.String
		}

		childs = append(childs, child)
	}

	return childs, nil
}

func (categories_repo *CategoriesMysql) GetCategoryMedias(ctx context.Context, category_id string) ([]dungeon_models.Media, error) {
	var category_medias []dungeon_models.Media = make([]dungeon_models.Media, 0)

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `last_seen`, `main_category`, `media_thumbnail`, `type`, `downloaded_from` FROM `medias` WHERE `main_category`=? ORDER BY `name`")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var media dungeon_models.Media
		var time_reciever sql.NullTime
		var media_thumbnail_reciever sql.NullString
		var null_int_reciever sql.NullInt64

		err = rows.Scan(&media.Uuid, &media.Name, &time_reciever, &media.MainCategory, &media_thumbnail_reciever, &media.Type, &null_int_reciever)
		if err != nil {
			return nil, err
		}
		if time_reciever.Valid {
			media.LastSeen = time_reciever.Time
		}

		media.MediaThumbnail = ""
		if media_thumbnail_reciever.Valid {
			media.MediaThumbnail = media_thumbnail_reciever.String
		}

		if null_int_reciever.Valid {
			media.DownloadedFrom = null_int_reciever.Int64
		}

		category_medias = append(category_medias, media)
	}

	return category_medias, nil
}

func (categories_repo *CategoriesMysql) GetMediaIdentity(ctx context.Context, media_uuid string) (*dungeon_models.MediaIdentity, error) {
	var media_identity *dungeon_models.MediaIdentity

	stmt, err := categories_repo.db.PrepareContext(ctx, `
		SELECT
			m.uuid, m.name, m.last_seen, m.main_category, m.media_thumbnail, m.type, m.downloaded_from,
			c.uuid, c.fullpath,
			cc.uuid, cc.fs_path
		FROM medias m
		INNER JOIN categorys c ON m.main_category=c.uuid
		INNER JOIN categories_clusters cc ON c.cluster=cc.uuid
		WHERE m.uuid=?
	`)
	if err != nil {
		return media_identity, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentity: Error preparing statement"))
	}
	defer stmt.Close()

	var time_reciever sql.NullTime
	var downloaded_from_reciever sql.NullInt64

	row := stmt.QueryRowContext(ctx, media_uuid)

	media_identity = new(dungeon_models.MediaIdentity)
	media_identity.Media = new(dungeon_models.Media)

	var media_thumbnail_reciever sql.NullString

	err = row.Scan(
		&media_identity.Media.Uuid,
		&media_identity.Media.Name,
		&time_reciever,
		&media_identity.Media.MainCategory,
		&media_thumbnail_reciever,
		&media_identity.Media.Type,
		&downloaded_from_reciever,
		&media_identity.CategoryUUID,
		&media_identity.CategoryPath,
		&media_identity.ClusterUUID,
		&media_identity.ClusterPath,
	)
	if err != nil {
		return media_identity, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentity: Error scanning row"))
	}

	if time_reciever.Valid {
		media_identity.Media.LastSeen = time_reciever.Time
	}

	media_identity.Media.MediaThumbnail = ""

	if media_thumbnail_reciever.Valid {
		media_identity.Media.MediaThumbnail = media_thumbnail_reciever.String
	}

	if downloaded_from_reciever.Valid {
		media_identity.Media.DownloadedFrom = downloaded_from_reciever.Int64
	}

	return media_identity, nil
}

func (categories_repo *CategoriesMysql) GetMediaIdentityList(ctx context.Context, media_uuids []string) ([]dungeon_models.MediaIdentity, error) {
	var media_identities []dungeon_models.MediaIdentity = make([]dungeon_models.MediaIdentity, len(media_uuids))

	stmt, err := categories_repo.db.PrepareContext(ctx, `
		SELECT
			m.uuid, m.name, m.last_seen, m.main_category, m.media_thumbnail, m.type, m.downloaded_from,
			c.uuid, c.fullpath,
			cc.uuid, cc.fs_path
		FROM medias m
		INNER JOIN categorys c ON m.main_category=c.uuid
		INNER JOIN categories_clusters cc ON c.cluster=cc.uuid
		WHERE m.uuid = ?
	`)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentityList: Error preparing statement"))
	}
	defer stmt.Close()

	var time_reciever sql.NullTime
	var media_thumbnail_reciever sql.NullString
	var downloaded_from_reciever sql.NullInt64

	for h, media_uuid := range media_uuids {
		row := stmt.QueryRowContext(ctx, media_uuid)

		media_identity := new(dungeon_models.MediaIdentity)
		media_identity.Media = new(dungeon_models.Media)

		err = row.Scan(
			&media_identity.Media.Uuid,
			&media_identity.Media.Name,
			&time_reciever,
			&media_identity.Media.MainCategory,
			&media_thumbnail_reciever,
			&media_identity.Media.Type,
			&downloaded_from_reciever,
			&media_identity.CategoryUUID,
			&media_identity.CategoryPath,
			&media_identity.ClusterUUID,
			&media_identity.ClusterPath,
		)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetMediaIdentityList: Error scanning row"))
		}

		if time_reciever.Valid {
			media_identity.Media.LastSeen = time_reciever.Time
		}

		media_identity.Media.MediaThumbnail = ""

		if media_thumbnail_reciever.Valid {
			media_identity.Media.MediaThumbnail = media_thumbnail_reciever.String
		}

		if downloaded_from_reciever.Valid {
			media_identity.Media.DownloadedFrom = downloaded_from_reciever.Int64
		}

		media_identities[h] = *media_identity
	}

	return media_identities, nil
}

func (categories_repo *CategoriesMysql) GetCategoryMediaIdentities(ctx context.Context, category_uuid string) ([]dungeon_models.MediaIdentity, error) {
	var media_identities []dungeon_models.MediaIdentity = make([]dungeon_models.MediaIdentity, 0)

	stmt, err := categories_repo.db.PrepareContext(ctx, `
		SELECT
			m.uuid, m.name, m.last_seen, m.main_category, m.media_thumbnail, m.type, m.downloaded_from,
			c.uuid, c.fullpath,
			cc.uuid, cc.fs_path
		FROM medias m
		INNER JOIN categorys c ON m.main_category=c.uuid
		INNER JOIN categories_clusters cc ON c.cluster=cc.uuid
		WHERE c.uuid=?
	`)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetCategoryMediaIdentities: Error preparing statement"))
	}
	defer stmt.Close()

	var time_reciever sql.NullTime
	var media_thumbnail_reciever sql.NullString
	var downloaded_from_reciever sql.NullInt64

	rows, err := stmt.QueryContext(ctx, category_uuid)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.GetCategoryMediaIdentities: Error querying"))
	}
	defer rows.Close()

	for rows.Next() {
		var media_identity dungeon_models.MediaIdentity
		media_identity.Media = new(dungeon_models.Media)

		err = rows.Scan(
			&media_identity.Media.Uuid,
			&media_identity.Media.Name,
			&time_reciever,
			&media_identity.Media.MainCategory,
			&media_thumbnail_reciever,
			&media_identity.Media.Type,
			&downloaded_from_reciever,
			&media_identity.CategoryUUID,
			&media_identity.CategoryPath,
			&media_identity.ClusterUUID,
			&media_identity.ClusterPath,
		)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("In CategoriesService.CategoriesMysql.Join: Error scanning row"))
		}

		if time_reciever.Valid {
			media_identity.Media.LastSeen = time_reciever.Time
		}

		media_identity.Media.MediaThumbnail = ""

		if media_thumbnail_reciever.Valid {
			media_identity.Media.MediaThumbnail = media_thumbnail_reciever.String
		}

		media_identities = append(media_identities, media_identity)
	}

	return media_identities, nil
}

func (categories_repo *CategoriesMysql) GetCategoryContent(ctx context.Context, category_id string) (*dungeon_models.CategoryLeaf, error) {
	var category_leaf *dungeon_models.CategoryLeaf = new(dungeon_models.CategoryLeaf)

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `fullpath`, `parent`, `cluster`, `category_thumbnail` FROM `categorys` WHERE `uuid`=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, category_id)

	var empty_parent_reciever sql.NullString
	var null_string_category_thumbnail sql.NullString

	err = row.Scan(&category_leaf.Uuid, &category_leaf.Name, &category_leaf.Fullpath, &empty_parent_reciever, &category_leaf.Cluster, &null_string_category_thumbnail)
	if err != nil {
		return nil, err
	}

	category_leaf.Parent = "" // The default value but explicit is better than implicit

	if empty_parent_reciever.Valid {
		category_leaf.Parent = empty_parent_reciever.String
	}

	category_leaf.CategoryThumbnail = ""

	if null_string_category_thumbnail.Valid {
		category_leaf.CategoryThumbnail = null_string_category_thumbnail.String
	}

	category_leaf.InnerCategories, err = categories_repo.GetCategoryChildsByID(ctx, category_id)
	if err != nil {
		return nil, err
	}

	category_leaf.Content, err = categories_repo.GetCategoryMedias(ctx, category_id)
	if err != nil {
		return nil, err
	}

	return category_leaf, nil
}

func (categories_repo *CategoriesMysql) GetCategory(ctx context.Context, category_id string) (dungeon_models.Category, error) {
	var category dungeon_models.Category

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `fullpath`, `parent`, `cluster`, `category_thumbnail` FROM `categorys` WHERE uuid=?")
	if err != nil {
		return category, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, category_id)

	var empty_parent_reciever sql.NullString
	var null_string_category_thumbnail sql.NullString

	err = row.Scan(&category.Uuid, &category.Name, &category.Fullpath, &empty_parent_reciever, &category.Cluster, &null_string_category_thumbnail)
	if err != nil {
		return category, err
	}

	if empty_parent_reciever.Valid {
		category.Parent = empty_parent_reciever.String
	}

	category.CategoryThumbnail = ""

	if null_string_category_thumbnail.Valid {
		category.CategoryThumbnail = null_string_category_thumbnail.String
	}

	return category, nil
}

func (categories_repo *CategoriesMysql) GetCategories(ctx context.Context, category_ids []string) ([]dungeon_models.Category, error) {
	var categories []dungeon_models.Category = make([]dungeon_models.Category, 0)

	for _, category_id := range category_ids {
		category, err := categories_repo.GetCategory(ctx, category_id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (categories_repo *CategoriesMysql) GetClusterCategories(ctx context.Context, cluster_id string) ([]dungeon_models.Category, error) {
	var categories []dungeon_models.Category = make([]dungeon_models.Category, 0)

	stmt, err := categories_repo.db.PrepareContext(ctx, "SELECT `uuid`, `name`, `fullpath`, `parent`, `cluster`, `category_thumbnail` FROM `categorys` WHERE `cluster`=? ORDER BY `name`")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, cluster_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category dungeon_models.Category
		var empty_parent_reciever sql.NullString
		var null_string_category_thumbnail sql.NullString

		err = rows.Scan(&category.Uuid, &category.Name, &category.Fullpath, &empty_parent_reciever, &category.Cluster, &null_string_category_thumbnail)
		if err != nil {
			return nil, err
		}

		if empty_parent_reciever.Valid {
			category.Parent = empty_parent_reciever.String
		}

		if null_string_category_thumbnail.Valid {
			category.CategoryThumbnail = null_string_category_thumbnail.String
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (categories_repo *CategoriesMysql) GetAllCategories(ctx context.Context) ([]dungeon_models.Category, error) {
	var categories []dungeon_models.Category = make([]dungeon_models.Category, 0)

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `fullpath`, `parent`, `cluster`, `category_thumbnail` FROM `categorys` ORDER BY `name`")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category dungeon_models.Category
		var empty_parent_reciever sql.NullString
		var null_string_category_thumbnail sql.NullString

		err = rows.Scan(&category.Uuid, &category.Name, &category.Fullpath, &empty_parent_reciever, &category.Cluster, &null_string_category_thumbnail)
		if err != nil {
			return nil, err
		}

		if empty_parent_reciever.Valid {
			category.Parent = empty_parent_reciever.String
		}

		if null_string_category_thumbnail.Valid {
			category.CategoryThumbnail = null_string_category_thumbnail.String
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (categories_repo *CategoriesMysql) GetCategoryContentByFullpath(ctx context.Context, category_path string, category_cluster string) (*dungeon_models.Category, error) {
	var category *dungeon_models.Category = new(dungeon_models.Category)

	stmt, err := categories_repo.db.Prepare("SELECT `uuid`, `name`, `fullpath`, `parent`, `cluster`, `category_thumbnail` FROM `categorys` WHERE `fullpath`=? AND `cluster`=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, category_path, category_cluster)

	var empty_parent_reciever sql.NullString
	var null_string_category_thumbnail sql.NullString

	err = row.Scan(&category.Uuid, &category.Name, &category.Fullpath, &empty_parent_reciever, &category.Cluster, &null_string_category_thumbnail)
	if err != nil {
		return nil, err
	}

	if empty_parent_reciever.Valid {
		category.Parent = empty_parent_reciever.String
	}

	if null_string_category_thumbnail.Valid {
		category.CategoryThumbnail = null_string_category_thumbnail.String
	}

	return category, nil
}

func (categories_repo *CategoriesMysql) GetCategoryFSBranch(ctx context.Context, category_id string) ([]dungeon_models.MediaWeakIdentity, error) {
	var medias []dungeon_models.MediaWeakIdentity = make([]dungeon_models.MediaWeakIdentity, 0)

	query := `
		WITH RECURSIVE category_tree AS (
			SELECT uuid, fullpath, parent 
			FROM categorys 
			WHERE uuid=?

			UNION ALL	

			SELECT c.uuid, c.fullpath, c.parent
			FROM category_tree p 
			JOIN categorys c 
			ON p.uuid = c.parent
		)
		SELECT 
			IFNULL(m.uuid, '') as media_uuid, 
			IFNULL(m.name, '') as media_name, 
			ct.uuid as category_uuid, 
			ct.fullpath as category_path 
		FROM 
			category_tree ct 
		LEFT JOIN 
			medias m 
		ON 
			ct.uuid = m.main_category`

	stmt, err := categories_repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var media_identity dungeon_models.MediaWeakIdentity

		err = rows.Scan(&media_identity.MediaUUID, &media_identity.MediaName, &media_identity.CategoryUUID, &media_identity.CategoryPath)
		if err != nil {
			return nil, err
		}

		medias = append(medias, media_identity)
	}

	return medias, nil
}

func (categories_repo *CategoriesMysql) DeleteCategoryMedias(ctx context.Context, medias []dungeon_models.Media) error {
	var err error

	tx, err := categories_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := categories_repo.db.Prepare("DELETE FROM `medias` WHERE `uuid`=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, media := range medias {
		_, err = stmt.ExecContext(ctx, media.Uuid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (categories_repo *CategoriesMysql) DeleteCategory(ctx context.Context, category_id string) error {
	var err error

	tx, err := categories_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := categories_repo.db.Prepare("DELETE FROM `categorys` WHERE `uuid`=? LIMIT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, category_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (categories_repo *CategoriesMysql) IsCategoryEmpty(ctx context.Context, category_id string) (bool, error) {
	var err error

	medias_count_stmt, err := categories_repo.db.Prepare("SELECT COUNT(*) FROM `medias` WHERE `main_category`=?")
	if err != nil {
		return false, err
	}
	defer medias_count_stmt.Close()

	category_count_stmt, err := categories_repo.db.Prepare("SELECT COUNT(*) FROM `categorys` WHERE `parent`=?")
	if err != nil {
		return false, err
	}
	defer category_count_stmt.Close()

	var medias_count int
	var category_count int

	row := medias_count_stmt.QueryRowContext(ctx, category_id)

	err = row.Scan(&medias_count)
	if err != nil {
		return false, err
	}

	row = category_count_stmt.QueryRowContext(ctx, category_id)

	err = row.Scan(&category_count)
	if err != nil {
		return false, err
	}

	return medias_count == 0 && category_count == 0, nil
}

func (categories_repo *CategoriesMysql) InsertCategory(ctx context.Context, category dungeon_models.Category) error {
	var err error

	stmt, err := categories_repo.db.Prepare("INSERT INTO `categorys` (`uuid`, `name`, `fullpath`, `parent`, `cluster`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, category.Uuid, category.Name, category.Fullpath, category.Parent, category.Cluster)
	if err != nil {
		return err
	}

	return nil
}

func (categories_repo *CategoriesMysql) UpdateMedia(ctx context.Context, media dungeon_models.Media) error {
	var err error

	stmt, err := categories_repo.db.Prepare("UPDATE `medias` SET `name`=?, `last_seen`=?, `main_category`=?, `media_thumbnail`=? WHERE `uuid`=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var media_thumbnail *string

	if media.MediaThumbnail != "" {
		media_thumbnail = &media.MediaThumbnail
	}

	_, err = stmt.ExecContext(ctx, media.Name, media.LastSeen, media.MainCategory, media_thumbnail, media.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func (categories_repo *CategoriesMysql) UpdateMedias(ctx context.Context, medias []dungeon_models.Media) error {
	var err error

	tx, err := categories_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, media := range medias {
		err = categories_repo.UpdateMedia(ctx, media)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (categories_repo *CategoriesMysql) UpdateCategoryName(ctx context.Context, category dungeon_models.Category, new_name string) error {
	var tx *sql.Tx
	var err error

	tx, err = categories_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rename_category_stmt, err := tx.Prepare("UPDATE `categorys` SET `name`=? WHERE `uuid`=?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rename_category_stmt.Close()

	_, err = rename_category_stmt.ExecContext(ctx, new_name, category.Uuid)
	if err != nil {
		return tx.Rollback()
	}

	rename_categories_fullpath_stmt, err := tx.Prepare("UPDATE `categorys` SET `fullpath`=REPLACE(`fullpath`, ?, ?) WHERE `fullpath` LIKE ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rename_categories_fullpath_stmt.Close()

	var new_fullpath string = helpers.RenameFsPath(category.Fullpath, new_name)

	// if new_fullpath != "" {
	// 	tx.Rollback()
	// 	return fmt.Errorf("Testing: Query would be 'UPDATE `categorys` SET `fullpath`=REPLACE(`fullpath`, '%s', '%s') WHERE `fullpath` LIKE '%s'", category.Fullpath, new_fullpath, category.Fullpath+"%")
	// }

	_, err = rename_categories_fullpath_stmt.ExecContext(ctx, category.Fullpath, new_fullpath, category.Fullpath+"%")
	if err != nil {
		return tx.Rollback()
	}

	err = tx.Commit()

	return err
}

func (categories_repo *CategoriesMysql) UpdateCategoryParent(ctx context.Context, category dungeon_models.Category, new_parent dungeon_models.Category) error {
	var tx *sql.Tx
	var err error

	var new_category_path string = fmt.Sprintf("%s/", filepath.Join(new_parent.Fullpath, category.Name))

	tx, err = categories_repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	update_category_stmt, err := tx.Prepare("UPDATE `categorys` SET `parent`=?, `fullpath`=? WHERE `uuid`=?")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Error rolling back transaction: %s. Oiginal Error was: %s", rollback_err, err.Error())
		}

		return fmt.Errorf("Error preparing update category statement: %s", err.Error())
	}
	defer update_category_stmt.Close()

	_, err = update_category_stmt.ExecContext(ctx, new_parent.Uuid, new_category_path, category.Uuid)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Error rolling back transaction: %s. Oiginal Error was: %s", rollback_err, err.Error())
		}
		return fmt.Errorf("Error executing update category statement: %s", err.Error())
	}

	update_childs_stmt, err := tx.Prepare("UPDATE `categorys` SET `fullpath`=REPLACE(`fullpath`, ?, ?) WHERE `fullpath` LIKE ?")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Error rolling back transaction: %s. Oiginal Error was: %s", rollback_err, err.Error())
		}

		return fmt.Errorf("Error preparing update childs statement: %s", err.Error())
	}
	defer update_childs_stmt.Close()

	_, err = update_childs_stmt.ExecContext(ctx, category.Fullpath, new_category_path, category.Fullpath+"%")
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			return fmt.Errorf("Error rolling back transaction: %s. Oiginal Error was: %s", rollback_err, err.Error())
		}
		return fmt.Errorf("Error executing update childs statement: %s", err.Error())
	}

	return tx.Commit()
}

func (categories_repo *CategoriesMysql) UpdateCategoryThumbnail(ctx context.Context, category_uuid, new_thumbnail string) error {
	var err error

	stmt, err := categories_repo.db.Prepare("UPDATE `categorys` SET `category_thumbnail`=? WHERE `uuid`=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, new_thumbnail, category_uuid)
	if err != nil {
		return err
	}

	return nil
}
