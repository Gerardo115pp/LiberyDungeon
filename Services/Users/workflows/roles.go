package workflows

import (
	service_models "libery_users_service/models"
)

// Returns the highest hierarchy present in the roles slice. the smaller the number, the higher the hierarchy.
func GetHighestRoleHierarchy(roles []service_models.RoleTaxonomy) int {
	var hierarchy_set bool = false
	var highest_hierarchy int

	for _, role := range roles {
		if !hierarchy_set {
			highest_hierarchy = role.RoleHierarchy
			hierarchy_set = true
			continue
		}

		if role.RoleHierarchy < highest_hierarchy {
			highest_hierarchy = role.RoleHierarchy
		}
	}

	return highest_hierarchy
}

// Returns a slice of grants compiled from the provided roles, the items in the slice are guaranteed to be unique
// even if the exist in more the one of the roles.
func CompileUserGrants(roles []service_models.RoleTaxonomy) []string {
	grants_map := make(map[string]struct{})
	var grants []string = make([]string, 0)

	for _, role := range roles {
		for _, grant := range role.RoleGrants {
			if _, ok := grants_map[grant]; !ok {
				grants_map[grant] = struct{}{}
				grants = append(grants, grant)
			}
		}
	}

	return grants
}
