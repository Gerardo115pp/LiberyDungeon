import { getClusterUserDefinedTags } from "@models/DungeonTags";
import { get, writable } from "svelte/store";

/**
* All the available user-defined tags in the current cluster as list of TaxonomyTags objects.
 * @type {import('svelte/store').Writable<import('@models/DungeonTags').TaxonomyTags[]>}
 */
export const cluster_tags = writable([]);

/**
* The cluster for which the `cluster_tags` were fetched for.
 * @type {import('svelte/store').Writable<string>}
 */
export const last_cluster_domain = writable("");

/**
 * Whether the cluster_tags correctness with respect to the current_cluster has been checked.
 * @type {import('svelte/store').Writable<boolean>}
 */ 
export let cluster_tags_checked = writable(false);

/**
* Refreshes the `cluster_tags` store with the given cluster_uuid. If the cluster_uuid is the same as `tags_cluster_domain` then the store is not refreshed.
* Unless `force` is passed and set to true. Returns whether the operation was successful or not(the cluster had at least one TaxonomyTags object).
 * @param {string} cluster_uuid
 * @param {boolean} [force=false]
 * @returns {Promise<boolean>}
 */
export const refreshClusterTags = async (cluster_uuid, force = false) => {
    if (cluster_uuid === get(last_cluster_domain) && force !== true && get(cluster_tags).length > 0) {
        return true;
    }

    return refreshClusterTagsNoCheck(cluster_uuid);
}

/**
 * Refreshes the `cluster_tags` store with the given cluster_uuid without using the get function to check cache. Checking the cache externally and then calling this function is recommended as
 * the get function involves subscribing and unsubscribing immediately to read the value. This can be slow for multiple calls.
 * @param {string} cluster_uuid
 * @returns {Promise<boolean>}
 */
export const refreshClusterTagsNoCheck = async (cluster_uuid) => {
    console.log("Refreshing tags for cluster: " + cluster_uuid);
    const new_taxonomy_tags = await getClusterUserDefinedTags(cluster_uuid);

    if (new_taxonomy_tags.length === 0) {
        cluster_tags_checked.set(false);
        return false;
    }

    cluster_tags_checked.set(true);
    
    last_cluster_domain.set(cluster_uuid);

    cluster_tags.set(new_taxonomy_tags);

    return true;
}

/**
 * Resets the state of the dungeon_tags store.
 * @returns {void}
 */
export const resetDungeonTagsStore = () => {
    cluster_tags.set([]);
    last_cluster_domain.set("");
    cluster_tags_checked.set(false);
}