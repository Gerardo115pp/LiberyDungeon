package dungeon_models

import (
	"cmp"
	"fmt"
	"libery-dungeon-libs/helpers"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type CategoryIdentity struct {
	Category    *Category `json:"category"`
	ClusterUUID string    `json:"cluster_uuid"`
	ClusterPath string    `json:"cluster_path"`
}

func (ci *CategoryIdentity) ToWeakIdentity() *CategoryWeakIdentity {
	return &CategoryWeakIdentity{
		CategoryUUID: ci.Category.Uuid,
		CategoryPath: ci.Category.Fullpath,
		ClusterUUID:  ci.ClusterUUID,
		ClusterPath:  ci.ClusterPath,
	}
}

func (ci *CategoryIdentity) ToClusterWeakIdentity() *CategoryClusterWeakIdentity {
	return &CategoryClusterWeakIdentity{
		ClusterUUID:   ci.ClusterUUID,
		ClusterFsPath: ci.ClusterPath,
	}
}

type CategoryWeakIdentity struct {
	CategoryUUID string `json:"category_uuid"`
	CategoryPath string `json:"category_path"`
	ClusterUUID  string `json:"cluster_uuid"`
	ClusterPath  string `json:"cluster_path"`
}

func (weak_category_identity CategoryWeakIdentity) ToMediaIdentity(media *Media) *MediaIdentity {
	return &MediaIdentity{
		Media:        media,
		CategoryUUID: weak_category_identity.CategoryUUID,
		CategoryPath: weak_category_identity.CategoryPath,
		ClusterUUID:  weak_category_identity.ClusterUUID,
		ClusterPath:  weak_category_identity.ClusterPath,
	}
}

func CreateNewCategoryIdentity(category *Category, cluster *CategoryCluster) *CategoryIdentity {
	category_identity := new(CategoryIdentity)

	category_identity.Category = category
	category_identity.ClusterUUID = cluster.Uuid
	category_identity.ClusterPath = cluster.FsPath

	return category_identity
}

type Category struct {
	Uuid              string `json:"uuid"`
	Name              string `json:"name"`
	Fullpath          string `json:"fullpath"`
	Parent            string `json:"parent"` // Id of another Category, if empty, it's a main Category
	Cluster           string `json:"cluster"`
	CategoryThumbnail string `json:"category_thumbnail"`
}

func (c *Category) CopyContent(other Category) {
	c.Uuid = other.Uuid
	c.Name = other.Name
	c.Fullpath = other.Fullpath
	c.Parent = other.Parent
	c.Cluster = other.Cluster
	c.CategoryThumbnail = other.CategoryThumbnail
}

func (c *Category) RecalculateHash() string {
	recalculation_timestamp := time.Now().Unix()

	return helpers.GenerateSha1ID(fmt.Sprintf("%s-%s-%d", c.Name, c.Parent, recalculation_timestamp))
}

type ChildCategory struct {
	Name              string `json:"name"`
	Uuid              string `json:"uuid"`
	Fullpath          string `json:"fullpath"`
	CategoryThumbnail string `json:"category_thumbnail"`
}

type CategoryLeaf struct {
	InnerCategories []ChildCategory `json:"inner_categories"`
	Content         []Media         `json:"content"`
	Category
}

// Check if all the content inside the leaf is a single series of medias(e.g a tv show). it does this by
// checking if the content has the same name with a single numeric variant at roughly the same index. e.g:
// episode-1080p-001.webm, episode-1080p-002.webm, episode-1080p-003.webm.
// If the content is a series, it will sort the content by the numeric variant.
func (cl *CategoryLeaf) SortContentSeries() bool {
	var common_shared_prefix string = cl.getCommonPrefix()

	if common_shared_prefix == "" {
		return false
	}

	var sorted_content []Media = make([]Media, len(cl.Content))

	copy(sorted_content, cl.Content)

	slices.SortFunc(sorted_content, func(a, b Media) int {
		series_number_a := getMediaSeriesNumber(a.Name, common_shared_prefix)
		series_number_b := getMediaSeriesNumber(b.Name, common_shared_prefix)

		return cmp.Compare(series_number_a, series_number_b)
	})

	cl.Content = sorted_content

	return true
}

func (cl CategoryLeaf) getCommonPrefix() string {
	if len(cl.Content) < 2 {
		return ""
	}

	var common_shared_prefix string = ""
	var last_seen_number_index int = -1

	first_media_name := cl.Content[0].Name
	second_media_name := cl.Content[1].Name

	for k := 0; k < len(first_media_name) && k < len(second_media_name); k++ {

		if first_media_name[k] != second_media_name[k] {
			if !(first_media_name[k] >= '0' && first_media_name[k] <= '9') && !(second_media_name[k] >= '0' && second_media_name[k] <= '9') {
				return ""
			}
			break
		}

		if first_media_name[k] >= '0' && first_media_name[k] <= '9' {
			if last_seen_number_index == -1 {

				last_seen_number_index = k
			}
		} else {
			last_seen_number_index = -1
		}

		common_shared_prefix += string(first_media_name[k])
	}

	if last_seen_number_index != -1 {
		common_shared_prefix = common_shared_prefix[:last_seen_number_index]
	}

	for h := 0; h < len(cl.Content); h++ {
		if !strings.HasPrefix(cl.Content[h].Name, common_shared_prefix) {
			return ""
		}
	}

	return common_shared_prefix
}

func getMediaSeriesNumber(media_name string, common_shared_prefix string) int {
	var numeric_part string = ""
	var ord int

	for h := len(common_shared_prefix); h < len(media_name); h++ {
		if media_name[h] < '0' || media_name[h] > '9' {
			break
		}

		numeric_part += string(media_name[h])
	}

	ord, err := strconv.Atoi(numeric_part)
	if err != nil {
		echo.EchoErr(err)
		return -1
	}

	return ord
}

func CreateNewCategory(name string, parent string, parent_path string, cluster_id string) Category {
	var new_category Category

	category_path := path.Join(parent_path, name)

	if !strings.HasSuffix(category_path, "/") {
		category_path += "/"
	}

	new_category.Parent = parent
	new_category.Name = name
	new_category.Fullpath = category_path
	new_category.Cluster = cluster_id

	creation_timestamp := time.Now().Unix()

	new_category.Uuid = helpers.GenerateSha1ID(fmt.Sprintf("%s-%s-%d", new_category.Name, new_category.Parent, creation_timestamp))

	return new_category
}
