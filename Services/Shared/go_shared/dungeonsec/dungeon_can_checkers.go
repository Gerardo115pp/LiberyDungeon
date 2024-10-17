package dungeonsec

var canGrant UserCanChecker = factory_UserCanChecker(PlatformGrant_GrantPrivileges, false)
var canReadUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ReadUsers, true)
var canModifyUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ModifyUsers, true)
var canDeleteUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_DeleteUsers, true)
var canCreateUsers UserCanChecker = factory_UserCanChecker(PlatformGrant_ModifyUsers, true)
var canViewPrivateClusters UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_ReadPrivate, true)
var canAlterPrivateClusters UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_AlterPrivate, true)
var canUploadFiles UserCanChecker = factory_UserCanChecker(PlatformGrant_UploadFiles, true)
var canContentAlter UserCanChecker = factory_UserCanChecker(PlatformGrant_ClustersContent_Alter, true)

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

func CanAlterPrivateClusters(grants []string) bool {
	return canAlterPrivateClusters(grants)
}

func CanUploadFiles(grants []string) bool {
	return canUploadFiles(grants)
}

func CanContentAlter(grants []string) bool {
	return canContentAlter(grants)
}
