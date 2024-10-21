import { getClusterTags } from "@models/DungeonTags";
import { get, writable } from "svelte/store";

/**
* All the available tags in the current cluster as list of TaxonomyTags objects.
 * @type {import('svelte/store').Writable<import('@models/DungeonTags').TaxonomyTags[]>>}
 */
export const cluster_tags = writable([]);

/**
* The cluster for which the `cluster_tags` were fetched for.
 * @type {import('svelte/store').Writable<string>}
 */
export const last_cluster_domain = writable("");

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

    const new_taxonomy_tags = await getClusterTags(cluster_uuid);

    if (new_taxonomy_tags.length === 0) {
        return false;
    }
    
    last_cluster_domain.set(cluster_uuid);

    cluster_tags.set(new_taxonomy_tags);

    return true;
}