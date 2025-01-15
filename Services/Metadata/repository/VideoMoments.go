package repository

import (
	"context"
	video_moment_models "libery-metadata-service/models/video_moments"
)

type VideoMomentsRepository interface {
	AddVideoCTX(ctx context.Context, video video_moment_models.Video) error
	AddVideo(video video_moment_models.Video) error
	AddVideoMomentCTX(ctx context.Context, video_moment video_moment_models.VideoMoment) (int, error)
	AddVideoMoment(video_moment video_moment_models.VideoMoment) (int, error)
	DeleteVideoMomentCTX(ctx context.Context, video_moment video_moment_models.VideoMoment) error
	DeleteVideoMoment(video_moment video_moment_models.VideoMoment) error
	GetVideoCTX(ctx context.Context, video_uuid string, cluster_uuid string) (*video_moment_models.Video, error)
	GetVideo(video_uuid string, cluster_uuid string) (*video_moment_models.Video, error)
	GetVideoMomentCTX(ctx context.Context, moment_id int) (*video_moment_models.VideoMoment, error)
	GetVideoMoment(moment_id int) (*video_moment_models.VideoMoment, error)
	GetVideoMomentsCTX(ctx context.Context, video *video_moment_models.Video) ([]video_moment_models.VideoMoment, error)
	GetVideoMoments(video *video_moment_models.Video) ([]video_moment_models.VideoMoment, error)
	GetClusterVideosCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.Video, error)
	GetClusterVideos(cluster_uuid string) ([]video_moment_models.Video, error)
	GetClusterMomentsCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.VideoMoment, error)
	GetClusterMoments(cluster_uuid string) ([]video_moment_models.VideoMoment, error)
	GetClusterVideoMomentsCTX(ctx context.Context, cluster_uuid string) ([]video_moment_models.VideoMoments, error)
	GetClusterVideoMoments(cluster_uuid string) ([]video_moment_models.VideoMoments, error)
	UpdateVideoMomentCTX(ctx context.Context, video_moment video_moment_models.VideoMoment) error
	UpdateVideoMoment(video_moment video_moment_models.VideoMoment) error
}

var VideoMomentsRepo VideoMomentsRepository

func SetVideoMomentsRepo(video_moments_repo VideoMomentsRepository) {
	VideoMomentsRepo = video_moments_repo
}
