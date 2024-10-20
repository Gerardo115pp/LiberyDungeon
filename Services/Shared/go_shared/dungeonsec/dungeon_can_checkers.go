package dungeonsec

var canGrant UserCanChecker = factory_UserCanChecker(PlatformGrant_GrantPrivileges, false)
var canReadUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ReadUsers, true)
var canModifyUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ModifyUsers, true)
var canDeleteUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_DeleteUsers, true)
var canCreateUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ModifyUsers, true)
var canViewPrivateClusters UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_ReadPrivate, true)
var canViewContent UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_Read, true)
var canAlterPrivateClusters UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_AlterPrivate, true)
var canUploadFiles UserCanChecker = factory_UserCanChecker(PlatformGrant_UploadFiles, true)
var canContentAlter UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_Alter, true)
var canDungeonTagsCreate UserCanChecker = factory_UserCanChecker(PlatformGrant_DungeonTags_Create, true)
var canDungeonTagsTag UserCanChecker = factory_UserCanChecker(PlatformGrant_DungeonTags_Tag, true)
var canDungeonTagsUntag UserCanChecker = factory_UserCanChecker(PlatformGrant_DungeonTags_Untag, true)
var canDungeonTagsTaxonomyCreate UserCanChecker = factory_UserCanChecker(PlatformGrant_DungeonTags_TaxonomyCreate, true)

func CanGrant(grants []string) bool {
	return canGrant(grants)
}

func CanReadUsers(grants []string) bool {
	return canReadUsers(grants)
}

func CanModifyUsers(grants []string) bool {
	return canModifyUsers(grants)
}

func CanDeleteUsers(grants []string) bool {
	return canDeleteUsers(grants)
}

func CanCreateUsers(grants []string) bool {
	return canCreateUsers(grants)
}

func CanViewPrivateClusters(grants []string) bool {
	return canViewPrivateClusters(grants)
}

func CanViewContent(grants []string) bool {
	return canViewContent(grants)
}

func CanAlterPrivateClusters(grants []string) bool {
	return canAlterPrivateClusters(grants)
}

func CanUploadFiles(grants []string) bool {
	return canUploadFiles(grants)
}

func CanContentAlter(grants []string) bool {
	return canContentAlter(grants)
}

func CanDungeonTagsCreate(grants []string) bool {
	return canDungeonTagsCreate(grants)
}

func CanDungeonTagsTag(grants []string) bool {
	return canDungeonTagsTag(grants)
}

func CanDungeonTagsUntag(grants []string) bool {
	return canDungeonTagsUntag(grants)
}

func CanDungeonTagsTaxonomyCreate(grants []string) bool {
	return canDungeonTagsTaxonomyCreate(grants)
}
