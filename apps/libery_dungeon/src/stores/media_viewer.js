import { writable } from "svelte/store";
import { media_change_types } from "@models/WorkManagers";
import { MediaChangesManager } from "@models/WorkManagers";
import { InnerCategory } from "@models/Categories";

/**
 * The index corresponding to the active media on stores/categories_tree.current_category.content array of Medias
 * @type {import('svelte/store').Writable<number>}
*/
export const active_media_index = writable(0);

/**
 * Whether or not the media random navigation is enabled, which means that the next media will be a random media
 * @type {import('svelte/store').Writable<boolean>}
 * @default false
*/
export const random_media_navigation = writable(false);

/**
 * whether the media navigation tries to skip deleted medias 
 * @type {import('svelte/store').Writable<boolean>}  
*/
export let skip_deleted_medias = writable(true);

/**
 * The previous media index, used to go back to the previous media when random navigation is enabled
 * @type {import('svelte/store').Writable<number>}
 * @default 0
*/
export const previous_media_index = writable(0);

/**
 * @typedef {"Moved" | "Deleted" | "Normal"} MediaChangeType
*/

/**
 * The change made to the active media
 * @type {import('svelte/store').Writable<MediaChangeType>}
 * @default "Normal"
 * @see media_change_types
*/
export const active_media_change = writable(media_change_types.NORMAL);

/** @type {import('svelte/store').Writable<MediaChangesManager>} */
export let media_changes_manager = writable(new MediaChangesManager());

/** @type {import('svelte/store').Writable<InnerCategory>} the category that medias will be moved to on auto move mode */
export let auto_move_category = writable(null);


/** @type {import('svelte/store').Writable<boolean>} whether or not the auto move mode is enabled */
export let auto_move_on = writable(false);

/** @type {import('svelte/store').Writable<string>} a store containing the name of the media viewer hotkeys context */
export let media_viewer_hotkeys_context_name = writable(undefined);


/**
 * Automute enabled
 * @type {import('svelte/store').Writable<boolean>}
 * @default true
 */ 
export let automute_enabled = writable(true);

