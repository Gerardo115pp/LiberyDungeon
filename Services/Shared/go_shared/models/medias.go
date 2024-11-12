package dungeon_models

import (
	"fmt"
	"libery-dungeon-libs/helpers"
	"path/filepath"
	"time"
)

type Media struct {
	Uuid           string    `json:"uuid"`
	Name           string    `json:"name"`
	LastSeen       time.Time `json:"last_seen"`
	MainCategory   string    `json:"main_category"`
	MediaThumbnail string    `json:"media_thumbnail"`
	Type           MediaType `json:"type"`
	DownloadedFrom int64     `json:"downloaded_from"` // Id of the thread download that created this media this media
}

type MediaIdentity struct {
	Media        *Media `json:"media"`
	CategoryUUID string `json:"category_uuid"`
	CategoryPath string `json:"category_path"`
	ClusterUUID  string `json:"cluster_uuid"`
	ClusterPath  string `json:"cluster_path"`
}

func (media_identity MediaIdentity) ToWeakIdentity() *MediaWeakIdentity {
	return &MediaWeakIdentity{
		MediaUUID:    media_identity.Media.Uuid,
		MediaName:    media_identity.Media.Name,
		CategoryUUID: media_identity.CategoryUUID,
		CategoryPath: media_identity.CategoryPath,
	}
}

type MediaWeakIdentity struct {
	MediaUUID    string `json:"media_uuid"`
	MediaName    string `json:"media_name"`
	CategoryUUID string `json:"category_uuid"`
	CategoryPath string `json:"category_path"`
}

func (weak_identity MediaWeakIdentity) String() string {
	return fmt.Sprintf("{ MediaUUID: '%s', MediaName: '%s', CategoryUUID: '%s', CategoryPath: '%s' }", weak_identity.MediaUUID, weak_identity.MediaName, weak_identity.CategoryUUID, weak_identity.CategoryPath)
}

func (weak_identity MediaWeakIdentity) ToClusterWeakIdentity() *CategoryClusterWeakIdentity {
	return &CategoryClusterWeakIdentity{
		ClusterUUID:   weak_identity.CategoryUUID,
		ClusterFsPath: weak_identity.CategoryPath,
	}
}

type MediaType string

const (
	Image MediaType = "IMAGE"
	Video MediaType = "VIDEO"
)

func CreateNewMedia(name string, main_category string, is_video bool, downloaded_from int64) *Media {
	var new_media *Media = new(Media)

	new_media.Name = name
	new_media.MainCategory = main_category
	new_media.Type = Image

	if is_video {
		new_media.Type = Video
	}

	new_media.DownloadedFrom = downloaded_from
	new_media.LastSeen = time.Now()

	new_media.Uuid = helpers.GenerateSha1ID(fmt.Sprintf("%s-%s-%d", new_media.Name, new_media.MainCategory, new_media.LastSeen.Unix()))

	return new_media
}

func (media_identity MediaIdentity) AbsFilename() string {
	return filepath.Join(media_identity.ClusterPath, media_identity.CategoryPath, media_identity.Media.Name)
}

func CreateNewMediaIdentity(media *Media, category *Category, cluster *CategoryCluster) *MediaIdentity {
	media_identity := new(MediaIdentity)

	media_identity.Media = media
	media_identity.CategoryUUID = category.Uuid
	media_identity.CategoryPath = category.Fullpath
	media_identity.ClusterUUID = cluster.Uuid
	media_identity.ClusterPath = cluster.FsPath

	return media_identity
}
