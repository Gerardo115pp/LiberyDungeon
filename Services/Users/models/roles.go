package models

import (
	"libery-dungeon-libs/dungeonsec"
	"slices"
)

// Default user roles
const SUPER_ADMIN_ROLE = "super_admin"
const ADMIN_ROLE = "admin"
const VISITOR_ROLE = "visitor"

type UserWithRoles struct {
	User
	Roles []string `json:"roles"`
}

type RoleTaxonomy struct {
	RoleLabel     string   `json:"role_label"`
	RoleHierarchy int      `json:"role_hierarchy"` // Lower number means higher hierarchy. max hierarchy is 0 which is super admin.
	RoleGrants    []string `json:"role_grants"`
}

func (rt *RoleTaxonomy) AddUniqueGrants(new_grants []string) {
	grants_map := rt.getGrantsMap()

	for _, new_grant := range new_grants {
		if _, exists := grants_map[new_grant]; !exists {
			rt.RoleGrants = append(rt.RoleGrants, new_grant)
		}
	}
}

func (rt RoleTaxonomy) HasGrant(grant string) bool {
	return slices.Contains(rt.RoleGrants, grant)
}

func (rt RoleTaxonomy) getGrantsMap() map[string]struct{} {
	grants_map := make(map[string]struct{})

	for _, grant := range rt.RoleGrants {
		grants_map[grant] = struct{}{}
	}

	return grants_map
}

// Default roles grants. Roles with higher hierarchy inherit the grants of roles with lower hierarchy.

var SuperAdminGrants = []string{
	dungeonsec.PlatformGrant_GrantPrivileges,
}

var AdminGrants = []string{
	dungeonsec.PlatformGrant_ALL_PRIVILEGES,
}

var VisitorGrants = []string{
	dungeonsec.PlatformGrant_ClustersContent_Read,
}
