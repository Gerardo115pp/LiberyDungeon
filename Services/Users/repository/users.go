package repository

import (
	"context"
	service_models "libery_users_service/models"
)

type UsersRepository interface {
	AddUser(user *service_models.User) error
	AddUserCTX(ctx context.Context, user *service_models.User) error
	AddUserToRole(user *service_models.User, role string) error
	AddUserToRoleCTX(ctx context.Context, user *service_models.User, role string) error
	AddGrant(grant string) error
	AddGrantCTX(ctx context.Context, grant string) error
	AddGrantToRole(role_label string, grant string, propagate_grant_up bool) error
	AddGrantToRoleCTX(ctx context.Context, role_label string, grant string, propagate_grant_up bool) error
	CreateNewRole(role service_models.RoleTaxonomy) error
	CreateNewRoleCTX(ctx context.Context, role service_models.RoleTaxonomy) error
	DeleteUserByUuid(user_uuid string) error
	DeleteUserByUuidCTX(ctx context.Context, user_uuid string) error
	DeleteGrant(grant string) error
	DeleteGrantCTX(ctx context.Context, grant string) error
	DeleteUserFromRole(user *service_models.User, role string) error
	DeleteUserFromRoleCTX(ctx context.Context, user *service_models.User, role string) error
	DeleteRole(role_label string) error
	DeleteRoleCTX(ctx context.Context, role_label string) error
	DeleteGrantFromRole(role_label string, grant string) error
	DeleteGrantFromRoleCTX(ctx context.Context, role_label string, grant string) error
	FindRolesBelowHierarchy(role_hierarchy int) ([]service_models.RoleTaxonomy, error)
	FindRolesBelowHierarchyCTX(ctx context.Context, role_hierarchy int) ([]service_models.RoleTaxonomy, error)
	GetUserByUsername(username string) (user *service_models.User, err error)
	GetUserByUsernameCTX(ctx context.Context, username string) (user *service_models.User, err error)
	GetUserByUuid(user_uuid string) (user *service_models.User, err error)
	GetUserByUuidCTX(ctx context.Context, user_uuid string) (user *service_models.User, err error)
	GetAllUsers() ([]*service_models.User, error)
	GetAllUsersCTX(ctx context.Context) ([]*service_models.User, error)
	GetUserGrants(user *service_models.User) ([]string, error)
	GetUserGrantsCTX(ctx context.Context, user *service_models.User) ([]string, error)
	GetUserRoles(user *service_models.User) ([]service_models.RoleTaxonomy, error)
	GetUserRolesCTX(ctx context.Context, user *service_models.User) ([]service_models.RoleTaxonomy, error)
	GetAllRoles() ([]service_models.RoleTaxonomy, error)
	GetAllRolesCTX(ctx context.Context) ([]service_models.RoleTaxonomy, error)
	GetRoleGrants(role_label string) ([]string, error)
	GetRoleGrantsCTX(ctx context.Context, role_label string) ([]string, error)
	GetRole(role_label string) (service_models.RoleTaxonomy, error)
	GetRoleCTX(ctx context.Context, role_label string) (service_models.RoleTaxonomy, error)
	GetAllGrants() ([]string, error)
	GetAllGrantsCTX(ctx context.Context) ([]string, error)
	UpdateUser(user *service_models.User) error
	UpdateUserCTX(ctx context.Context, user *service_models.User) error
}

var UsersRepo UsersRepository

func SetUsersRepository(repo UsersRepository) {
	UsersRepo = repo
}
