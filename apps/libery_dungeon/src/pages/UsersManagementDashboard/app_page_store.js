import { writable } from "svelte/store";

/**
 * All the existing roles in the platform.
 * @type {import("svelte/store").Writable<string[]>}
 */
export const all_roles = writable([]);

/**
 * All the existing users in the platform.
 * @type {import("svelte/store").Writable<import('@models/Users').UserEntry[]>}
 */
export const all_users = writable([]);

/**
 * All the grants available to be assigned to roles
 * @type {import("svelte/store").Writable<string[]>}
 */ 
export const all_grants = writable([]);

/**
 * Whether role mode or user mode is currently enabled. Which changes the displayed tools and behavior of certain hotkeys.
 * @type {import("svelte/store").Writable<boolean>}
 */
export const role_mode_enabled = writable(false);