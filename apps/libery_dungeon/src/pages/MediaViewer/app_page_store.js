import { writable } from "svelte/store";
import { media_change_types } from "@models/WorkManagers";
import { MediaChangesManager } from "@models/WorkManagers";
import { InnerCategory } from "@models/Categories";
import { Stack, StackBuffer } from "@libs/utils";

/**
 * @type {import('svelte/store').Writable<boolean>} the name of a category been created, if the value is n null then no category is being created
 * @default false
 */
export const create_subcategory = writable(false);

/**
 * Whether the media tagging tool is mounted or not.
 * @type {import('svelte/store').Writable<boolean>}
 */
export const media_tagging_tool_mounted = writable(false);

/**
 * Whether the tagged medias tool is mounted or not.
 * @type {import('svelte/store').Writable<boolean>}
 * @default true
 */
export const tagged_medias_tool_mounted = writable(false);

export const resetMediaViewerPageStore = () => {
    create_subcategory.set(false);
    media_tagging_tool_mounted.set(false);
    tagged_medias_tool_mounted.set(false);
}

/* ---------------------------------- State --------------------------------- */



/**
 * The index corresponding to the active media on stores/categories_tree.current_category.content array of Medias
 * @type {import('svelte/store').Writable<number>}
 */
export const active_media_index = writable(0);

/**
 * @type {import('svelte/store').Writable<import('@models/Medias').Media | null>}
 */
export const shared_active_media = writable(null);

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
 * Used in non-linear navigation or in modes where media navigation is not dependent on the current category contents. These are the medias that have been visited in the past
 * @type {import('svelte/store').Writable<StackBuffer<string>>}
 */
export const previous_medias = writable(new StackBuffer(200));

/**
 * These are medias that were popped from the previous_medias stack, and are used to go back to the already visited 'next' medias. for example, on random navigation, the user goes back to the previous media. this pops an
 * element from the previous_medias stack and adds it to this stack, so that when the user goes forward again, he/she can see the same 'next' medias until this stack is empty again(when element of this stack are popped in
 * random mode, they are added back to the previous_medias stack)
 * @type {import('svelte/store').Writable<Stack<string>>}
 */
export const static_next_medias = writable(new Stack()); 

/**
 * The previous media index, used to go back to the previous media when random navigation is enabled
 * @type {import('svelte/store').Writable<number>}
 * @default 0
 */
export const previous_media_index = writable(0);

/**
 * The change made to the active media
 * @type {import('svelte/store').Writable<import("@models/WorkManagers").MediaChangeType>}
 * @default "Normal"
 * @see media_change_types
 */
export const active_media_change = writable(media_change_types.NORMAL);

/**
 * @type {import('svelte/store').Writable<MediaChangesManager>}
 */
export let media_changes_manager = writable(new MediaChangesManager());

/**
 * @type {import('svelte/store').Writable<InnerCategory | null>}
 * the category that medias will be moved to on auto move mode 
 */
export let auto_move_category = writable(null);

/**
 * @type {import('svelte/store').Writable<boolean>}
 * whether or not the auto move mode is enabled 
 */
export let auto_move_on = writable(false);

/**
 * @type {import('svelte/store').Writable<string | undefined>}
 * a store containing the name of the media viewer hotkeys context 
 */
export let media_viewer_hotkeys_context_name = writable(undefined);

/**
 * Automute enabled
 * @type {import('svelte/store').Writable<boolean>}
 * @default true
 */ 
export let automute_enabled = writable(true);
