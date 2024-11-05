import {
    GetIsInitialSetupRequest,
    PostCreateInitialUserRequest,
    PostCreateUserRequest,
    GetUserSignAccessRequest,
    GetUserAccessTokenValidationRequest,
    GetUserIdentityRequest,
    GetUserSignOutRequest,
    GetAllUsersRequest,
    GetAllRoleLabelsRequest,
    GetAllGrantsRequest,
    PostCreateGrantRequest,
    PostLinkGrantToRoleRequest,
    GetRoleTaxonomyRequest,
    GetRoleTaxonomiesBelowHierarchyRequest,
    PostCreateRoleRequest,
    PatchUserRolesRequest,
    DeleteUserFromRoleRequest,
    DeleteGrantFromRoleRequest,
    GetUserRolesRequest,
    DeleteUserRequest,
    DeleteGrantRequest,
    DeleteRoleRequest, 
    PutChangeUserPasswordRequest,
    PutChangeUsernameRequest
} from "@libs/DungeonsCommunication/services_requests/users_requests";


const ALL_PRIVILEGES_GRANT = "ALL_PRIVILEGES";
const GRANT_OPTION = "grant_privileges";

/**
 * Creates a function to check if a user has a given grant. and if included_in_all_privileges is true,
 * it will first check if the user has the ALL_PRIVILEGES grant, if it does, it will return 
 * true(the returned function, not this function). The UserIdentity must be binded, is expected to be
 * the function's 'this' context.
 * @this UserIdentity
 * @param {string} grant_label
 * @param {boolean} included_in_all_privileges
 * @returns {function(): boolean}
 */
function factory_UserCanChecker (grant_label, included_in_all_privileges) {

    /**
     * @memberof UserIdentity
     * @returns {boolean}
     */
    const checker = () => {

        let user_can = false;
        /** @type {UserIdentity} */ // This is needed to deal with vscode stupidity
        const user_identity = this;

        if (included_in_all_privileges) {
            user_can = user_identity.Grants.includes(ALL_PRIVILEGES_GRANT);
        }

        if (!user_can) {
            user_can = user_identity.Grants.includes(grant_label);
        }

        return user_can;
    }

    return checker;
}

/**
* @typedef {Object} UserIdentityParams
 * @property {string} uuid
 * @property {string} username
 * @property {number} role_hierarchy
 * @property {string[]} grants
*/

/**
* @typedef {Object} UserEntry
 * @property {string} uuid
 * @property {string} username
*/

/**
* @typedef {Object} RoleTaxonomyParams
 * @property {string} role_label
 * @property {number} role_hierarchy
 * @property {string[]} role_grants
*/

/**
 * A user identity helps the frontend know what the user can do and how to render different parts of the application. it
 * is not, however, a security feature, the information a UserIdentity instance should always be verified. is not used to authenticate the user in any way or event passed to the backend. If
 * you are looking for authentication related code, this is kept in a jwt signed token stored in a cookie inaccessible to the
 * frontend. 
 */
export class UserIdentity {
    /**
     * The user uuid is the unique identifier for the user. can be used to fetch metadata about the user from the backend
     * @type {string}
     */
    #the_uuid;

    /**
     * The username of the user. this is the name that is displayed to the user in the frontend
     * @type {string}
     */
    #the_username;

    /**
     * The role hierarchy of the user. the smaller the number, the higher the role. 0 been the highest privilege role.
     * @type {number}
     */
    #the_role_hierarchy;

    /**
     * The grants of the user. these are a list of capabilities that the user has. ALL_PRIVILEGES means the user can do anything except grant privileges, 
     * that also requires the 'grant_privileges' grant. grants can be any random string that any service, be it core service or plugin service, can define.
     * @type {string[]}
     */
    #the_grants;

    /**
     * @param {UserIdentityParams} param0
     */
    constructor({ uuid, username, role_hierarchy, grants }) {
        this.#the_uuid = uuid;
        this.#the_username = username;
        this.#the_role_hierarchy = role_hierarchy;
        this.#the_grants = grants;
    }

    /**
     * The user uuid is the unique identifier for the user. can be used to fetch metadata about the user from the backend
     * @type {string}
     */
    get UUID() {
        return this.#the_uuid;
    }

    /**
     * The username of the user. this is the name that is displayed to the user in the frontend
     * @type {string}
     */
    get Username() {
        return this.#the_username;
    }

    /**
     * The role hierarchy of the user. the smaller the number, the higher the role. 0 been the highest privilege role.
     * @type {number}
     */
    get RoleHierarchy() {
        return this.#the_role_hierarchy;
    }

    /**
     * The grants of the user. these are a list of capabilities that the user has. ALL_PRIVILEGES means the user can do anything except grant privileges, 
     * that also requires the 'grant_privileges' grant. grants can be any random string that any service, be it core service or plugin service, can define.
     * @type {string[]}
     */
    get Grants() {
        return this.#the_grants;
    } 

    #userCanChecker = factory_UserCanChecker.bind(this);
    
    /*=============================================
    =            Grant Checkers            =
    =============================================*/
    
    
    
        /**
         * Check if the user can modify other users. which would also mean it can read users(the opposite is not true).
         * @returns {boolean}
         */
        canModifyUsers = this.#userCanChecker("modify_users", true);

        /**
         * Checks if the user can give privileges to roles, attach roles to users or any other action that requires the 'grant_option' grant.
         * @returns {boolean}
         */
        canGrant = this.#userCanChecker("grant_privileges", false);

        /**
         * Checks whether a user can view the content of the trashcan.
         * @returns {boolean}
         */
        canViewTrashcan = this.#userCanChecker("read_trashcan", true);

        /**
         * Checks whether a user can modify the content of the trashcan. Meaning if it can restore elements from the trashcan.
         * This does not mean the user can delete elements, for that it would need the empty_trashcan grant.
         * @returns {boolean}
         */
        canModifyTrashcan = this.#userCanChecker("modify_trashcan", true);

        /**
         * Checks if the user can delete elements from the trashcan. This is the highest level of trashcan access.
         * So is safe to assume that if the user can empty the trashcan, it can also view and restore elements from it.
         * @returns {boolean}
         */
        canEmptyTrashcan = this.#userCanChecker("empty_trashcan", true);

        /**
         * Checks if the user can download files into the platform. this means whether it can interact with the platform's integrated media sources(e.g 4chan)
         * Any user that can read content can download the files(so long as they are not private).
         * @returns {boolean}
         */
        canDownloadFiles = this.#userCanChecker("download_files", true);

        /**
         * Checks if the user can upload files into the platform. 
         * @returns {boolean}
         */
        canUploadFiles = this.#userCanChecker("upload_files", true);

        /**
         * Checks if the user can create new Dungeons/CategoryClusters.
         * @returns {boolean}
         */
        canCreateDungeons = this.#userCanChecker("clusters_create", true);

        /**
         * Check if the user sync clusters. Which means the CategoriesService will re-scan the cluster's fs root for missing files or check if
         * any files in the db are missing.
         * @returns {boolean}
         */
        canSyncClusters = this.#userCanChecker("clusters_sync", true);
        
        /**
         * Check if the user can delete a cluster. This will remove the cluster from the database but not touch the files in the fs.
         * @returns {boolean}
         */
        canDropClusters = this.#userCanChecker("clusters_drop", true);

        /**
         * Check if the user can modify a cluster's content. This means the user can move medias around, rename categories, move categories, delete medias and categories.
         * This is only valid for public clusters, private clusters have their own variation of this grant.
         * @returns {boolean}
         */
        canPublicContentAlter = this.#userCanChecker("clusters_content_alter", true);

        /**
         * Checks if the user can read content from public clusters. 
         * Hint: This is the lowest level of access and every role should at least have this grant.
         * @returns {boolean}
         */
        canReadPublicContent = this.#userCanChecker("clusters_content_read", true);

        /**
         * Check if the user read content from private clusters.  
         * @returns {boolean}
         */
        canReadPrivateContent = this.#userCanChecker("clusters_content_read_private", true);

        /**
         * Check if the user can alter private clusters content. This means the user can move medias around, rename categories, move categories, delete medias and categories.
         * @returns {boolean}
         */
        canPrivateContentAlter = this.#userCanChecker("clusters_content_alter_private", true);

    /*=====  End of Grant Checkers  ======*/

}

/**
 * A user account information, this is used to allow super admins to modify a user account. There is no server endpoint to retrieve all this information 
 * at once, you need access to a UserEntry or a UserIdentity to get this information. to get a list of user entries use the `getAllUsers` function. to get the
 * roles of a user use the `getUserRoles` function. 
 */
export class UserAccount {
    /**
     * @type {string}
     */
    #the_uuid;

    /**
     * @type {string}
     */
    #the_username;

    /**
     * @type {RoleTaxonomy[]}
     */
    #the_roles;

    /**
     * @type {UserIdentity}
     */
    #the_identity;

    /**
     * @type {number}
     */
    #the_highest_role_hierarchy;

    /**
     * @param {string} uuid
     * @param {string} username
     * @param {RoleTaxonomy[]} roles
     */
    constructor(uuid, username, roles) {
        this.#the_uuid = uuid;
        this.#the_username = username;
        this.#the_roles = roles;
        this.#the_highest_role_hierarchy = findHighestRoleHierarchy(roles);

        let compiled_grants = compileGrants(roles);

        this.#the_identity = new UserIdentity({
            uuid: uuid,
            username: username,
            role_hierarchy: this.#the_highest_role_hierarchy,
            grants: compiled_grants
        });
    }

    /**
     * The account's unique identifier.
     * @type {string}
     */
    get UUID() {
        return this.#the_uuid;
    }

    /**
     * The account's username.
     * @type {string}
     */
    get Username() {
        return this.#the_username;
    }

    /**
     * The account's roles.
     * @type {RoleTaxonomy[]}
     */
    get Roles() {
        return this.#the_roles;
    }

    /**
     * The user's highest role hierarchy. smaller number means higher role, 0 been the highest.
     * @type {number}
     */
    get HighestRoleHierarchy() {
        return this.#the_highest_role_hierarchy;
    }

    /**
     * The user's identity. this is used to determine what the user can do and how to render different parts of the application.
     * @type {UserIdentity}
     */
    get Identity() {
        return this.#the_identity;
    }

    /**
     * Checks if the user has a given role by its role label.
     * @param {string} role_label
     * @returns {boolean}
     */
    isInRole(role_label) {
        return this.#the_roles.some(r => r.RoleLabel === role_label);
    }    
}

/**
 * A role taxonomy is: a set of grants that a role has, users gain grants by been assigned to a role(a user can have multiple roles).
 * the role hierarchy is a way to ensure lower hierarchy never have more grants than higher hierarchy roles. when a new grant is
 * add to a role, all the roles with a higher hierarchy will automatically have that grant as well. Also it can serve as a hint
 * to what the user can do, for example, a role with a hierarchy of 0 is a super admin, it can do anything in the platform and is
 * the only role that can grant privileges to other roles. other roles can create users, but those users will be visitors by default(which
 * is the lowest role in the platform, unless changed), and they can only be added to other roles by a role with the grant_option
 * which as said, is only the super admin role.
 */
export class RoleTaxonomy {
    /** @type {string} */
    #the_role_label;

    /** @type {number} */
    #the_role_hierarchy;

    /** @type {string[]} */
    #the_role_grants;

    /**
     * @param {RoleTaxonomyParams} param0
     */    
    constructor({ role_label, role_hierarchy, role_grants }) {
        this.#the_role_label = role_label;
        this.#the_role_hierarchy = role_hierarchy;
        this.#the_role_grants = role_grants;
    }

    /**
     * The role's identifier. is human readable
     * @type {string}
     */
    get RoleLabel() {
        return this.#the_role_label;
    }

    /**
     * The role hierarchy of the role. the smaller the number, the higher the role. 0 being the highest privilege role.
     * @type {number}
     */
    get RoleHierarchy() {
        return this.#the_role_hierarchy;
    }

    /**
     * The grants of the role. these are a list of capabilities that the role has. ALL_PRIVILEGES means the role can do anything except grant privileges,
     * that also requires the 'grant_privileges' grant. grants can be any random string that any service, be it core service or plugin service, can define.
     * @type {string[]}
     */
    get RoleGrants() {
        return this.#the_role_grants;
    }
}

/**
 * Checks if a given grant is an admin grant. basically a grant that if removed the system would be rendered unusable.
 * @param {string} grant
 * @returns {boolean}
 */
export const isSystemGrant = (grant) => {
    return grant === ALL_PRIVILEGES_GRANT || grant === GRANT_OPTION;
}

/**
 * Compiles all the grants of a list of role taxonomies into a string list.
 * @param {RoleTaxonomy[]} role_taxonomies
 * @returns {string[]}
 */
export const compileGrants = (role_taxonomies) => {
    /** @type {string[]} */
    let grants = [];

    role_taxonomies.forEach((taxonomy) => {
        grants = grants.concat(taxonomy.RoleGrants);
    });

    return grants;
}

/**
 * creates a user with the provided credentials
 * @param {string} username
 * @param {string} secret
 * @returns {Promise<boolean>}
 */
export const createNewUser = async (username, secret) => {
    let was_created = false;
    let request = new PostCreateUserRequest(username, secret);

    let response = await request.do();

    was_created =  response.Created;

    return was_created;
}

/**
 * Creates the initial super admin with the provided credentials and the initial setup secret. the server must be in initial setup mode which means that
 * there are no users registered. If this function is called when the server is not in this state, the request will fail even if a valid user claim token is set.
 * @param {string} username
 * @param {string} secret
 * @param {string} initial_setup_secret
 * @returns {Promise<boolean>}
 */
export const createInitialUser = async (username, secret, initial_setup_secret) => {
    let was_created = false;
    let request = new PostCreateInitialUserRequest(username, secret, initial_setup_secret);

    let response = await request.do();

    was_created = response.Created;

    return was_created;
}

/**
 * Returns whether the server has completed the initial setup. if not, a super admin user can be created without user authentication
 * but requires the `initial-setup-secret` which the user defines on the settings.json of the users service. 
 * @returns {Promise<boolean>}
 */
export const isUsersInInitialSetupMode = async () => {
    let is_initial_setup = false;
    let request = new GetIsInitialSetupRequest();

    let response = await request.do();

    if (response.Ok) {
        is_initial_setup = response.data.response;
    }

    return is_initial_setup;
}

/**
 * Logs the user in and returns a user identity. the user identity only serves as an informational object. when this function is called and succeeds,
 * the server will set a cookie with a signed jwt token that only it can read. this is how the services will determine what can and cannot be done by the user.
 * @param {string} username
 * @param {string} secret
 * @returns {Promise<UserIdentity | null>}
 */
export const loginPlatformUser = async (username, secret) => {
    let user_identity = null;
    let request = new GetUserSignAccessRequest(username, secret);

    let response = await request.do();

    if (response.Ok && response.data.granted) {
        user_identity = new UserIdentity(response.data.user_data);
    }

    return user_identity;
}

/**
 * Logs the user out. this will invalidate the user's access token and remove the cookie that the server uses to authenticate the user.
 * @returns {Promise<void>}
 */
export const logoutPlatformUser = async () => {
    let request = new GetUserSignOutRequest();

    await request.do();
}

/**
 * Validate the user's access token. this is used to check if the user is still logged in and the access token is not expired.
 * @returns {Promise<boolean>}
 */
export const validateUserAccessToken = async () => {
    let is_valid = false;
    let request = new GetUserAccessTokenValidationRequest();

    let response = await request.do();

    if (response.Ok) {
        is_valid = response.data.response;
    }

    return is_valid;
}

/**
 * Fetches the user identity of the currently logged in user. this is used to get the user's metadata and capabilities.
 * @returns {Promise<UserIdentity | null>}
 */
export const getCurrentUserIdentity = async () => {
    let user_identity = null;
    let request = new GetUserIdentityRequest();

    let response = await request.do();

    if (response.Ok) {
        user_identity = new UserIdentity(response.data);
    }

    return user_identity;
}

/**
 * Fetches all the users in the platform as a list of user entries. this is used to display the users in the frontend. and 
 * requires the user to have the 'read_users' grant.
 * @returns {Promise<UserEntry[]>}
 */
export const getAllUsers = async () => {
    /**
     * @type {UserEntry[]}
     */
    let users = [];
    let request = new GetAllUsersRequest();

    let response = await request.do();

    if (response.Ok) {
        users = response.data;
    }

    return users;
}

/**
 * Fetches all the role labels(strings) in the platform. this is used to display the existing roles in the frontend. 
 * Requires the user to have the 'read_users' grant.
 * @returns {Promise<string[]>} 
 */
export const getAllRoleLabels = async () => {
    /**
     * @type {string[]}
     */
    let role_labels = [];
    let request = new GetAllRoleLabelsRequest();

    let response = await request.do();

    if (response.Ok) {
        role_labels = response.data;
    }

    return role_labels;
}

/**
 * Fetches all the grants in the platform. this is are the capabilities that can be assigned to a role. requires the user to have the 'grant_option' grant.
 * @returns {Promise<string[]>}
 */
export const getAllGrants = async () => {
    /** @type {string[]} */
    let grants = [];
    let request = new GetAllGrantsRequest();

    let response = await request.do();

    if (response.Ok) {
        grants = response.data;
    }

    return grants;
}

/**
 * Registers the given string as a grant to be used by roles in the system. a grant cannot be assigned to a role before it is registered.
 * requires the user to have the 'grant_option' grant.
 * @param {string} grant_label
 * @returns {Promise<boolean>}
 */
export const createGrant = async (grant_label) => {
    let was_created = false;
    let request = new PostCreateGrantRequest(grant_label);

    let response = await request.do();

    was_created = response.data;

    return was_created;
}

/**
 * Adds a grant to a role. this is used to give a role a new capability. Roles with higher hierarchy will automatically have this grant as well.
 * requires the user to have the 'grant_option' grant.
 * @param {string} role_label
 * @param {string} grant
 * @returns {Promise<boolean>}
 */
export const addGrantToRole = async (role_label, grant) => {
    let was_linked = false;
    let request = new PostLinkGrantToRoleRequest(role_label, grant);

    let response = await request.do();

    was_linked = response.data;

    return was_linked;
}

/**
 * Requests the role taxonomy of a given role label. requires the 'grant_option' grant.
 * @param {string} role_label
 * @returns {Promise<RoleTaxonomy | null>}
 */
export const getRoleTaxonomy = async (role_label) => {
    let role_taxonomy = null;
    let request = new GetRoleTaxonomyRequest(role_label);

    let response = await request.do();

    if (response.Ok) {
        role_taxonomy = new RoleTaxonomy(response.data);
    }

    return role_taxonomy;
}

/**
 * Use this request when you want to know what grants will a newly created role inherit. Get all role taxonomies that are directly below a given role hierarchy.
 * for example, assume the system has roles with the following hierarchy: [0, 2, 3, 7, 8, 8 , 10]. If 4 is passed then it will return 2 taxonomies, with hierarchies 8. 
 * if 9 is passed then it will return 1 taxonomy with hierarchy 10. you can use the grants in these taxonomies to know what grants a role will inherit.
 * as you could've guessed, this request requires the 'grant_option' grant.
 * @param {number} role_hierarchy
 * @returns {Promise<RoleTaxonomy[]>}
 */
export const getRoleTaxonomiesBelowHierarchy = async (role_hierarchy) => {
    /** @type {RoleTaxonomy[]} */
    let role_taxonomies = [];
    let request = new GetRoleTaxonomiesBelowHierarchyRequest(role_hierarchy);

    let response = await request.do();

    if (response.Ok) {
        role_taxonomies = response.data.map((taxonomy) => new RoleTaxonomy(taxonomy));
    }

    return role_taxonomies;
}

/**
 * Creates a new role with the given role label, hierarchy and grants. the grants must exist before the role is created otherwise the request will fail.
 * requires the 'grant_option' grant.
 * @param {string} role_label
 * @param {number} role_hierarchy
 * @param {string[]} role_grants
 * @returns {Promise<boolean>}
 */
export const createRole = async (role_label, role_hierarchy, role_grants) => {
    let was_created = false;
    const tmp_taxonomy = new RoleTaxonomy({ role_label, role_hierarchy, role_grants });

    let request = new PostCreateRoleRequest(tmp_taxonomy);

    let response = await request.do();

    was_created = response?.data ?? false;

    return was_created;
}

/**
 * Adds a user to role. requires the 'grant_option' grant.
 * @param {string} username
 * @param {string} role_label
 * @returns {Promise<boolean>}
 */
export const addUserToRole = async (username, role_label) => {
    let was_patched = false;
    let request = new PatchUserRolesRequest(username, role_label);

    let response = await request.do();

    was_patched = response.data;

    return was_patched; 
}

/**
 * Deletes the link between a user and a role. requires the 'grant_option' grant.
 * @param {string} username
 * @param {string} role_label
 * @returns {Promise<boolean>}
 */
export const removeUserFromRole = async (username, role_label) => {
    let was_removed = false;
    let request = new DeleteUserFromRoleRequest(username, role_label);

    let response = await request.do();

    was_removed = response.data;

    return was_removed;
}

/**
 * Removes a given grant from a role. requires the 'grant_option' grant.
 * @param {string} role_label
 * @param {string} grant
 * @returns {Promise<boolean>}
 */
export const removeGrantFromRole = async (role_label, grant) => {
    let was_removed = false;
    let request = new DeleteGrantFromRoleRequest(role_label, grant);

    let response = await request.do();

    was_removed = response.data;

    return was_removed;
}

/**
 * Gets all the role_labels of roles that a given username is part of. requires the 'grant_option' grant.
 * @param {string} username
 * @returns {Promise<string[]>}
 */
export const getUserRoles = async (username) => {
    /** @type {string[]} */
    let roles = [];
    let request = new GetUserRolesRequest(username);

    let response = await request.do();

    if (response.Ok) {
        roles = response.data;
    }

    return roles;
}

/**
 * Deletes a user account from the platform. requires the 'delete_users' grant.
 * @param {string} user_uuid
 * @returns {Promise<boolean>}
 */
export const deleteUser = async (user_uuid) => {
    let was_deleted = false;
    let request = new DeleteUserRequest(user_uuid);

    let response = await request.do();

    was_deleted = response.data;

    return was_deleted;
}

/**
 * Deletes a grant from the system. requires 'grant_option' grant. Deleting a grant does not mean that actions that require that grant will stop requiring it.
 * but it will ensure that no user(not including the super admin and admins as they have the ALL_PRIVILEGES grant which includes every grant except for grant_option, which 
 * only the super admin has) will have that grant as it will be removed from all roles. Also until is readded, any attempts to link it to a role will fail.
 * @param {string} grant_label
 * @returns {Promise<boolean>}
 */
export const deleteGrant = async (grant_label) => {
    let was_deleted = false;
    let request = new DeleteGrantRequest(grant_label);

    let response = await request.do();

    was_deleted = response.data;

    return was_deleted;
}

/**
 * Deletes a role from the system. all users will be unassociated from the role. a role with hierarchy 0 cannot be deleted. requires 'grant_option' grant.
 * @param {string} role_label
 * @returns {Promise<boolean>}
 */
export const deleteRole = async (role_label) => {
    let was_deleted = false;
    let request = new DeleteRoleRequest(role_label);

    let response = await request.do();

    was_deleted = response.data;

    return was_deleted;
}

/**
 * Changes the password of the user. requires the 'modify_users' grant.
 * @param {string} uuid
 * @param {string} username
 * @param {string} secret
 * @returns {Promise<boolean>}
 */
export const changeUserPassword = async (uuid, username, secret) => {
    let was_changed = false;
    let request = new PutChangeUserPasswordRequest(uuid, username, secret);

    let response = await request.do();

    was_changed = response.data;

    return was_changed;
}

/**
 * Changes the username of the user. requires the 'modify_users' grant.
 * @param {string} uuid
 * @param {string} username
 * @returns {Promise<boolean>}
 */
export const changeUserUsername = async (uuid, username) => {
    let was_changed = false;
    let request = new PutChangeUsernameRequest(uuid, username);

    let response = await request.do();

    was_changed = response.data;

    return was_changed;
}

/**
 * From a user entry object, fetches the all the necessary information to create a user account instance.
 * @param {UserEntry} user_entry
 * @returns {Promise<UserAccount | null>}
 */
export const getUserAccountFromUserEntry = async (user_entry) => {
    /** @type {string[]} */
    const user_role_labels = await getUserRoles(user_entry.username);

    const user_roles = [];

    for (let role_label of user_role_labels) {
        let role_taxonomy = await getRoleTaxonomy(role_label);

        if (role_taxonomy) {
            user_roles.push(role_taxonomy);
        }
    }

    const user_account = new UserAccount(user_entry.uuid, user_entry.username, user_roles);

    return user_account;
}

/**
 * Returns the highest role hierarchy in any of the roles the user has. smaller number means higher role, 0 been the highest.
 * In the case the user has no roles, it will return Infinity.
 * @param {RoleTaxonomy[]} the_roles
 * @returns {number}
 */
export const findHighestRoleHierarchy = the_roles => {
    let highest_found_hierarchy =  Infinity;

    for (let role of the_roles) {
        if (role.RoleHierarchy < highest_found_hierarchy) {
            highest_found_hierarchy = role.RoleHierarchy;
        }
    }

    return highest_found_hierarchy;
}