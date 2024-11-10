import { writable } from 'svelte/store';

/**
 * Whether the Category folder should use media thumbnails if they have any.
 * @type {import('svelte/store').Writable<boolean>}
 */
export const use_category_folder_thumbnails = writable(true);