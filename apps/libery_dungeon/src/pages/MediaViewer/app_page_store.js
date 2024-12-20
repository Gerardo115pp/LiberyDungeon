import { writable } from "svelte/store";
import { current_category } from "@stores/categories_tree";

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