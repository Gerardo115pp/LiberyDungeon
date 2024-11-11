import { writable } from "svelte/store";
import { ChanCatalogThread } from "@models/4Chan";

/**
 * @type {import('svelte/store').Writable<ChanCatalogThread | null>} the thread that was selected on the board catalog
 */
export const selected_thread = writable(null);

/**
 * @type {import('svelte/store').Writable<string | null>} the thread id that was selected on the board catalog
 * @description this is used to scroll to the thread when returning from the thread page. is reseted when the board changes
 */
export const selected_thread_id = writable(null);

selected_thread_id.subscribe(value => console.log("selected_thread_id: " + value));