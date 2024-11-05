import { writable } from "svelte/store";

/**
* A store for shared state among different TaxonomyTags components. This store assumes that there can only be one TaxonomyTags component been interacted with at a time. 
* Event if there are multiple instances mounted, if somehow you implement a feature that interacts with multiple instances at the same time, this store might not work as expected.
*/

/**
* The last DungeonTag focused by a TaxonomyTags hotkey action.
 * @type {import('svelte/store').Writable<import('@models/DungeonTags').DungeonTag | null>}
 */
export const last_keyboard_focused_tag = writable(null);
