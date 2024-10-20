package dungeonsec

import "slices"

const (
	PlatformGrant_ALL_PRIVILEGES               string = "ALL_PRIVILEGES"
	PlatformGrant_GrantPrivileges              string = "grant_privileges"
	PlatformGrant_ModifyUsers                  string = "modify_users"
	PlatformGrant_DeleteUsers                  string = "delete_users"
	PlatformGrant_ReadUsers                    string = "read_users"
	PlatformGrant_ModifyTrashcan               string = "modify_trashcan"
	PlatformGrant_EmptyTrashcan                string = "empty_trashcan"
	PlatformGrant_ReadTrashcan                 string = "read_trashcan"
	PlatformGrant_DownloadFiles                string = "download_files"
	PlatformGrant_UploadFiles                  string = "upload_files"
	PlatformGrant_ReadPrivateFiles             string = "read_private_files"
	PlatformGrant_Clusters_Create              string = "clusters_create"
	PlatformGrant_Clusters_Sync                string = "clusters_sync"
	PlatformGrant_Clusters_Drop                string = "clusters_drop"
	PlatformGrant_ClustersContent_Alter        string = "clusters_content_alter"
	PlatformGrant_ClustersContent_Read         string = "clusters_content_read"
	PlatformGrant_ClustersContent_ReadPrivate  string = "clusters_content_read_private"
	PlatformGrant_ClustersContent_AlterPrivate string = "clusters_content_alter_private"
	PlatformGrant_DungeonTags_Create           string = "dungeon_tags_create"
	PlatformGrant_DungeonTags_Tag              string = "dungeon_tags_tag"
	PlatformGrant_DungeonTags_Untag            string = "dungeon_tags_untag"
	PlatformGrant_DungeonTags_TaxonomyCreate   string = "dungeon_tags_taxonomy_create"
)

type UserCanChecker func([]string) bool

func factory_UserCanChecker(specific_grant string, included_in_all_privileges bool) UserCanChecker {
	return func(grants []string) bool {
		var user_can bool = slices.Contains(grants, PlatformGrant_ALL_PRIVILEGES) && included_in_all_privileges

		if !user_can {
			user_can = slices.Contains(grants, specific_grant)
		}

		return user_can
	}
}
