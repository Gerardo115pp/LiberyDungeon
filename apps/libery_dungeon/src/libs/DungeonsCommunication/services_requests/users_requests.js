import { HttpResponse, attributesToJson } from "../base";
import { users_server } from '../services';

export class GetIsInitialSetupRequest {
    static endpoint = `${users_server}/is-initial-setup`;

    /**
     * @returns {Promise<HttpResponse<import("../base").BooleanResponse>}
     */
    do = async () => {
        const response = await fetch(GetIsInitialSetupRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class PostCreateInitialUserRequest {
    static endpoint = `${users_server}/users`;

    /**
     * @param {string} username
     * @param {string} secret
     * @param {string} initial_setup_secret
     */
    constructor(username, secret, initial_setup_secret) {
        this.username = username;
        this.secret = secret;
        this._initial_setup_secret = initial_setup_secret;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('../base').SingleStringResponse>>}
     */
    do = async () => {
        const response = await fetch(`${PostCreateInitialUserRequest.endpoint}?initial-setup-secret=${this._initial_setup_secret}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class PostCreateUserRequest {
    static endpoint = `${users_server}/users`;

    /**
     * @param {string} username
     * @param {string} secret
     */
    constructor(username, secret) {
        this.username = username;
        this.secret = secret;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("../base").SingleStringResponse>>}
     */
    do = async () => {
        const response = await fetch(PostCreateUserRequest.endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class GetUserSignAccessRequest {
    static endpoint = `${users_server}/user-auth`;

    /**
     * @param {string} username
     * @param {string} secret
     */
    constructor(username, secret) {
        this.username = username;
        this.secret = secret;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<UserSignAccessResponse>>}
     * @typedef {Object} UserSignAccessResponse
     * @property {boolean} granted
     * @property {import('@models/Users').UserIdentityParams} user_data
     */
    do = async () => {
        /** @type {UserSignAccessResponse} */
        let data = {
            granted: false,
            user_data: {
                uuid: "",
                username: "",
                role_hierarchy: 0,
                grants: []
            }
        };

        const response = await fetch(`${GetUserSignAccessRequest.endpoint}?username=${this.username}&secret=${this.secret}`);

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class GetUserAccessTokenValidationRequest {
    static endpoint = `${users_server}/user-auth/verify`;

    /**
     * Requests the validation of a user access token(present as an http only cookie) the request will return 200
     * regardless of the token's validity. a boolan response indicates whether the token is valid or not.
     * @returns {Promise<HttpResponse<import("../base").BooleanResponse>}
     */
    do = async () => {
        const response = await fetch(GetUserAccessTokenValidationRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Requests the user's identity, relies on the user agent already having a valid access token stored.
 */
export class GetUserIdentityRequest {
    static endpoint = `${users_server}/users/identity`;

    /**
     * @returns {Promise<HttpResponse<import('@models/Users').UserIdentityParams>}
     */
    do = async () => {
        const response = await fetch(GetUserIdentityRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class GetUserSignOutRequest {
    static endpoint = `${users_server}/user-auth/logout`;

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const response = await fetch(GetUserSignOutRequest.endpoint);

        let logged_out = false;

        if (response.ok) {
            logged_out = true;
        }

        return new HttpResponse(response, logged_out);
    }
}

/**
 * Returns all the users as user entries(only the username and the uuid). to use this endpoint,
 * the user token must have the 'read_users' grant.
 */
export class GetAllUsersRequest {
    static endpoint = `${users_server}/users/read-all`;  

    /**
     * @returns {Promise<HttpResponse<import('@models/Users').UserEntry[]>}
     */
    do = async () => {
        const response = await fetch(GetAllUsersRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

export class GetAllRoleLabelsRequest {
    static endpoint = `${users_server}/roles/read-all`;

    /**
     * @returns {Promise<HttpResponse<string[]>}
     */
    do = async () => {
        const response = await fetch(GetAllRoleLabelsRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Returns all the grants which are just a string array. to use this endpoint, the user token must have the 'grant_option' grant, which only a
 * super admin has.
 */
export class GetAllGrantsRequest {
    static endpoint = `${users_server}/roles/grant/read-all`;

    /**
     * @returns {Promise<HttpResponse<string[]>}
     */
    do = async () => {
        const response = await fetch(GetAllGrantsRequest.endpoint);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Registers a new grant to be used in by any role. to use this endpoint, the user token must have the 'grant_option' grant, which only a
 * super admin has.
 */
export class PostCreateGrantRequest {
    static endpoint = `${users_server}/roles/grant`;

    /**
     * @param {string} new_grant
     */
    constructor(new_grant) {
        this.new_grant = new_grant;
    }


    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${PostCreateGrantRequest.endpoint}?new_grant=${this.new_grant}`;

        const response = await fetch(request_url, {
            method: "POST"
        });

        let created = false;

        if (response.status === 201) {
            created = true;
        }

        return new HttpResponse(response, created);
    }
}

/**
 * Links a grant to a role. to use this endpoint, the user token must have the 'grant_option' grant, which only a
 * super admin has.
 */
export class PostLinkGrantToRoleRequest {
    static endpoint = `${users_server}/roles/add-grant`;

    /**
     * @param {string} role_label
     * @param {string} grant
     */
    constructor(role_label, grant) {
        this.role_label = role_label;
        this.grant = grant;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${PostLinkGrantToRoleRequest.endpoint}?role_label=${this.role_label}&grant=${this.grant}`;

        const response = await fetch(request_url, {
            method: "POST"
        });

        let linked = false;

        if (response.status === 201) {
            linked = true;
        }

        return new HttpResponse(response, linked);
    }
}

/**
 * Requests the role taxonomy of a given role label. requires the 'grant_option' grant.
 */
export class GetRoleTaxonomyRequest {
    static endpoint = `${users_server}/roles/role`;

    /**
     * @param {string} role_label
     */
    constructor(role_label) {
        this.role_label = role_label;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Users').RoleTaxonomyParams>>}
     */
    do = async () => {
        const request_url = `${GetRoleTaxonomyRequest.endpoint}?role_label=${this.role_label}`;

        const response = await fetch(request_url);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Use this request when you want to know what grants will a newly created role inherit. Get all role taxonomies that are directly below a given role hierarchy.
 * for example, assume the system has roles with the following hierarchy: [0, 2, 3, 7, 8, 8 , 10]. If 4 is passed then it will return 2 taxonomies, with hierarchies 8. 
 * if 9 is passed then it will return 1 taxonomy with hierarchy 10. you can use the grants in these taxonomies to know what grants a role will inherit.
 * as you could've guessed, this request requires the 'grant_option' grant.
 */
export class GetRoleTaxonomiesBelowHierarchyRequest {

    static endpoint = `${users_server}/roles/below-hierarchy`;

    /**
     * @param {number} hierarchy
     */
    constructor(hierarchy) {
        this.hierarchy = hierarchy;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Users').RoleTaxonomyParams[]>}
     */
    do = async () => {
        const request_url = `${GetRoleTaxonomiesBelowHierarchyRequest.endpoint}?hierarchy=${this.hierarchy}`;

        const response = await fetch(request_url);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Creates a new role from a given taxonomy that is not already in the system(otherwise it will fail). requires the 'grant_option' grant.
 */
export class PostCreateRoleRequest {
    static endpoint = `${users_server}/roles/role`;

    /**
     * 
     * @param {import('@models/Users').RoleTaxonomy} role_taxonomy 
     */
    constructor(role_taxonomy) {
        this.role_label = role_taxonomy.RoleLabel;
        this.role_hierarchy = role_taxonomy.RoleHierarchy;
        this.role_grants = role_taxonomy.RoleGrants;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const response = await fetch(PostCreateRoleRequest.endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let created = false;

        if (response.status === 201) {
            created = true;
        }

        return new HttpResponse(response, created);
    }
}

/**
 * Adds a user to role. requires the 'grant_option' grant.
 */    
export class PatchUserRolesRequest {
    static endpoint = `${users_server}/users/role`;

    /**
     * @param {string} username
     * @param {string} role_label
     */
    constructor(username, role_label) {
        this.username = username;
        this.role_label = role_label;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${PatchUserRolesRequest.endpoint}?username=${this.username}&role_label=${this.role_label}`;

        const response = await fetch(request_url, {
            method: "PATCH"
        });

        let added = false;

        if (response.status === 200) {
            added = true;
        }

        return new HttpResponse(response, added);
    }
}

/**
 * Deletes the link between a user and a role. requires the 'grant_option' grant.
 */
export class DeleteUserFromRoleRequest {
    static endpoint = `${users_server}/users/role`;

    /**
     * @param {string} username
     * @param {string} role_label
     */
    constructor(username, role_label) {
        this.username = username;
        this.role_label = role_label;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${DeleteUserFromRoleRequest.endpoint}?username=${this.username}&role_label=${this.role_label}`;

        const response = await fetch(request_url, {
            method: "DELETE"
        });

        let deleted = false;

        if (response.status === 200) {
            deleted = true;
        }

        return new HttpResponse(response, deleted);
    }
}

/**
 * Deletes a grant from a role. requires the 'grant_option' grant.
 */
export class DeleteGrantFromRoleRequest {
    static endpoint = `${users_server}/roles/remove-grant`;

    /**
     * @param {string} role_label
     * @param {string} grant
     */
    constructor(role_label, grant) {
        this.role_label = role_label;
        this.grant = grant;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${DeleteGrantFromRoleRequest.endpoint}?role_label=${this.role_label}&grant=${this.grant}`;

        const response = await fetch(request_url, {
            method: "DELETE"
        });

        let deleted = false;

        if (response.status === 204) {
            deleted = true;
        }

        return new HttpResponse(response, deleted);
    }
}

/**
 * Gets all the role_labels of roles that a given username is part of. requires the 'grant_option' grant.
 */
export class GetUserRolesRequest {
    static endpoint = `${users_server}/users/roles`;

    /**
     * @param {string} username
     */
    constructor(username) {
        this.username = username;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<string[]>}
     */
    do = async () => {
        const request_url = `${GetUserRolesRequest.endpoint}?username=${this.username}`;

        const response = await fetch(request_url);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Deletes a user account. requires 'delete_users' grant.
 */
export class DeleteUserRequest {
    static endpoint = `${users_server}/users/user`;

    /**
     * @param {string} user_uuid
     */
    constructor(user_uuid) {
        this.user_uuid = user_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${DeleteUserRequest.endpoint}?user_uuid=${this.user_uuid}`;

        const response = await fetch(request_url, {
            method: "DELETE"
        });

        let deleted = false;

        if (response.status === 200) {
            deleted = true;
        }

        return new HttpResponse(response, deleted);
    }
}

/**
 * Deletes a grant from the system. requires 'grant_option' grant. Deleting a grant does not mean that actions that require that grant will stop requiring it.
 * but it will ensure that no user(not including the super admin and admins as they have the ALL_PRIVILEGES grant which includes every grant except for grant_option, which 
 * only the super admin has) will have that grant as it will be removed from all roles. Also until is readded, any attempts to link it to a role will fail.
 */ 
export class DeleteGrantRequest {
    static endpoint = `${users_server}/roles/grant`;

    /**
     * @param {string} grant
     */
    constructor(grant) {
        this.grant = grant;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${DeleteGrantRequest.endpoint}?grant=${this.grant}`;

        const response = await fetch(request_url, {
            method: "DELETE"
        });

        let deleted = false;

        if (response.status === 204) {
            deleted = true;
        }

        return new HttpResponse(response, deleted);
    }
}

/**
 * Deletes a role from the system. all users will be unassociated from the role. a role with hierarchy 0 cannot be deleted. requires 'grant_option' grant.
 */
export class DeleteRoleRequest {
    static endpoint = `${users_server}/roles/role`;

    /**
     * @param {string} role_label
     */
    constructor(role_label) {
        this.role_label = role_label;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = `${DeleteRoleRequest.endpoint}?role_label=${this.role_label}`;

        const response = await fetch(request_url, {
            method: "DELETE"
        });

        let deleted = false;

        if (response.status === 204) {
            deleted = true;
        }

        return new HttpResponse(response, deleted);
    }
}

/**
 * Changes a user account's password. requires the 'modify_users' grant. 
 */
export class PutChangeUserPasswordRequest {
    static endpoint = `${users_server}/users/user/secret`;

    /**
     * @param {string} uuid
     * @param {string} username
     * @param {string} secret_hash
     */
    constructor(uuid, username, secret_hash) {
        this.uuid = uuid;
        this.username = username;
        this.secret_hash = secret_hash;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = PutChangeUserPasswordRequest.endpoint;

        const response = await fetch(request_url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let changed = false;

        if (response.status === 204) {
            changed = true;
        }

        return new HttpResponse(response, changed);
    }
}

/**
 * Changes a user account's username. requires the 'modify_users' grant.
 */
export class PutChangeUsernameRequest {
    static endpoint = `${users_server}/users/user/username`;

    /**
     * @param {string} uuid
     * @param {string} username
     */
    constructor(uuid, username) {
        this.uuid = uuid;
        this.username = username;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>}
     */
    do = async () => {
        const request_url = PutChangeUsernameRequest.endpoint;

        const response = await fetch(request_url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let changed = false;

        if (response.status === 204) {
            changed = true;
        }

        return new HttpResponse(response, changed);
    }
}