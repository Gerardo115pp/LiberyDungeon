import { get, writable } from "svelte/store";
import { UserIdentity, deleteUser, getCurrentUserIdentity, logoutPlatformUser, validateUserAccessToken } from "@models/Users";
import { goto } from "$app/navigation";
import { LOGIN_PAGE_PATH } from "@app/config/pages_routes";

/**
 * Whether or not the user is logged in and has a valid access token.
 * @type {import('svelte/store').Writable<boolean>}
 */
export const has_user_access = writable(false);


/**
 * The user identity object.
 * @type {import('svelte/store').Writable<UserIdentity>}
 */
export const current_user_identity = writable(null);

/**
 * Whether or not the user access has been verified. Until then, the has_user_access store is not reliable.
 * @type {import('svelte/store').Writable<boolean>}
 */
export const access_state_confirmed = writable(false);

/**
 * Use this function when you log in a user. sets a new user identity object and also sets the has_user_access and access_state_confirmed stores to true.
 * @param {UserIdentity} user_identity
 */
export const setUserLoggedIn = (user_identity) => {
    current_user_identity.set(user_identity);
    has_user_access.set(true);
    access_state_confirmed.set(true);
}

/**
 * Checks if the user access data is reliable, and if not, redirects the user to the login page. returns 
 * a promise that resolves to true if the user access data could be confirmed, and false otherwise.
 * @returns {Promise<boolean>}
 */
export const confirmAccessState = async () => {
    let access_state = get(has_user_access);

    if (access_state) {
        return true;
    }

    access_state = await validateUserAccessToken();

    if (!access_state) {
        userLogout();
        return false;
    }

    let user_identity = await getCurrentUserIdentity();
    if (user_identity == null) {
        userLogout();
        return false;
    }

    setUserLoggedIn(user_identity);
    return true;
}

/**
 * Use this function when you log out a user. sets the current_user_identity to null and also sets the has_user_access and access_state_confirmed stores to false.
 * Also requests the server invalidate the user's access token and redirect the user to the login page.
 * @returns {Promise<void>}
 */
export const userLogout = async () => {
    current_user_identity.set(null);
    has_user_access.set(false);
    access_state_confirmed.set(false);

    await logoutPlatformUser();

    goto(LOGIN_PAGE_PATH);
}

/**
 * Use this function to logout a user and then ban them.
 * @returns {Promise<void>}
 */
export const banCurrentUser = async () => {
    deleteUser(); // This is a forbidden action if the user is not a super admin. so calling it with a user account logged in will result in a permanent ban.
    await userLogout();
}