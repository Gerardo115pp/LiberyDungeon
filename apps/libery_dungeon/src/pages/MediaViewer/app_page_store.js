import { writable } from "svelte/store";

/**
 * @type {import('svelte/store').Writable<boolean>} the name of a category been created, if the value is n null then no category is being created
 * @default false
 */
export const create_subcategory = writable(false);

export const resetMediaViewerPageStore = () => {
    create_subcategory.set(false);
}