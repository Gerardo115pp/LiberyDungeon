package sqlite_users

import (
	"context"
	"database/sql"
	"fmt"
	"libery-dungeon-libs/dungeonsec"
	service_models "libery_users_service/models"
	"slices"
	"sync"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type UsersDB struct {
	db_conn    *sql.DB
	update_mux sync.Mutex
}

func NewUsersDB() (*UsersDB, error) {
	var users_db *UsersDB = new(UsersDB)

	db, err := openDB()
	if err != nil {
		return nil, err
	}

	users_db.db_conn = db

	db.Exec("PRAGMA foreign_keys = ON")

	return users_db, nil
}

func (users_db *UsersDB) AddUserCTX(ctx context.Context, user *service_models.User) error {
	stmt, err := users_db.db_conn.PrepareContext(ctx, "INSERT INTO `users`(`uuid`, `username`, `password`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.UUID, user.Username, user.SecretHash)
	if err != nil {
		return err
	}

	return nil
}

func (users_db *UsersDB) AddUser(user *service_models.User) error {
	return users_db.AddUserCTX(context.Background(), user)
}

func (users_db *UsersDB) AddUserToRoleCTX(ctx context.Context, user *service_models.User, role string) error {
	stmt, err := users_db.db_conn.PrepareContext(ctx, "INSERT INTO `user_roles`(`user`, `role`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.UUID, role)
	if err != nil {
		return err
	}

	return nil
}

func (users_db *UsersDB) AddUserToRole(user *service_models.User, role string) error {
	return users_db.AddUserToRoleCTX(context.Background(), user, role)
}

func (users_db *UsersDB) AddGrantCTX(ctx context.Context, grant string) error {
	add_grant_stmt, err := users_db.db_conn.PrepareContext(ctx, "INSERT INTO `grants`(`grant_label`) VALUES (?)")
	if err != nil {
		return err
	}
	defer add_grant_stmt.Close()

	_, err = add_grant_stmt.ExecContext(ctx, grant)

	return err
}

func (users_db *UsersDB) AddGrant(grant string) error {
	return users_db.AddGrantCTX(context.Background(), grant)
}

func (users_db *UsersDB) AddGrantToRoleCTX(ctx context.Context, role_label string, grant string, propagate_grant_up bool) error {
	add_grant_to_role_stmt, err := users_db.db_conn.PrepareContext(ctx, "INSERT INTO `role_grants`(`role`, `grant`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer add_grant_to_role_stmt.Close()

	_, err = add_grant_to_role_stmt.ExecContext(ctx, role_label, grant)
	if err != nil {
		return err
	}

	if propagate_grant_up {
		var role service_models.RoleTaxonomy

		role, err = users_db.GetRoleCTX(ctx, role_label)
		if err != nil {
			return err
		}

		go users_db.updateHigherRoles(role.RoleHierarchy, role.RoleGrants)
	}

	return err
}

func (users_db *UsersDB) AddGrantToRole(role_label string, grant string, propagate_grant_up bool) error {
	return users_db.AddGrantToRoleCTX(context.Background(), role_label, grant, propagate_grant_up)
}

func (users_db *UsersDB) allGrantsExistCTX(ctx context.Context, grants []string) (bool, error) {
	all_grants, err := users_db.GetAllGrantsCTX(ctx)
	if err != nil {
		return false, err
	}

	var all_exist bool = true

	for _, grant := range grants {
		all_exist = slices.Contains(all_grants, grant)
		echo.EchoDebug(fmt.Sprintf("Checking if grant<%s> exists: %t", grant, all_exist))

		if !all_exist {
			break
		}
	}

	return all_exist, nil
}

func (users_db *UsersDB) CreateNewRoleCTX(ctx context.Context, role service_models.RoleTaxonomy) error {
	add_role_stmt, err := users_db.db_conn.PrepareContext(ctx, "INSERT INTO `roles`(`role_label`, `role_hierarchy`) VALUES (?, ?)")
	if err != nil {
		echo.Echo(echo.OrangeFG, fmt.Sprintf("In Users.CreateNewRoleCTX while preparing statement: %s", err.Error()))
		return err
	}
	defer add_role_stmt.Close()

	if role.RoleHierarchy < 0 {
		return fmt.Errorf("Role hierarchy must be a positive integer")
	}

	if all_grants_exist, err := users_db.allGrantsExistCTX(ctx, role.RoleGrants); err != nil || !all_grants_exist {
		return fmt.Errorf("One or more grants do not exist")
	}

	_, err = add_role_stmt.ExecContext(ctx, role.RoleLabel, role.RoleHierarchy)
	if err != nil {
		echo.Echo(echo.OrangeFG, fmt.Sprintf("In Users.CreateNewRoleCTX while adding role: %s", err.Error()))
		return err
	}

	err = users_db.populateNewRoleGrantsCTX(ctx, &role)
	if err != nil {
		echo.Echo(echo.OrangeFG, fmt.Sprintf("In Users.CreateNewRoleCTX while populating new role grants: %s", err.Error()))
		return err
	}

	for _, grant := range role.RoleGrants {
		err = users_db.AddGrantToRoleCTX(ctx, role.RoleLabel, grant, false)
		if err != nil {
			echo.Echo(echo.OrangeFG, fmt.Sprintf("In Users.CreateNewRoleCTX while adding grant to role: %s", err.Error()))
			return err
		}
	}

	go users_db.updateHigherRoles(role.RoleHierarchy, role.RoleGrants)

	return nil
}

func (users_db *UsersDB) CreateNewRole(role service_models.RoleTaxonomy) error {
	return users_db.CreateNewRoleCTX(context.Background(), role)
}

func (users_db *UsersDB) DeleteUserByUuidCTX(ctx context.Context, user_uuid string) error {
	delete_user_stmt, err := users_db.db_conn.PrepareContext(ctx, "DELETE FROM `users` WHERE `uuid` = ?")
	if err != nil {
		return err
	}
	defer delete_user_stmt.Close()

	_, err = delete_user_stmt.ExecContext(ctx, user_uuid)

	return err
}

func (users_db *UsersDB) DeleteUserByUuid(user_uuid string) error {
	return users_db.DeleteUserByUuidCTX(context.Background(), user_uuid)
}

func (users_db *UsersDB) DeleteGrantCTX(ctx context.Context, grant string) error {
	delete_grant_stmt, err := users_db.db_conn.PrepareContext(ctx, "DELETE FROM `grants` WHERE `grant_label` = ?")
	if err != nil {
		return err
	}
	defer delete_grant_stmt.Close()

	_, err = delete_grant_stmt.ExecContext(ctx, grant)

	return err
}

func (users_db *UsersDB) DeleteGrant(grant string) error {
	return users_db.DeleteGrantCTX(context.Background(), grant)
}

func (users_db *UsersDB) DeleteGrantFromRoleCTX(ctx context.Context, role_label string, grant string) error {
	delete_grant_stmt, err := users_db.db_conn.PrepareContext(ctx, "DELETE FROM `role_grants` WHERE `role` = ? AND `grant` = ?")
	if err != nil {
		return err
	}
	defer delete_grant_stmt.Close()

	_, err = delete_grant_stmt.ExecContext(ctx, role_label, grant)

	return err
}

func (users_db *UsersDB) DeleteGrantFromRole(role_label string, grant string) error {
	return users_db.DeleteGrantFromRoleCTX(context.Background(), role_label, grant)
}

func (users_db *UsersDB) DeleteUserFromRoleCTX(ctx context.Context, user *service_models.User, role string) error {
	delete_user_role_stmt, err := users_db.db_conn.PrepareContext(ctx, "DELETE FROM `user_roles` WHERE `user` = ? AND `role` = ?")
	if err != nil {
		return err
	}
	defer delete_user_role_stmt.Close()

	_, err = delete_user_role_stmt.ExecContext(ctx, user.UUID, role)

	return err
}

func (users_db *UsersDB) DeleteUserFromRole(user *service_models.User, role string) error {
	return users_db.DeleteUserFromRoleCTX(context.Background(), user, role)
}

func (users_db *UsersDB) DeleteRoleCTX(ctx context.Context, role_label string) error {
	delete_role_stmt, err := users_db.db_conn.PrepareContext(ctx, "DELETE FROM `roles` WHERE `role_label` = ?")
	if err != nil {
		return err
	}
	defer delete_role_stmt.Close()

	_, err = delete_role_stmt.ExecContext(ctx, role_label)

	return err
}

func (users_db *UsersDB) DeleteRole(role_label string) error {
	return users_db.DeleteRoleCTX(context.Background(), role_label)
}

// Returns all roles in the hierarchy directly below the given role_hierarchy. e.g if 3 is passed and in the db there are roles with the hierarchies [0, 1, 7, 10] then
// all the roles with a hierarchy 7 will be returned(remember, the lower the hierarchy number, the higher hierarchy the role is)
func (users_db *UsersDB) FindRolesBelowHierarchyCTX(ctx context.Context, role_hierarchy int) ([]service_models.RoleTaxonomy, error) {
	var roles []service_models.RoleTaxonomy = make([]service_models.RoleTaxonomy, 0)

	get_labels_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `role_label` FROM `roles` WHERE `role_hierarchy` = (SELECT `role_hierarchy` FROM `roles` WHERE `role_hierarchy` > ? ORDER BY `role_hierarchy` LIMIT 1)")
	if err != nil {
		return roles, err
	}
	defer get_labels_stmt.Close()

	rows, err := get_labels_stmt.QueryContext(ctx, role_hierarchy)
	if err != nil {
		return roles, err
	}
	defer rows.Close()

	var role_label string
	var role service_models.RoleTaxonomy

	for rows.Next() {
		err = rows.Scan(&role_label)
		if err != nil {
			return roles, err
		}

		role, err = users_db.GetRoleCTX(ctx, role_label)
		if err != nil {
			return roles, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (users_db *UsersDB) FindRolesBelowHierarchy(role_hierarchy int) ([]service_models.RoleTaxonomy, error) {
	return users_db.FindRolesBelowHierarchyCTX(context.Background(), role_hierarchy)
}

func (users_db *UsersDB) findRolesAboveHierarchyCTX(ctx context.Context, role_hierarchy int) ([]service_models.RoleTaxonomy, error) {
	var roles []service_models.RoleTaxonomy = make([]service_models.RoleTaxonomy, 0)

	get_labels_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `role_label` FROM `roles` WHERE `role_hierarchy` < ?")
	if err != nil {
		return roles, err
	}
	defer get_labels_stmt.Close()

	rows, err := get_labels_stmt.QueryContext(ctx, role_hierarchy)
	if err != nil {
		return roles, err
	}
	defer rows.Close()

	var role_label string
	var role service_models.RoleTaxonomy

	for rows.Next() {
		err = rows.Scan(&role_label)
		if err != nil {
			return roles, err
		}

		role, err = users_db.GetRoleCTX(ctx, role_label)
		if err != nil {
			return roles, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (users_db *UsersDB) GetUserByUsername(username string) (user *service_models.User, err error) {
	return users_db.GetUserByUsernameCTX(context.Background(), username)
}

func (users_db *UsersDB) GetUserByUsernameCTX(ctx context.Context, username string) (user *service_models.User, err error) {
	var user_uuid string
	var user_secret_hash string

	stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `uuid`, `username`, `password` FROM `users` WHERE `username` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, username).Scan(&user_uuid, &username, &user_secret_hash)
	if err != nil {
		return nil, err
	}

	var found_user *service_models.User = new(service_models.User)

	found_user.UUID = user_uuid
	found_user.Username = username
	found_user.SecretHash = user_secret_hash

	return found_user, nil
}

func (users_db *UsersDB) GetUserByUuid(user_uuid string) (user *service_models.User, err error) {
	return users_db.GetUserByUuidCTX(context.Background(), user_uuid)
}

func (users_db *UsersDB) GetUserByUuidCTX(ctx context.Context, user_uuid string) (user *service_models.User, err error) {
	user = new(service_models.User)

	stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `uuid`, `username`, `password` FROM `users` WHERE `uuid` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user_uuid).Scan(&user.UUID, &user.Username, &user.SecretHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (users_db *UsersDB) GetAllUsersCTX(ctx context.Context) ([]*service_models.User, error) {
	var users []*service_models.User = make([]*service_models.User, 0)

	rows, err := users_db.db_conn.QueryContext(ctx, "SELECT `uuid`, `username`, `password` FROM `users`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user *service_models.User = new(service_models.User)

		err = rows.Scan(&user.UUID, &user.Username, &user.SecretHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (users_db *UsersDB) GetAllUsers() ([]*service_models.User, error) {
	return users_db.GetAllUsersCTX(context.Background())
}

func (users_db *UsersDB) GetRoleGrantsCTX(ctx context.Context, role_label string) ([]string, error) {
	var grants []string = make([]string, 0)

	role_grants_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `grant` FROM `role_grants` WHERE `role` = ?")
	if err != nil {
		return nil, err
	}
	defer role_grants_stmt.Close()

	rows, err := role_grants_stmt.QueryContext(ctx, role_label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var grant string

		err = rows.Scan(&grant)
		if err != nil {
			return nil, err
		}

		grants = append(grants, grant)
	}

	return grants, nil
}

func (users_db *UsersDB) GetRoleGrants(role_label string) ([]string, error) {
	return users_db.GetRoleGrantsCTX(context.Background(), role_label)
}

func (users_db *UsersDB) GetRoleCTX(ctx context.Context, role_label string) (service_models.RoleTaxonomy, error) {
	var role service_models.RoleTaxonomy

	role_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `role_hierarchy` FROM `roles` WHERE `role_label` = ?")
	if err != nil {
		return role, err
	}
	defer role_stmt.Close()

	err = role_stmt.QueryRowContext(ctx, role_label).Scan(&role.RoleHierarchy)
	if err != nil {
		return role, err
	}

	role.RoleLabel = role_label

	var role_grants []string
	role_grants, err = users_db.GetRoleGrantsCTX(ctx, role_label)
	if err != nil {
		return role, err
	}

	role.RoleGrants = role_grants

	return role, nil
}

func (users_db *UsersDB) GetRole(role_label string) (service_models.RoleTaxonomy, error) {
	return users_db.GetRoleCTX(context.Background(), role_label)
}

func (users_db *UsersDB) GetAllGrantsCTX(ctx context.Context) ([]string, error) {
	var grants []string = make([]string, 0)

	rows, err := users_db.db_conn.QueryContext(ctx, "SELECT `grant_label` FROM `grants`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var grant string

		err = rows.Scan(&grant)
		if err != nil {
			return nil, err
		}

		grants = append(grants, grant)
	}

	return grants, nil
}

func (users_db *UsersDB) GetAllGrants() ([]string, error) {
	return users_db.GetAllGrantsCTX(context.Background())
}

func (users_db *UsersDB) GetUserRolesCTX(ctx context.Context, user *service_models.User) ([]service_models.RoleTaxonomy, error) {
	var roles []service_models.RoleTaxonomy = make([]service_models.RoleTaxonomy, 0)

	user_roles_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `role` FROM `user_roles` WHERE `user` = ?")
	if err != nil {
		return nil, err
	}
	defer user_roles_stmt.Close()

	role_data_stmt, err := users_db.db_conn.PrepareContext(ctx, "SELECT `role_hierarchy` FROM `roles` WHERE `role_label` = ?")
	if err != nil {
		return nil, err
	}
	defer role_data_stmt.Close()

	rows, err := user_roles_stmt.QueryContext(ctx, user.UUID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var role_label string

		err = rows.Scan(&role_label)
		if err != nil {
			return nil, err
		}

		var role_hierarchy int

		err = role_data_stmt.QueryRowContext(ctx, role_label).Scan(&role_hierarchy)
		if err != nil {
			return nil, err
		}

		var role_grants []string
		role_grants, err = users_db.GetRoleGrantsCTX(ctx, role_label)
		if err != nil {
			return nil, err
		}

		var role_taxonomy service_models.RoleTaxonomy = service_models.RoleTaxonomy{
			RoleLabel:     role_label,
			RoleHierarchy: role_hierarchy,
			RoleGrants:    role_grants,
		}

		roles = append(roles, role_taxonomy)
	}

	return roles, nil
}

func (users_db *UsersDB) GetUserRoles(user *service_models.User) ([]service_models.RoleTaxonomy, error) {
	return users_db.GetUserRolesCTX(context.Background(), user)
}

func (users_db *UsersDB) GetAllRolesCTX(ctx context.Context) ([]service_models.RoleTaxonomy, error) {
	var roles []service_models.RoleTaxonomy = make([]service_models.RoleTaxonomy, 0)

	rows, err := users_db.db_conn.QueryContext(ctx, "SELECT `role_label`, `role_hierarchy` FROM `roles`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role_label string
		var role_hierarchy int

		err = rows.Scan(&role_label, &role_hierarchy)
		if err != nil {
			return nil, err
		}

		var role_grants []string
		role_grants, err = users_db.GetRoleGrantsCTX(ctx, role_label)
		if err != nil {
			return nil, err
		}

		var role_taxonomy service_models.RoleTaxonomy = service_models.RoleTaxonomy{
			RoleLabel:     role_label,
			RoleHierarchy: role_hierarchy,
			RoleGrants:    role_grants,
		}

		roles = append(roles, role_taxonomy)
	}

	return roles, nil
}

func (users_db *UsersDB) GetAllRoles() ([]service_models.RoleTaxonomy, error) {
	return users_db.GetAllRolesCTX(context.Background())
}

func (users_db *UsersDB) GetUserGrantsCTX(ctx context.Context, user *service_models.User) ([]string, error) {
	var user_roles []service_models.RoleTaxonomy

	user_roles, err := users_db.GetUserRolesCTX(ctx, user)
	if err != nil {
		return nil, err
	}

	var grants []string = make([]string, 0)

	for _, role := range user_roles {
		grants = append(grants, role.RoleGrants...)
	}

	return grants, nil
}

func (users_db *UsersDB) GetUserGrants(user *service_models.User) ([]string, error) {
	return users_db.GetUserGrantsCTX(context.Background(), user)
}

// Call it after creating a new role to make sure higher roles get the possibly new grants.
func (users_db *UsersDB) updateHigherRoles(from_hierarchy int, new_grants []string) error {
	users_db.update_mux.Lock()
	defer users_db.update_mux.Unlock()
	echo.Echo(echo.CyanFG, "Updating higher roles")

	var higher_roles []service_models.RoleTaxonomy

	higher_roles, err := users_db.findRolesAboveHierarchyCTX(context.Background(), from_hierarchy)
	if err != nil {
		return err
	}

	var filtered_grants []string = make([]string, 0)

	// Filter out grants that are not inheritable.(ALL_PRIVILEGES and GRANT_OPTION)
	for _, grant := range new_grants {
		if grant != dungeonsec.PlatformGrant_ALL_PRIVILEGES && grant != dungeonsec.PlatformGrant_GrantPrivileges {
			filtered_grants = append(filtered_grants, grant)
		}
	}

	for _, higher_role := range higher_roles {
		for _, new_grant := range filtered_grants {
			if !higher_role.HasGrant(new_grant) {
				err = users_db.AddGrantToRole(higher_role.RoleLabel, new_grant, false)
				if err != nil {
					return err
				}
			}
		}
	}

	echo.Echo(echo.CyanFG, fmt.Sprintf("Updated %d higher roles", len(higher_roles)))

	return nil
}

// Receives a user and updates an entire on the database matching the user's UUID. this can change a user's username and password.
func (users_db *UsersDB) UpdateUserCTX(ctx context.Context, user *service_models.User) error {
	update_user_stmt, err := users_db.db_conn.PrepareContext(ctx, "UPDATE `users` SET `username` = ?, `password` = ? WHERE `uuid` = ?")
	if err != nil {
		return err
	}
	defer update_user_stmt.Close()

	_, err = update_user_stmt.ExecContext(ctx, user.Username, user.SecretHash, user.UUID)

	return err
}

func (users_db *UsersDB) UpdateUser(user *service_models.User) error {
	return users_db.UpdateUserCTX(context.Background(), user)
}

// Fetches all the roles below this role's hierarchy and adds all of their grants to new role.
func (users_db *UsersDB) populateNewRoleGrantsCTX(ctx context.Context, role *service_models.RoleTaxonomy) error {
	var below_roles []service_models.RoleTaxonomy

	below_roles, err := users_db.FindRolesBelowHierarchyCTX(ctx, role.RoleHierarchy)
	if err != nil {
		return err
	}

	var inherited_grants []string = make([]string, 0)

	for _, below_role := range below_roles {
		inherited_grants = append(inherited_grants, below_role.RoleGrants...)
	}

	role.AddUniqueGrants(inherited_grants)

	return nil
}

func (users_db *UsersDB) Close() error {
	return users_db.db_conn.Close()
}
