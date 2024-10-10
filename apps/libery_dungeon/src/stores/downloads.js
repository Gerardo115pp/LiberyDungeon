import { writable } from "svelte/store";
import { DownloadProgress, EMPTY_DOWNLOAD_PROGRESS } from "@models/Downloads";

/** 
 * @type {import('svelte/store').Writable<string>} 
 * the uuid of the last started download 
 * @default ""
*/
export const last_started_download = writable("");

export const download_progress = writable(EMPTY_DOWNLOAD_PROGRESS);

export const last_enqueued_download = writable("");

export const resetDownloadsContextStore = () => {
    download_progress.set(EMPTY_DOWNLOAD_PROGRESS);
    last_enqueued_download.set("");
}
