import { writable } from "svelte/store";

/**
* All the available tags in the current cluster as list of TaxonomyTags objects.
 * @type {import('svelte/store').Writable<import('@models/DungeonTags').TaxonomyTags[]>>}
 */
export const cluster_tags = writable([]);