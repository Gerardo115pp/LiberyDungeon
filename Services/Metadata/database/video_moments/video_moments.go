package video_moments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"libery-dungeon-libs/libs/dungeon_sqlite_opener"
	app_config "libery-metadata-service/Config"
	video_moment_models "libery-metadata-service/models/video_moments"
)

type VideoMomentsDB struct {
	db_conn *sql.DB
}

func NewVideoMomentsDB() *VideoMomentsDB {
	var video_moments_db *VideoMomentsDB = new(VideoMomentsDB)

	var sqlite_opener *dungeon_sqlite_opener.DungeonSqliteOpener
	sqlite_opener = dungeon_sqlite_opener.NewDungeonSqliteOpener("video_moments.db", "video_moments.sql", app_config.OPERATION_DATA_PATH)

	db, err := sqlite_opener.OpenDB(true)
	if err != nil {
		panic(err)
	}

	video_moments_db.db_conn = db

	return video_moments_db
}

func (video_moments_db VideoMomentsDB) AddVideoCTX(ctx context.Context, video video_moment_models.Video) error {
	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "INSERT INTO videos (video_uuid, video_cluster) VALUES (?, ?)")
	if err != nil {
		return errors.Join(fmt.Errorf("In database/video_moments/video_moments.AddVideoCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, video.VideoUUID, video.VideoCluster)
	if err != nil {
		return errors.Join(fmt.Errorf("In database/video_moments/video_moments.AddVideoCTX: While executing statement."), err)
	}

	return nil
}

func (video_moments_db VideoMomentsDB) AddVideo(video video_moment_models.Video) error {
	return video_moments_db.AddVideoCTX(context.Background(), video)
}

func (video_moments_db VideoMomentsDB) AddVideoMomentCTX(ctx context.Context, video_moment video_moment_models.VideoMoment) error {
	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "INSERT INTO video_moments (video_uuid, moment_time) VALUES (?, ?)")
	if err != nil {
		return errors.Join(fmt.Errorf("In database/video_moments/video_moments.AddVideoMomentCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, video_moment.VideoUUID, video_moment.MomentTime)
	if err != nil {
		return errors.Join(fmt.Errorf("In database/video_moments/video_moments.AddVideoMomentCTX: While executing statement."), err)
	}

	return nil
}

func (video_moments_db VideoMomentsDB) AddVideoMoment(video_moment video_moment_models.VideoMoment) error {
	return video_moments_db.AddVideoMomentCTX(context.Background(), video_moment)
}

func (video_moments_db VideoMomentsDB) GetVideoCTX(ctx context.Context, video_uuid string, cluster_uuid string) (*video_moment_models.Video, error) {
	var video video_moment_models.Video

	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "SELECT video_uuid, video_cluster FROM videos WHERE video_uuid = ? AND video_cluster = ?")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, video_uuid, cluster_uuid).Scan(&video.VideoUUID, &video.VideoCluster)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoCTX: While executing statement."), err)
	}

	return &video, nil
}

func (video_moments_db VideoMomentsDB) GetVideo(video_uuid string, cluster_uuid string) (*video_moment_models.Video, error) {
	return video_moments_db.GetVideoCTX(context.Background(), video_uuid, cluster_uuid)
}

func (video_moments_db VideoMomentsDB) GetVideoMomentCTX(ctx context.Context, video *video_moment_models.Video, moment_id int) (*video_moment_models.VideoMoment, error) {
	var video_moment video_moment_models.VideoMoment

	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "SELECT `id`, `video_uuid`, `moment_time` FROM `video_moments` WHERE `video_uuid` = ? AND `id` = ?")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoMomentCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, video.VideoUUID, moment_id).Scan(&video_moment.ID, &video_moment.VideoUUID, &video_moment.MomentTime)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoMomentCTX: While executing statement."), err)
	}

	return &video_moment, nil
}

func (video_moments_db VideoMomentsDB) GetVideoMoment(video *video_moment_models.Video, moment_id int) (*video_moment_models.VideoMoment, error) {
	return video_moments_db.GetVideoMomentCTX(context.Background(), video, moment_id)
}

func (video_moments_db VideoMomentsDB) GetVideoMomentsCTX(ctx context.Context, video *video_moment_models.Video) ([]video_moment_models.VideoMoment, error) {
	var video_moments []video_moment_models.VideoMoment

	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "SELECT `id`, `video_uuid`, `moment_time` FROM `video_moments` WHERE `video_uuid` = ?")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoMomentsCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, video.VideoUUID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoMomentsCTX: While executing statement."), err)
	}
	defer rows.Close()

	for rows.Next() {
		var video_moment video_moment_models.VideoMoment

		err = rows.Scan(&video_moment.ID, &video_moment.VideoUUID, &video_moment.MomentTime)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetVideoMomentsCTX: While scanning rows."), err)
		}

		video_moments = append(video_moments, video_moment)
	}

	return video_moments, nil
}

func (video_moments_db VideoMomentsDB) GetVideoMoments(video *video_moment_models.Video) ([]video_moment_models.VideoMoment, error) {
	return video_moments_db.GetVideoMomentsCTX(context.Background(), video)
}

func (video_moments_db VideoMomentsDB) GetClusterVideosCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.Video, error) {
	var videos []video_moment_models.Video

	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "SELECT `video_uuid`, `video_cluster` FROM `videos` WHERE `video_cluster` = ?")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterVideosCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, cluster_uuid)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterVideosCTX: While executing statement."), err)
	}
	defer rows.Close()

	for rows.Next() {
		var video video_moment_models.Video

		err = rows.Scan(&video.VideoUUID, &video.VideoCluster)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterVideosCTX: While scanning rows."), err)
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (video_moments_db VideoMomentsDB) GetClusterVideos(cluster_uuid string) ([]video_moment_models.Video, error) {
	return video_moments_db.GetClusterVideosCTX(context.Background(), cluster_uuid)
}

func (video_moments_db VideoMomentsDB) GetClusterMomentsCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.VideoMoment, error) {
	var video_moments []video_moment_models.VideoMoment

	stmt, err := video_moments_db.db_conn.PrepareContext(ctx, "SELECT `id`, `video_uuid`, `moment_time` FROM `video_moments` WHERE `video_uuid` IN (SELECT `video_uuid` FROM `videos` WHERE `video_cluster` = ?)")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterMomentsCTX: While preparing statement."), err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, cluster_uuid)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterMomentsCTX: While executing statement."), err)
	}
	defer rows.Close()

	for rows.Next() {
		var video_moment video_moment_models.VideoMoment

		err = rows.Scan(&video_moment.ID, &video_moment.VideoUUID, &video_moment.MomentTime)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterMomentsCTX: While scanning rows."), err)
		}

		video_moments = append(video_moments, video_moment)
	}

	return video_moments, nil
}

func (video_moments_db VideoMomentsDB) GetClusterMoments(cluster_uuid string) ([]video_moment_models.VideoMoment, error) {
	return video_moments_db.GetClusterMomentsCTX(context.Background(), cluster_uuid)
}

func (video_moments_db VideoMomentsDB) GetClusterVideoMomentsCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.VideoMoments, error) {
	var video_moments []video_moment_models.VideoMoments

	videos, err := video_moments_db.GetClusterVideosCTX(ctx, cluster_uuid)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterVideoMomentsCTX: While getting videos."), err)
	}

	for _, video := range videos {
		var video_moments_instance video_moment_models.VideoMoments

		video_moments_instance.Video = video

		moments, err := video_moments_db.GetVideoMomentsCTX(ctx, &video)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("In database/video_moments/video_moments.GetClusterVideoMomentsCTX: While getting video moments."), err)
		}

		video_moments_instance.Moments = moments

		video_moments = append(video_moments, video_moments_instance)
	}

	return video_moments, nil
}

func (video_moments_db VideoMomentsDB) GetClusterVideoMoments(cluster_uuid string) ([]video_moment_models.VideoMoments, error) {
	return video_moments_db.GetClusterVideoMomentsCTX(context.Background(), cluster_uuid)
}
